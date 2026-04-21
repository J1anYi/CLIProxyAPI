package auth

import (
	"sync"
	"testing"
	"time"
)

func TestRateLimitStats_IncrementModelArts81101(t *testing.T) {
	stats := NewRateLimitStats()

	// Initial state
	date, modelArts81101, modelArts81011, decodeServer, error406, ctxLen := stats.GetStats()
	if modelArts81101 != 0 {
		t.Errorf("initial modelArts81101 count = %d, want 0", modelArts81101)
	}

	// Increment
	stats.IncrementModelArts81101()
	date, modelArts81101, modelArts81011, decodeServer, error406, ctxLen = stats.GetStats()
	if modelArts81101 != 1 {
		t.Errorf("after increment, modelArts81101 count = %d, want 1", modelArts81101)
	}
	if modelArts81011 != 0 {
		t.Errorf("modelArts81011 count = %d, want 0", modelArts81011)
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

func TestRateLimitStats_IncrementModelArts81011(t *testing.T) {
	stats := NewRateLimitStats()

	// Initial state
	_, modelArts81101, modelArts81011, _, _, _ := stats.GetStats()
	if modelArts81011 != 0 {
		t.Errorf("initial modelArts81011 count = %d, want 0", modelArts81011)
	}

	// Increment
	stats.IncrementModelArts81011()
	_, modelArts81101, modelArts81011, _, _, _ = stats.GetStats()
	if modelArts81011 != 1 {
		t.Errorf("after increment, modelArts81011 count = %d, want 1", modelArts81011)
	}
	if modelArts81101 != 0 {
		t.Errorf("modelArts81101 count = %d, want 0", modelArts81101)
	}
}

func TestRateLimitStats_IncrementDecodeServerOverloaded(t *testing.T) {
	stats := NewRateLimitStats()

	// Initial state
	_, _, _, decodeServer, _, _ := stats.GetStats()
	if decodeServer != 0 {
		t.Errorf("initial decodeServer count = %d, want 0", decodeServer)
	}

	// Increment
	stats.IncrementDecodeServerOverloaded()
	_, modelArts81101, _, decodeServer, _, _ := stats.GetStats()
	if decodeServer != 1 {
		t.Errorf("after increment, decodeServer count = %d, want 1", decodeServer)
	}
	if modelArts81101 != 0 {
		t.Errorf("modelArts81101 count = %d, want 0", modelArts81101)
	}
}

func TestRateLimitStats_IncrementError406(t *testing.T) {
	stats := NewRateLimitStats()

	stats.IncrementError406()
	_, modelArts81101, modelArts81011, decodeServer, error406, ctxLen := stats.GetStats()
	if error406 != 1 {
		t.Errorf("after increment, error406 count = %d, want 1", error406)
	}
	if modelArts81101 != 0 || modelArts81011 != 0 || decodeServer != 0 || ctxLen != 0 {
		t.Errorf("other counters should be 0, got modelArts81101=%d, modelArts81011=%d, decodeServer=%d, ctxLen=%d", modelArts81101, modelArts81011, decodeServer, ctxLen)
	}
}

func TestRateLimitStats_IncrementContextLengthExceeded(t *testing.T) {
	stats := NewRateLimitStats()

	stats.IncrementContextLengthExceeded()
	_, modelArts81101, modelArts81011, decodeServer, error406, ctxLen := stats.GetStats()
	if ctxLen != 1 {
		t.Errorf("after increment, ctxLen count = %d, want 1", ctxLen)
	}
	if modelArts81101 != 0 || modelArts81011 != 0 || decodeServer != 0 || error406 != 0 {
		t.Errorf("other counters should be 0, got modelArts81101=%d, modelArts81011=%d, decodeServer=%d, error406=%d", modelArts81101, modelArts81011, decodeServer, error406)
	}
}

func TestRateLimitStats_AllCounters(t *testing.T) {
	stats := NewRateLimitStats()

	// Increment multiple
	stats.IncrementModelArts81101()
	stats.IncrementModelArts81101()
	stats.IncrementModelArts81011()
	stats.IncrementModelArts81011()
	stats.IncrementModelArts81011()
	stats.IncrementDecodeServerOverloaded()
	stats.IncrementError406()
	stats.IncrementContextLengthExceeded()

	_, modelArts81101, modelArts81011, decodeServer, error406, ctxLen := stats.GetStats()
	if modelArts81101 != 2 {
		t.Errorf("modelArts81101 count = %d, want 2", modelArts81101)
	}
	if modelArts81011 != 3 {
		t.Errorf("modelArts81011 count = %d, want 3", modelArts81011)
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
		modelArts81011:         50,
		decodeServerOverloaded: 30,
		error406:               20,
		contextLengthExceeded:  10,
	}

	// GetStats should trigger rollover
	date, modelArts81101, modelArts81011, decodeServer, error406, ctxLen := stats.GetStats()

	expectedDate := time.Now().Format("2006-01-02")
	if date != expectedDate {
		t.Errorf("date = %s, want %s", date, expectedDate)
	}
	if modelArts81101 != 0 {
		t.Errorf("after rollover, modelArts81101 count = %d, want 0", modelArts81101)
	}
	if modelArts81011 != 0 {
		t.Errorf("after rollover, modelArts81011 count = %d, want 0", modelArts81011)
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
				switch idx % 5 {
				case 0:
					stats.IncrementModelArts81101()
				case 1:
					stats.IncrementModelArts81011()
				case 2:
					stats.IncrementDecodeServerOverloaded()
				case 3:
					stats.IncrementError406()
				case 4:
					stats.IncrementContextLengthExceeded()
				}
			}
		}(i)
	}

	wg.Wait()

	_, modelArts81101, modelArts81011, decodeServer, error406, ctxLen := stats.GetStats()
	expectedPerCounter := int64(numGoroutines / 5 * incrementsPerGoroutine)

	if modelArts81101 != expectedPerCounter {
		t.Errorf("modelArts81101 count = %d, want %d", modelArts81101, expectedPerCounter)
	}
	if modelArts81011 != expectedPerCounter {
		t.Errorf("modelArts81011 count = %d, want %d", modelArts81011, expectedPerCounter)
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
