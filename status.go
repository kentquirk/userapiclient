// Implements a set of useful status-related tools.

package userapiclient

import (
	"fmt"
	"net/http"
)

// Type status holds the status return from an http request.
type Status struct {
	Code   int
	Reason string
}

// StatusError represents a Status as an error object.
type StatusError struct {
	Status
}

// NewStatus constructs a Status object; if no reason is provided, it uses the
// standard one.
func NewStatus(code int, reason string) Status {
	r := reason
	if r == "" {
		r = http.StatusText(code)
	}
	s := Status{code, r}
	return s
}

// String() converts a status to a printable string.
func (s Status) String() string { return fmt.Sprintf("%d %s", s.Code, s.Reason) }

// Error() renders a StatusError.
func (s *StatusError) Error() string {
	return s.Status.String()
}
