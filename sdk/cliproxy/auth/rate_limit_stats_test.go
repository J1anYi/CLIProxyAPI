package auth

import (
	"sync"
	"testing"
	"time"
)

func TestRateLimitStats_IncrementModelArts81101(t *testing.T) {
	stats := NewRateLimitStats()

	// Initial state
	date, modelArts, decodeServer, error406, ctxLen := stats.GetStats()
	if modelArts != 0 {
		t.Errorf("initial modelArts count = %d, want 0", modelArts)
	}

	// Increment
	stats.IncrementModelArts81101()
	date, modelArts, decodeServer, error406, ctxLen = stats.GetStats()
	if modelArts != 1 {
		t.Errorf("after increment, modelArts count = %d, want 1", modelArts)
	}
	if decodeServer != 0 {
		t.Errorf("decodeServer count = %d, want 0", decodeServer)
	}
	if error406 != 0 {
		t.Errorf("error406 count = %d, want 0", error406)
	}
	if ctxLen != 0 {
		t.Errorf("ctxLen count = %d, want 0", ctxLen)
	}

	// Verify date is today
	expectedDate := time.Now().Format("2006-01-02")
	if date != expectedDate {
		t.Errorf("date = %s, want %s", date, expectedDate)
	}
}

func TestRateLimitStats_IncrementDecodeServerOverloaded(t *testing.T) {
	stats := NewRateLimitStats()

	// Initial state
	_, modelArts, decodeServer, _, _ := stats.GetStats()
	if decodeServer != 0 {
		t.Errorf("initial decodeServer count = %d, want 0", decodeServer)
	}

	// Increment
	stats.IncrementDecodeServerOverloaded()
	_, modelArts, decodeServer, _, _ = stats.GetStats()
	if decodeServer != 1 {
		t.Errorf("after increment, decodeServer count = %d, want 1", decodeServer)
	}
	if modelArts != 0 {
		t.Errorf("modelArts count = %d, want 0", modelArts)
	}
}

func TestRateLimitStats_IncrementError406(t *testing.T) {
	stats := NewRateLimitStats()

	stats.IncrementError406()
	_, modelArts, decodeServer, error406, ctxLen := stats.GetStats()
	if error406 != 1 {
		t.Errorf("after increment, error406 count = %d, want 1", error406)
	}
	if modelArts != 0 || decodeServer != 0 || ctxLen != 0 {
		t.Errorf("other counters should be 0, got modelArts=%d, decodeServer=%d, ctxLen=%d", modelArts, decodeServer, ctxLen)
	}
}

func TestRateLimitStats_IncrementContextLengthExceeded(t *testing.T) {
	stats := NewRateLimitStats()

	stats.IncrementContextLengthExceeded()
	_, modelArts, decodeServer, error406, ctxLen := stats.GetStats()
	if ctxLen != 1 {
		t.Errorf("after increment, ctxLen count = %d, want 1", ctxLen)
	}
	if modelArts != 0 || decodeServer != 0 || error406 != 0 {
		t.Errorf("other counters should be 0, got modelArts=%d, decodeServer=%d, error406=%d", modelArts, decodeServer, error406)
	}
}

func TestRateLimitStats_BothCounters(t *testing.T) {
	stats := NewRateLimitStats()

	// Increment multiple
	stats.IncrementModelArts81101()
	stats.IncrementModelArts81101()
	stats.IncrementDecodeServerOverloaded()
	stats.IncrementError406()
	stats.IncrementContextLengthExceeded()

	_, modelArts, decodeServer, error406, ctxLen := stats.GetStats()
	if modelArts != 2 {
		t.Errorf("modelArts count = %d, want 2", modelArts)
	}
	if decodeServer != 1 {
		t.Errorf("decodeServer count = %d, want 1", decodeServer)
	}
	if error406 != 1 {
		t.Errorf("error406 count = %d, want 1", error406)
	}
	if ctxLen != 1 {
		t.Errorf("ctxLen count = %d, want 1", ctxLen)
	}
}

func TestRateLimitStats_DateRollover(t *testing.T) {
	stats := &RateLimitStats{
		date:                   "2020-01-01", // Old date
		modelArts81101:         100,
		decodeServerOverloaded: 50,
		error406:               30,
		contextLengthExceeded:  20,
	}

	// GetStats should trigger rollover
	date, modelArts, decodeServer, error406, ctxLen := stats.GetStats()

	expectedDate := time.Now().Format("2006-01-02")
	if date != expectedDate {
		t.Errorf("date = %s, want %s", date, expectedDate)
	}
	if modelArts != 0 {
		t.Errorf("after rollover, modelArts count = %d, want 0", modelArts)
	}
	if decodeServer != 0 {
		t.Errorf("after rollover, decodeServer count = %d, want 0", decodeServer)
	}
	if error406 != 0 {
		t.Errorf("after rollover, error406 count = %d, want 0", error406)
	}
	if ctxLen != 0 {
		t.Errorf("after rollover, ctxLen count = %d, want 0", ctxLen)
	}
}

func TestRateLimitStats_ConcurrentIncrements(t *testing.T) {
	stats := NewRateLimitStats()

	// Run concurrent increments
	var wg sync.WaitGroup
	numGoroutines := 100
	incrementsPerGoroutine := 100

	// Distribute across all counter types
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				switch idx % 4 {
				case 0:
					stats.IncrementModelArts81101()
				case 1:
					stats.IncrementDecodeServerOverloaded()
				case 2:
					stats.IncrementError406()
				case 3:
					stats.IncrementContextLengthExceeded()
				}
			}
		}(i)
	}

	wg.Wait()

	_, modelArts, decodeServer, error406, ctxLen := stats.GetStats()
	expectedPerCounter := int64(numGoroutines / 4 * incrementsPerGoroutine)

	if modelArts != expectedPerCounter {
		t.Errorf("modelArts count = %d, want %d", modelArts, expectedPerCounter)
	}
	if decodeServer != expectedPerCounter {
		t.Errorf("decodeServer count = %d, want %d", decodeServer, expectedPerCounter)
	}
	if error406 != expectedPerCounter {
		t.Errorf("error406 count = %d, want %d", error406, expectedPerCounter)
	}
	if ctxLen != expectedPerCounter {
		t.Errorf("ctxLen count = %d, want %d", ctxLen, expectedPerCounter)
	}
}
