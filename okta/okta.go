package okta

// Portions of this code are inspired by and borrowed from https://github.com/google/go-github
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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
	headerRequestID     = "X-Okta-Request-Id"
	envDebug            = "GO_OKTA_DEBUG"
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

	Apps   *AppsService
	Groups *GroupsService
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
	Prev string `json:"prev"`
	Next string `json:"next"`
	Self string `json:"self"`
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

	c := &Client{
		UserAgent:  userAgent,
		BaseURL:    baseURL,
		apiToken:   apiToken,
		httpClient: httpClient,
	}

	c.common.client = c
	c.Apps = (*AppsService)(&c.common)
	c.Groups = (*GroupsService)(&c.common)

	return c, nil
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

	// If we are in debug mode, log the request prior to adding the auth header.
	if os.Getenv(envDebug) != "" {
		reqDump, _ := httputil.DumpRequest(req, true)
		log.Printf("Request:\n %s\n", reqDump)
	}

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

	// If we are in debug mode, log the response.
	if os.Getenv(envDebug) != "" {
		respDump, _ := httputil.DumpResponse(resp, true)
		log.Printf("Response:\n %s\n", respDump)
	}

	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = url.String()
				return nil, e
			}
		}

		return nil, err
	}
	defer resp.Body.Close()

	rateLimit := parseRate(resp)
	c.rateMu.Lock()
	c.rateLimits[rateLimitCategory] = rateLimit
	c.rateMu.Unlock()

	response := &Response{Response: resp}

	response.Pagination = Pagination{}
	response.populatePageValues()

	response.Rate = rateLimit
	response.OktaRequestID = resp.Header.Get(headerRequestID)

	err = checkResponseForErrors(resp)
	if err != nil {
		return response, err
	}

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

// populatePageValues parses the HTTP Link response headers and populates the
// various pagination link values in the Response.
func (r *Response) populatePageValues() {
	if links, ok := r.Response.Header["Link"]; ok && len(links) > 0 {
		for _, link := range links {
			segments := strings.Split(strings.TrimSpace(link), ";")

			// link must at least have href and rel
			if len(segments) < 2 {
				continue
			}

			// ensure href is properly formatted
			if !strings.HasPrefix(segments[0], "<") || !strings.HasSuffix(segments[0], ">") {
				continue
			}

			// pull out the URL
			url, err := url.Parse(segments[0][1 : len(segments[0])-1])
			if err != nil {
				continue
			}

			for _, segment := range segments[1:] {
				switch strings.TrimSpace(segment) {
				case `rel="next"`:
					r.Pagination.Next = url.String()
				case `rel="prev"`:
					r.Pagination.Prev = url.String()
				case `rel="self"`:
					r.Pagination.Self = url.String()
				}

			}
		}
	}
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

// checkResponseForErrors checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range or equal to 202 Accepted.
// API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other
// response body will be silently ignored.
//
// The error type will be *RateLimitError for rate limit exceeded errors,
// *AcceptedError for 202 Accepted status codes,
// and *TwoFactorAuthError for two-factor authentication errors.
func checkResponseForErrors(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	switch {
	case r.StatusCode == http.StatusForbidden && r.Header.Get(headerRateRemaining) == "0":
		return &RateLimitError{
			Rate:     parseRate(r),
			Response: errorResponse.Response,
		}
	default:
		return errorResponse
	}
}
