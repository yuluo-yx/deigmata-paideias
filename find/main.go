package main

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed web/*
var webFS embed.FS

const (
	statusRunning  = "running"
	statusFinished = "finished"
	statusError    = "error"
)

type IPRecord struct {
	IP       string `json:"ip"`
	Location string `json:"location"`
}

type Node struct {
	Domain   string     `json:"domain"`
	Records  []IPRecord `json:"records"`
	Children []*Node    `json:"children"`
}

type Scan struct {
	ID         string    `json:"scanId"`
	Root       string    `json:"root"`
	Status     string    `json:"status"`
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt,omitempty"`
	Tree       *Node     `json:"tree"`
	Error      string    `json:"error,omitempty"`
}

type ScanRequest struct {
	Domain   string `json:"domain"`
	MaxDepth int    `json:"maxDepth"`
	Workers  int    `json:"workers"`
}

type Task struct {
	Domain string
	Depth  int
	Parent *Node
}

type ScanStore struct {
	mu    sync.RWMutex
	scans map[string]*Scan
}

type GeoResolver struct {
	provider string
	client   *http.Client
	cache    map[string]string
	mu       sync.RWMutex
	limit    chan struct{}
	apiURL   string
}

func main() {
	wordlist := loadWordlist("wordlist.txt")
	store := &ScanStore{scans: make(map[string]*Scan)}
	geo := newGeoResolver()

	mux := http.NewServeMux()
	webRoot, err := fs.Sub(webFS, "web")
	if err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(webRoot)))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var req ScanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}

		req.Domain = normalizeDomain(req.Domain)
		if req.Domain == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "domain required"})
			return
		}

		if req.MaxDepth <= 0 {
			req.MaxDepth = 2
		}
		if req.Workers <= 0 {
			req.Workers = 20
		}

		scan := store.newScan(req.Domain)
		go runScan(scan, req.MaxDepth, req.Workers, wordlist, store, geo)

		writeJSON(w, http.StatusAccepted, map[string]string{"scanId": scan.ID, "status": scan.Status})
	})

	mux.HandleFunc("/tree/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/tree/")
		if id == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "scan id required"})
			return
		}

		scan, ok := store.getScan(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "scan not found"})
			return
		}

		writeJSON(w, http.StatusOK, scan)
	})

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (s *ScanStore) newScan(root string) *Scan {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := randomID()
	scan := &Scan{
		ID:        id,
		Root:      root,
		Status:    statusRunning,
		StartedAt: time.Now(),
		Tree:      &Node{Domain: root, Records: []IPRecord{}},
	}

	s.scans[id] = scan
	return scan
}

func (s *ScanStore) getScan(id string) (*Scan, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	scan, ok := s.scans[id]
	return scan, ok
}

func runScan(scan *Scan, maxDepth int, workers int, wordlist []string, store *ScanStore, geo *GeoResolver) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := make(chan Task, workers*2)
	var workersWG sync.WaitGroup
	var tasksWG sync.WaitGroup

	visited := make(map[string]struct{})
	visitedMu := sync.Mutex{}

	rootNode := scan.Tree
	tasksWG.Add(1)
	queue <- Task{Domain: scan.Root, Depth: 0, Parent: nil}

	worker := func() {
		defer workersWG.Done()
		for task := range queue {
			func() {
				defer tasksWG.Done()

				if task.Domain == "" {
					return
				}

				if !markVisited(&visitedMu, visited, task.Domain) {
					return
				}

				ips := resolveIPs(ctx, task.Domain)
				records := buildRecords(ctx, geo, ips)

				var node *Node
				if task.Depth == 0 {
					rootNode.Records = records
					node = rootNode
				} else if len(records) > 0 {
					node = &Node{Domain: task.Domain, Records: records}
					if task.Parent != nil {
						appendChild(task.Parent, node)
					}
				}

				if task.Depth < maxDepth {
					parentForChildren := node
					if task.Depth == 0 {
						parentForChildren = rootNode
					}
					if parentForChildren != nil {
						for _, word := range wordlist {
							select {
							case <-ctx.Done():
								return
							default:
							}
							sub := word + "." + task.Domain
							tasksWG.Add(1)
							queue <- Task{Domain: sub, Depth: task.Depth + 1, Parent: parentForChildren}
						}
					}
				}
			}()
		}
	}

	workersWG.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}

	go func() {
		tasksWG.Wait()
		close(queue)
	}()

	workersWG.Wait()

	store.mu.Lock()
	scan.Status = statusFinished
	scan.FinishedAt = time.Now()
	store.mu.Unlock()
}

func resolveIPs(ctx context.Context, domain string) []string {
	resolver := net.DefaultResolver
	ips, err := resolver.LookupIP(ctx, "ip", domain)
	if err != nil {
		return nil
	}

	results := make([]string, 0, len(ips))
	seen := make(map[string]struct{})
	for _, ip := range ips {
		text := ip.String()
		if _, ok := seen[text]; ok {
			continue
		}
		seen[text] = struct{}{}
		results = append(results, text)
	}
	return results
}

func buildRecords(ctx context.Context, geo *GeoResolver, ips []string) []IPRecord {
	records := make([]IPRecord, 0, len(ips))
	for _, ip := range ips {
		records = append(records, IPRecord{
			IP:       ip,
			Location: geo.Lookup(ctx, ip),
		})
	}
	return records
}

func appendChild(parent *Node, child *Node) {
	if parent == nil || child == nil {
		return
	}
	parent.Children = append(parent.Children, child)
}

func markVisited(mu *sync.Mutex, visited map[string]struct{}, domain string) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := visited[domain]; ok {
		return false
	}
	visited[domain] = struct{}{}
	return true
}

func loadWordlist(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{"www", "api", "app", "cdn", "static", "admin", "portal"}
	}

	lines := strings.Split(string(data), "\n")
	words := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		words = append(words, line)
	}
	return words
}

func normalizeDomain(input string) string {
	value := strings.TrimSpace(strings.ToLower(input))
	if value == "" {
		return ""
	}

	if strings.Contains(value, "://") {
		parsed, err := url.Parse(value)
		if err == nil {
			value = parsed.Hostname()
		}
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	value = strings.TrimPrefix(value, "//")
	if strings.Contains(value, "/") {
		value = strings.Split(value, "/")[0]
	}
	if strings.Contains(value, ":") {
		parts := strings.Split(value, ":")
		value = parts[0]
	}

	return strings.TrimSpace(value)
}

func randomID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return hex.EncodeToString(buf)
}

func newGeoResolver() *GeoResolver {
	provider := strings.ToLower(strings.TrimSpace(os.Getenv("GEO_PROVIDER")))
	if provider == "" {
		provider = "ipapi"
	}
	return &GeoResolver{
		provider: provider,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		cache:  make(map[string]string),
		limit:  make(chan struct{}, 5),
		apiURL: strings.TrimSpace(os.Getenv("GEO_API_URL")),
	}
}

func (g *GeoResolver) Lookup(ctx context.Context, ip string) string {
	if g.provider == "none" {
		return ""
	}

	g.mu.RLock()
	if val, ok := g.cache[ip]; ok {
		g.mu.RUnlock()
		return val
	}
	g.mu.RUnlock()

	g.limit <- struct{}{}
	defer func() { <-g.limit }()

	location := g.lookupIPAPI(ctx, ip)

	g.mu.Lock()
	g.cache[ip] = location
	g.mu.Unlock()

	return location
}

func (g *GeoResolver) lookupIPAPI(ctx context.Context, ip string) string {
	url := g.apiURL
	if url == "" {
		url = "https://ipapi.co/" + ip + "/json/"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ""
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	type ipapiResp struct {
		City        string `json:"city"`
		Region      string `json:"region"`
		CountryName string `json:"country_name"`
		Error       bool   `json:"error"`
		Reason      string `json:"reason"`
	}

	var payload ipapiResp
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
	}

	if payload.Error {
		return ""
	}

	parts := []string{}
	for _, part := range []string{payload.City, payload.Region, payload.CountryName} {
		if strings.TrimSpace(part) != "" {
			parts = append(parts, part)
		}
	}

	return strings.Join(parts, ", ")
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
