package ghapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const baseURL = "https://api.github.com"

type Client struct {
	httpClient *http.Client
	token      string
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}

type RateLimit struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	ResetAt   time.Time `json:"reset_at"`
}

type RepoContribution struct {
	Repository   string `json:"repository"`
	IssuesOpened int    `json:"issues_opened"`
	PRsOpened    int    `json:"prs_opened"`
	PRsMerged    int    `json:"prs_merged"`
	Total        int    `json:"total"`
}

type ContributionItem struct {
	Title      string `json:"title"`
	URL        string `json:"url"`
	Repository string `json:"repository"`
	Merged     bool   `json:"merged,omitempty"`
}

type OrgContributionStats struct {
	User                string             `json:"user"`
	Org                 string             `json:"org"`
	IssuesOpened        int                `json:"issues_opened"`
	PRsOpened           int                `json:"prs_opened"`
	PRsMerged           int                `json:"prs_merged"`
	Total               int                `json:"total"`
	RepositoryBreakdown []RepoContribution `json:"repository_breakdown"`
	Issues              []ContributionItem `json:"issues"`
	PRs                 []ContributionItem `json:"prs"`
	RepoStatsTruncated  bool               `json:"repo_stats_truncated"`
	RateLimit           RateLimit          `json:"rate_limit"`
}

type searchItem struct {
	RepositoryURL string `json:"repository_url"`
	HTMLURL       string `json:"html_url"`
	Title         string `json:"title"`
}

type searchResponse struct {
	TotalCount int          `json:"total_count"`
	Items      []searchItem `json:"items"`
	Message    string       `json:"message"`
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 20 * time.Second},
		token:      strings.TrimSpace(token),
	}
}

func (c *Client) FetchOrgContributionStats(ctx context.Context, user, org string) (OrgContributionStats, error) {
	stats := OrgContributionStats{User: user, Org: org}

	issuesQuery := fmt.Sprintf("is:issue author:%s org:%s", user, org)
	prsQuery := fmt.Sprintf("is:pr author:%s org:%s", user, org)
	mergedQuery := fmt.Sprintf("is:pr author:%s org:%s is:merged", user, org)

	issuesCount, issuesByRepo, issueItems, rate, issuesTruncated, err := c.searchAndGroupByRepo(ctx, issuesQuery)
	if err != nil {
		return OrgContributionStats{}, err
	}
	allPRCount, allPRByRepo, prItems, _, prsTruncated, err := c.searchAndGroupByRepo(ctx, prsQuery)
	if err != nil {
		return OrgContributionStats{}, err
	}
	mergedCount, mergedByRepo, mergedItems, _, mergedTruncated, err := c.searchAndGroupByRepo(ctx, mergedQuery)
	if err != nil {
		return OrgContributionStats{}, err
	}

	prsCount := allPRCount - mergedCount
	if prsCount < 0 {
		prsCount = 0
	}
	prsByRepo := subtractRepoCounts(allPRByRepo, mergedByRepo)

	mergedSet := make(map[string]struct{}, len(mergedItems))
	for _, item := range mergedItems {
		mergedSet[item.URL] = struct{}{}
	}

	for idx := range prItems {
		_, ok := mergedSet[prItems[idx].URL]
		prItems[idx].Merged = ok
	}

	stats.IssuesOpened = issuesCount
	stats.PRsOpened = prsCount
	stats.PRsMerged = mergedCount
	stats.Total = issuesCount + prsCount + mergedCount
	stats.RepositoryBreakdown = mergeRepoContributions(issuesByRepo, prsByRepo, mergedByRepo)
	stats.Issues = issueItems
	stats.PRs = prItems
	stats.RepoStatsTruncated = issuesTruncated || prsTruncated || mergedTruncated
	stats.RateLimit = rate

	return stats, nil
}

func subtractRepoCounts(all map[string]int, merged map[string]int) map[string]int {
	result := make(map[string]int, len(all))
	for repo, allCount := range all {
		remain := allCount - merged[repo]
		if remain > 0 {
			result[repo] = remain
		}
	}
	return result
}

