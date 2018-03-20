package okta

import (
	"fmt"
	"net/http"
	"time"
)

// ErrorResponse represents a response from the Okta API when an error occurs.
type ErrorResponse struct {
	Response *http.Response
	Code     string       `json:"errorCode"`
	Summary  string       `json:"errorSummary"`
	Link     string       `json:"errorLink"`
	ID       string       `json:"errorId"`
	Causes   []ErrorCause `json:"errorCauses"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: (%d) %s - %v - %s %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Code, r.Summary, r.ID, r.Causes)
}

// ErrorCause represents on cause for an error
type ErrorCause struct {
	Summary string `json:"errorSummary"`
}

func (e *ErrorCause) Error() string {
	return fmt.Sprintf("%s", e.Summary)
}

// RateLimitError represents an error when RateLimits are exceeded.
type RateLimitError struct {
	Rate     Rate           // Rate specifies last known rate limit for the client
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
}

func (r *RateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, formatRateReset(r.Rate.Reset.Sub(time.Now())))
}

// formatRateReset formats d to look like "[rate reset in 2s]" or
// "[rate reset in 87m02s]" for the positive durations. And like "[rate limit was reset 87m02s ago]"
// for the negative cases.
//
// Borrowed from: https://github.com/google/go-github/
func formatRateReset(d time.Duration) string {
	isNegative := d < 0
	if isNegative {
		d *= -1
	}
	secondsTotal := int(0.5 + d.Seconds())
	minutes := secondsTotal / 60
	seconds := secondsTotal - minutes*60

	var timeString string
	if minutes > 0 {
		timeString = fmt.Sprintf("%dm%02ds", minutes, seconds)
	} else {
		timeString = fmt.Sprintf("%ds", seconds)
	}

	if isNegative {
		return fmt.Sprintf("[rate limit was reset %v ago]", timeString)
	}
	return fmt.Sprintf("[rate reset in %v]", timeString)
}
