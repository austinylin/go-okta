package okta

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	userAgent           = "go-okta"
	headerRateLimit     = "X-Rate-Limit-Limit"
	headerRateRemaining = "X-Rate-Limit-Remaining"
	headerRateReset     = "X-Rate-Limit-Reset"
)

type contextKey string

type service struct {
	client *Client
}

// Client represents an Okta API client.
type Client struct {
	httpClient *http.Client
	apiToken   string
	UserAgent  string
	BaseURL    *url.URL
	rateMu     sync.Mutex
	rateLimits [categories]Rate // Rate limits for the client as determined by the most recent API calls.
	common     service          // Reuse a single struct instead of allocating one for each service on the heap.
}

// Response represents a response from the Okta API.
type Response struct {
	*http.Response
	Pagination
	Rate
	OktaRequestID string
}

// Pagination represents the pagination primiatives of the Okta API.
type Pagination struct {
	Prev  string `json:"prev"`
	Next  string `json:"next"`
	Limit int    `json:"limit"`
}

// NewClient creates a new Okta API client.
func NewClient(apiToken string, paramBaseURL string, httpClient *http.Client) (*Client, error) {
	if len(apiToken) == 0 {
		return nil, errors.New("API Token is not present")
	}
	if len(paramBaseURL) == 0 {
		return nil, errors.New("Base URL is not present")
	}
	baseURL, _ := url.Parse(paramBaseURL)

	if !strings.HasSuffix(baseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", baseURL)
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := &Client{
		UserAgent:  userAgent,
		BaseURL:    baseURL,
		apiToken:   apiToken,
		httpClient: httpClient,
	}

	return client, nil
}

// NewRequest creates a new *http.Request that can be used to query the Okta API.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Do executes an http.Request with context, and returns the result, optionally decoding the body into the
// provided interface.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	req = req.WithContext(ctx)

	// Auth
	req.Header.Set("Authorization", fmt.Sprintf("SSWS %s", c.apiToken))

	// Check rate limits before we actually make the request
	rateLimitCategory := ctx.Value(rateLimitCategoryCtxKey).(rateLimitCategory)
	if err := c.checkRateLimitBeforeDo(req, rateLimitCategory); err != nil {
		return &Response{
			Response: err.Response,
			Rate:     err.Rate,
		}, err
	}

	// actually send the request
	resp, err := c.httpClient.Do(req)

	rateLimit := parseRate(resp)
	c.rateMu.Lock()
	c.rateLimits[rateLimitCategory] = rateLimit
	c.rateMu.Unlock()

	response := &Response{Response: resp}
	// TODO: Pagination?
	response.Rate = rateLimit

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}

func parseRate(r *http.Response) Rate {
	var rate Rate
	if limit := r.Header.Get(headerRateLimit); limit != "" {
		rate.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := r.Header.Get(headerRateReset); reset != "" {
		if v, _ := strconv.ParseInt(reset, 10, 64); v != 0 {
			rate.Reset = Timestamp{Time: time.Unix(v, 0)}
		}
	}
	return rate
}

// checkRateLimitBeforeDo does not make any network calls, but uses existing knowledge from
// current client state in order to quickly check if *RateLimitError can be immediately returned
// from Client.Do, and if so, returns it so that Client.Do can skip making a network API call unnecessarily.
// Otherwise it returns nil, and Client.Do should proceed normally.
func (c *Client) checkRateLimitBeforeDo(req *http.Request, rateLimitCategory rateLimitCategory) *RateLimitError {
	c.rateMu.Lock()
	rate := c.rateLimits[rateLimitCategory]
	c.rateMu.Unlock()
	if rate.Remaining == 0 && time.Now().Before(rate.Reset.Time) {
		// Create a fake response.
		resp := &http.Response{
			Status:     http.StatusText(http.StatusForbidden),
			StatusCode: http.StatusForbidden,
			Request:    req,
			Header:     make(http.Header),
			Body:       ioutil.NopCloser(strings.NewReader("")),
		}
		return &RateLimitError{
			Rate:     rate,
			Response: resp,
			Message:  fmt.Sprintf("API rate limit of %v still exceeded until %v, not making remote request.", rate.Limit, rate.Reset),
		}
	}

	return nil
}
