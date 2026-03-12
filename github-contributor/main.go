package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/deigmata-paideias/github-contributor/internal/ghapi"
	"github.com/deigmata-paideias/github-contributor/internal/report"
)

func main() {
	var (
		user    = flag.String("user", "", "GitHub 用户名 (必填)")
		org     = flag.String("org", "", "GitHub 组织名 (必填)")
		output  = flag.String("output", "html", "输出格式: html 或 json")
		outFile = flag.String("out-file", "report.html", "输出文件路径")
		token   = flag.String("token", "", "GitHub Token (可选, 也可用环境变量 GITHUB_TOKEN)")
		timeout = flag.Duration("timeout", 20*time.Second, "请求超时时间")
		serve   = flag.Bool("serve", true, "当 output=html 时是否启动本地 Web 服务")
		port    = flag.Int("serve-port", 8080, "本地 Web 服务端口")
		open    = flag.Bool("open-browser", true, "启动服务后是否自动打开浏览器")
	)
	flag.Parse()

	if strings.TrimSpace(*user) == "" || strings.TrimSpace(*org) == "" {
		log.Fatalf("参数错误: --user 和 --org 均为必填")
	}

	resolvedToken := strings.TrimSpace(*token)
	if resolvedToken == "" {
		resolvedToken = strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	client := ghapi.NewClient(resolvedToken)
	stats, err := client.FetchOrgContributionStats(ctx, *user, *org)
	if err != nil {
		var apiErr *ghapi.APIError
		if errors.As(err, &apiErr) {
			log.Fatalf("GitHub API 错误: %s (status=%d)", apiErr.Message, apiErr.StatusCode)
		}
		log.Fatalf("查询失败: %v", err)
	}

	reportData := report.Data{
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05 MST"),
		Stats:       stats,
	}

	switch strings.ToLower(strings.TrimSpace(*output)) {
	case "html":
		if err := report.WriteHTML(*outFile, reportData); err != nil {
			log.Fatalf("生成 HTML 失败: %v", err)
		}
		fmt.Printf("HTML 报表已生成: %s\n", *outFile)
		if *serve {
			if err := serveHTMLReport(*outFile, *port, *open); err != nil {
				log.Fatalf("启动 Web 服务失败: %v", err)
			}
		}
	case "json":
		if err := report.WriteJSON(*outFile, reportData); err != nil {
			log.Fatalf("生成 JSON 失败: %v", err)
		}
		fmt.Printf("JSON 报表已生成: %s\n", *outFile)
	default:
		log.Fatalf("参数错误: 不支持的 --output 值 %q, 仅支持 html/json", *output)
	}
}

func serveHTMLReport(reportPath string, port int, autoOpen bool) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("无效端口: %d", port)
	}

	absPath, err := filepath.Abs(reportPath)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, absPath)
	})

	addr := fmt.Sprintf(":%d", port)
	url := fmt.Sprintf("http://localhost:%d", port)
	server := &http.Server{Addr: addr, Handler: mux}

	if autoOpen {
		_ = openBrowser(url)
	}

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stopCh)

	go func() {
		<-stopCh
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	log.Printf("Web 服务已启动: %s (Ctrl+C 退出)", url)
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
