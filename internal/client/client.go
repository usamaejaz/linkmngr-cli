package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		http: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

type Link struct {
	ID           int    `json:"id"`
	Domain       string `json:"domain,omitempty"`
	Link         string `json:"link,omitempty"`
	Slug         string `json:"slug,omitempty"`
	Destination  string `json:"destination,omitempty"`
	BrandID      int    `json:"brand_id,omitempty"`
	Clicks       int    `json:"clicks,omitempty"`
	UniqueClicks int    `json:"unique_clicks,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
}

type Brand struct {
	ID       int            `json:"id"`
	Name     string         `json:"name,omitempty"`
	Domains  []BrandDomain  `json:"domains,omitempty"`
	Users    []BrandUser    `json:"users,omitempty"`
	Settings map[string]any `json:"settings,omitempty"`
}

type BrandUser struct {
	Name        string   `json:"name,omitempty"`
	Email       string   `json:"email,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type BrandDomain struct {
	Domain string `json:"domain,omitempty"`
	Active bool   `json:"active,omitempty"`
	Type   string `json:"type,omitempty"`
}

type PaginatedResults[T any] struct {
	Total       int `json:"total"`
	PerPage     int `json:"perPage"`
	CurrentPage int `json:"currentPage"`
	NextPage    int `json:"nextPage"`
	LastPage    int `json:"lastPage"`
	Items       []T `json:"items"`
}

type ListLinksParams struct {
	Page    int
	BrandID int
	Domain  string
}

type CreateLinkRequest struct {
	Destination string `json:"destination"`
	Domain      string `json:"domain,omitempty"`
	Slug        string `json:"slug,omitempty"`
	BrandID     int    `json:"brand_id,omitempty"`
}

type StatsParams struct {
	Start    string
	End      string
	TimeUnit string
	GroupBy  string
	BrandID  int
}

type ListBrandsParams struct {
	Page int
}

type Page struct {
	ID             int            `json:"id"`
	Title          string         `json:"title,omitempty"`
	Image          string         `json:"image,omitempty"`
	Slug           string         `json:"slug,omitempty"`
	Description    string         `json:"description,omitempty"`
	CustomDomainID int            `json:"custom_domain_id,omitempty"`
	BrandID        int            `json:"brand_id,omitempty"`
	Clicks         int            `json:"clicks,omitempty"`
	UniqueClicks   int            `json:"unique_clicks,omitempty"`
	LastClickAt    string         `json:"last_click_at,omitempty"`
	Design         map[string]any `json:"design,omitempty"`
	FeaturedLinks  []any          `json:"featuredLinks,omitempty"`
	Items          []any          `json:"items,omitempty"`
	URL            string         `json:"url,omitempty"`
}

type ListPagesParams struct {
	Page           int
	BrandID        int
	Domain         string
	CustomDomainID int
	Slug           string
	Query          string
}

type BrandDomainInput struct {
	Domain string `json:"domain"`
}

type BrandObjectForRequest struct {
	Name    string             `json:"name,omitempty"`
	Domains []BrandDomainInput `json:"domains,omitempty"`
}

func (c *Client) ListLinks(ctx context.Context, params ListLinksParams) (PaginatedResults[Link], error) {
	q := url.Values{}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.BrandID > 0 {
		q.Set("brand_id", strconv.Itoa(params.BrandID))
	}
	if params.Domain != "" {
		q.Set("domain", params.Domain)
	}

	var out PaginatedResults[Link]
	if err := c.do(ctx, http.MethodGet, withQuery("/links", q), nil, &out); err != nil {
		return PaginatedResults[Link]{}, err
	}
	return out, nil
}

func (c *Client) GetLink(ctx context.Context, id int) (Link, error) {
	var out Link
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/links/%d", id), nil, &out); err != nil {
		return Link{}, err
	}
	return out, nil
}

func (c *Client) CreateLink(ctx context.Context, in CreateLinkRequest) (Link, error) {
	var out Link
	if err := c.do(ctx, http.MethodPost, "/links", in, &out); err != nil {
		return Link{}, err
	}
	return out, nil
}

func (c *Client) GetLinkStats(ctx context.Context, linkID int, params StatsParams) (map[string]any, error) {
	q := url.Values{}
	q.Set("start", params.Start)
	q.Set("end", params.End)
	q.Set("time_unit", params.TimeUnit)
	if params.GroupBy != "" {
		q.Set("group_by", params.GroupBy)
	}

	var out map[string]any
	if err := c.do(ctx, http.MethodGet, withQuery(fmt.Sprintf("/links/%d/stats", linkID), q), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) ListBrands(ctx context.Context, params ListBrandsParams) (PaginatedResults[Brand], error) {
	q := url.Values{}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}

	var out PaginatedResults[Brand]
	if err := c.do(ctx, http.MethodGet, withQuery("/brands", q), nil, &out); err != nil {
		return PaginatedResults[Brand]{}, err
	}
	return out, nil
}

func (c *Client) GetBrand(ctx context.Context, id int) (Brand, error) {
	var out Brand
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/brands/%d", id), nil, &out); err != nil {
		return Brand{}, err
	}
	return out, nil
}

