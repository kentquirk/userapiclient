// Implements a set of useful status-related tools

package userapiclient

import (
	"fmt"
	"net/http"
)

type Status struct {
	Code   int
	Reason string
}

// StatusError reprents a Status as an error.
type StatusError struct {
	Status
}

func NewStatus(code int, reason string) Status {
	r := reason
	if r == "" {
		r = http.StatusText(code)
	}
	s := Status{code, r}
	return s
}

func (s Status) String() string { return fmt.Sprintf("%d %s", s.Code, s.Reason) }

func (s *StatusError) Error() string {
	return s.Status.String()
}
