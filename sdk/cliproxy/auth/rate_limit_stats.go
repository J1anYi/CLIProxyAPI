package auth

import (
	"sync"
	"time"
)

// RateLimitStats tracks daily statistics for rate limit errors.
// It uses a single mutex for thread-safe operations including date rollover.
type RateLimitStats struct {
	mu                    sync.Mutex
	date                  string // YYYY-MM-DD format
	modelArts81101        int64
	decodeServerOverloaded int64
}

// NewRateLimitStats creates a new stats tracker initialized with today's date.
func NewRateLimitStats() *RateLimitStats {
	return &RateLimitStats{
		date: time.Now().Format("2006-01-02"),
	}
}

// IncrementModelArts81101 increments the counter for ModelArts.81101 errors.
// If the date has changed, counters are reset before incrementing.
func (s *RateLimitStats) IncrementModelArts81101() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkDateRollover()
	s.modelArts81101++
}

// IncrementDecodeServerOverloaded increments the counter for Decode server overloaded errors.
// If the date has changed, counters are reset before incrementing.
func (s *RateLimitStats) IncrementDecodeServerOverloaded() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkDateRollover()
	s.decodeServerOverloaded++
}

// GetStats returns the current date and counters.
// If the date has changed, counters are reset before returning.
func (s *RateLimitStats) GetStats() (date string, modelArts81101 int64, decodeServerOverloaded int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkDateRollover()
	return s.date, s.modelArts81101, s.decodeServerOverloaded
}

// checkDateRollover resets counters if the date has changed.
// Must be called with mutex already held.
func (s *RateLimitStats) checkDateRollover() {
	today := time.Now().Format("2006-01-02")
	if s.date != today {
		s.date = today
		s.modelArts81101 = 0
		s.decodeServerOverloaded = 0
	}
}