func (c *Client) CreateBrand(ctx context.Context, in BrandObjectForRequest) (Brand, error) {
	var out Brand
	if err := c.do(ctx, http.MethodPost, "/brands", in, &out); err != nil {
		return Brand{}, err
	}
	return out, nil
}

func (c *Client) UpdateBrand(ctx context.Context, id int, in BrandObjectForRequest) (Brand, error) {
	var out Brand
	if err := c.do(ctx, http.MethodPatch, fmt.Sprintf("/brands/%d", id), in, &out); err != nil {
		return Brand{}, err
	}
	return out, nil
}

func (c *Client) DeleteBrand(ctx context.Context, id int) error {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("/brands/%d", id), nil, nil)
}

func (c *Client) CheckBrandDomain(ctx context.Context, id int, domain string) (BrandDomain, error) {
	in := map[string]string{"domain": domain}
	var out BrandDomain
	if err := c.do(ctx, http.MethodPost, fmt.Sprintf("/brands/%d/check_domain", id), in, &out); err != nil {
		return BrandDomain{}, err
	}
	return out, nil
}

func (c *Client) GetAnalytics(ctx context.Context, params StatsParams) (map[string]any, error) {
	q := url.Values{}
	q.Set("start", params.Start)
	q.Set("end", params.End)
	q.Set("time_unit", params.TimeUnit)
	if params.GroupBy != "" {
		q.Set("group_by", params.GroupBy)
	}
	if params.BrandID > 0 {
		q.Set("brand_id", strconv.Itoa(params.BrandID))
	}

	var out map[string]any
	if err := c.do(ctx, http.MethodGet, withQuery("/analytics", q), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) ListDomains(ctx context.Context) ([]string, error) {
	var out any
	if err := c.do(ctx, http.MethodGet, "/domains", nil, &out); err != nil {
		return nil, err
	}
	domains, ok := extractDomains(out)
	if !ok {
		return nil, fmt.Errorf("unexpected /domains response shape")
	}
	return domains, nil
}

func extractDomains(v any) ([]string, bool) {
	switch t := v.(type) {
	case []any:
		out := make([]string, 0, len(t))
		for _, item := range t {
			switch it := item.(type) {
			case string:
				if strings.TrimSpace(it) != "" {
					out = append(out, it)
				}
			case map[string]any:
				if d, ok := it["domain"].(string); ok && strings.TrimSpace(d) != "" {
					out = append(out, d)
				}
			}
		}
		return out, true
	case map[string]any:
		for _, key := range []string{"items", "domains", "data"} {
			if nested, exists := t[key]; exists {
				if out, ok := extractDomains(nested); ok {
					return out, true
				}
			}
		}
	case []string:
		return t, true
	}
	return nil, false
}

func (c *Client) ListPages(ctx context.Context, params ListPagesParams) (PaginatedResults[Page], error) {
	q := url.Values{}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.BrandID > 0 {
		q.Set("brand_id", strconv.Itoa(params.BrandID))
	}
	if params.Domain != "" {
		q.Set("domain", params.Domain)
	}
	if params.CustomDomainID > 0 {
		q.Set("custom_domain_id", strconv.Itoa(params.CustomDomainID))
	}
	if params.Slug != "" {
		q.Set("slug", params.Slug)
	}
	if params.Query != "" {
		q.Set("q", params.Query)
	}

	var out PaginatedResults[Page]
	if err := c.do(ctx, http.MethodGet, withQuery("/pages", q), nil, &out); err != nil {
		return PaginatedResults[Page]{}, err
	}
	return out, nil
}

func (c *Client) GetPage(ctx context.Context, id int) (Page, error) {
	var out Page
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/pages/%d", id), nil, &out); err != nil {
		return Page{}, err
	}
	return out, nil
}

func (c *Client) GetPageStats(ctx context.Context, pageID int, params StatsParams) (map[string]any, error) {
	q := url.Values{}
	q.Set("start", params.Start)
	q.Set("end", params.End)
	q.Set("time_unit", params.TimeUnit)
	if params.GroupBy != "" {
		q.Set("group_by", params.GroupBy)
	}

	var out map[string]any
	if err := c.do(ctx, http.MethodGet, withQuery(fmt.Sprintf("/pages/%d/stats", pageID), q), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetPageHits(ctx context.Context, pageID int) ([]map[string]any, error) {
	var out []map[string]any
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/pages/%d/hits", pageID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) WhoAmI(ctx context.Context) (map[string]any, error) {
	var out map[string]any
	if err := c.do(ctx, http.MethodGet, "/auth/user", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) RevokeToken(ctx context.Context) error {
	return c.do(ctx, http.MethodPost, "/auth/revoke", nil, nil)
}

func (c *Client) RawRequest(ctx context.Context, method, path string, in any) (any, error) {
	if path == "" {
		return nil, fmt.Errorf("path is required")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	var out any
	if err := c.do(ctx, strings.ToUpper(method), path, in, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func withQuery(path string, q url.Values) string {
	if len(q) == 0 {
		return path
	}
	return path + "?" + q.Encode()
}

func (c *Client) do(ctx context.Context, method, path string, in any, out any) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("perform request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(respBody))
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return fmt.Errorf("api %s %s failed (%d): %s", method, path, resp.StatusCode, msg)
	}

	if out == nil || len(respBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("parse response JSON: %w", err)
	}
	return nil
}
