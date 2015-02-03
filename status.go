// Copyright 2015 Bowery, Inc.

package progress

// Status represents the download/upload progress.
// It is comprised of the current progress (bytes)
// and the total size (bytes)
type Status struct {
	Current int64
	Total   int64
}

// Completion returns the current completion in the range [0, 1].
func (s *Status) Completion() float64 {
	return float64(s.Current) / float64(s.Total)
}

// IsFinished returns a boolean indicating whether the
// requests status is complete.
func (s *Status) IsFinished() bool {
	return s.Current == s.Total
}