func (c *Client) searchAndGroupByRepo(ctx context.Context, query string) (int, map[string]int, []ContributionItem, RateLimit, bool, error) {
	const perPage = 100

	repoCount := map[string]int{}
	items := make([]ContributionItem, 0)
	fetched := 0
	page := 1
	total := 0
	var rate RateLimit
	truncated := false

	for {
		u := baseURL + "/search/issues?q=" + url.QueryEscape(query) + "&per_page=" + strconv.Itoa(perPage) + "&page=" + strconv.Itoa(page)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			return 0, nil, nil, RateLimit{}, false, err
		}

		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		req.Header.Set("User-Agent", "github-contributor-cli")
		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return 0, nil, nil, RateLimit{}, false, err
		}

		var body searchResponse
		decodeErr := json.NewDecoder(resp.Body).Decode(&body)
		resp.Body.Close()
		if decodeErr != nil {
			return 0, nil, nil, RateLimit{}, false, decodeErr
		}

		rate = parseRateLimit(resp.Header)
		if resp.StatusCode >= http.StatusMultipleChoices {
			msg := strings.TrimSpace(body.Message)
			if msg == "" {
				msg = "GitHub API 请求失败"
			}
			return 0, nil, nil, rate, false, &APIError{StatusCode: resp.StatusCode, Message: msg}
		}

		if page == 1 {
			total = body.TotalCount
			if total > 1000 {
				truncated = true
			}
		}

		for _, item := range body.Items {
			repo := repoFromAPIURL(item.RepositoryURL)
			if repo == "" {
				continue
			}
			repoCount[repo]++
			items = append(items, ContributionItem{
				Title:      strings.TrimSpace(item.Title),
				URL:        strings.TrimSpace(item.HTMLURL),
				Repository: repo,
			})
		}

		fetched += len(body.Items)
		maxFetch := total
		if maxFetch > 1000 {
			maxFetch = 1000
		}
		if fetched >= maxFetch || len(body.Items) == 0 {
			break
		}
		page++
	}

	return total, repoCount, items, rate, truncated, nil
}

func repoFromAPIURL(raw string) string {
	const prefix = "https://api.github.com/repos/"
	if strings.HasPrefix(raw, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(raw, prefix))
	}
	return ""
}

func mergeRepoContributions(issuesByRepo, prsByRepo, mergedByRepo map[string]int) []RepoContribution {
	bucket := map[string]*RepoContribution{}

	upsert := func(repo string) *RepoContribution {
		entry, ok := bucket[repo]
		if !ok {
			entry = &RepoContribution{Repository: repo}
			bucket[repo] = entry
		}
		return entry
	}

	for repo, count := range issuesByRepo {
		entry := upsert(repo)
		entry.IssuesOpened = count
	}
	for repo, count := range prsByRepo {
		entry := upsert(repo)
		entry.PRsOpened = count
	}
	for repo, count := range mergedByRepo {
		entry := upsert(repo)
		entry.PRsMerged = count
	}

	result := make([]RepoContribution, 0, len(bucket))
	for _, entry := range bucket {
		entry.Total = entry.IssuesOpened + entry.PRsOpened + entry.PRsMerged
		result = append(result, *entry)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Total == result[j].Total {
			return result[i].Repository < result[j].Repository
		}
		return result[i].Total > result[j].Total
	})

	return result
}

func parseRateLimit(h http.Header) RateLimit {
	var result RateLimit

	if v, err := strconv.Atoi(h.Get("X-RateLimit-Limit")); err == nil {
		result.Limit = v
	}
	if v, err := strconv.Atoi(h.Get("X-RateLimit-Remaining")); err == nil {
		result.Remaining = v
	}
	if v, err := strconv.ParseInt(h.Get("X-RateLimit-Reset"), 10, 64); err == nil && v > 0 {
		result.ResetAt = time.Unix(v, 0)
	}

	return result
}
