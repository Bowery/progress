// Copyright 2015 Bowery, Inc.

package progress

// Status represents the copy progress. It contains the current progress
// and total size in bytes.
type Status struct {
	Current  int64
	Total    int64
	finished bool
}

// Completion returns the current completion in the range [0, 1].
func (s *Status) Completion() float64 {
	return float64(s.Current) / float64(s.Total)
}

// IsFinished returns a boolean indicating whether the copy
// has completed.
func (s *Status) IsFinished() bool {
	return s.finished
}
