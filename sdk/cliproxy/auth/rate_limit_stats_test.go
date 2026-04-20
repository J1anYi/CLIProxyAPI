package auth

import (
	"sync"
	"testing"
	"time"
)

func TestRateLimitStats_IncrementModelArts81101(t *testing.T) {
	stats := NewRateLimitStats()

	// Initial state
	date, modelArts, decodeServer := stats.GetStats()
	if modelArts != 0 {
		t.Errorf("initial modelArts count = %d, want 0", modelArts)
	}

	// Increment
	stats.IncrementModelArts81101()
	date, modelArts, decodeServer = stats.GetStats()
	if modelArts != 1 {
		t.Errorf("after increment, modelArts count = %d, want 1", modelArts)
	}
	if decodeServer != 0 {
		t.Errorf("decodeServer count = %d, want 0", decodeServer)
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
	_, modelArts, decodeServer := stats.GetStats()
	if decodeServer != 0 {
		t.Errorf("initial decodeServer count = %d, want 0", decodeServer)
	}

	// Increment
	stats.IncrementDecodeServerOverloaded()
	_, modelArts, decodeServer = stats.GetStats()
	if decodeServer != 1 {
		t.Errorf("after increment, decodeServer count = %d, want 1", decodeServer)
	}
	if modelArts != 0 {
		t.Errorf("modelArts count = %d, want 0", modelArts)
	}
}

func TestRateLimitStats_BothCounters(t *testing.T) {
	stats := NewRateLimitStats()

	// Increment both
	stats.IncrementModelArts81101()
	stats.IncrementModelArts81101()
	stats.IncrementDecodeServerOverloaded()

	_, modelArts, decodeServer := stats.GetStats()
	if modelArts != 2 {
		t.Errorf("modelArts count = %d, want 2", modelArts)
	}
	if decodeServer != 1 {
		t.Errorf("decodeServer count = %d, want 1", decodeServer)
	}
}

func TestRateLimitStats_DateRollover(t *testing.T) {
	stats := &RateLimitStats{
		date:                  "2020-01-01", // Old date
		modelArts81101:        100,
		decodeServerOverloaded: 50,
	}

	// GetStats should trigger rollover
	date, modelArts, decodeServer := stats.GetStats()

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
}

func TestRateLimitStats_ConcurrentIncrements(t *testing.T) {
	stats := NewRateLimitStats()

	// Run concurrent increments
	var wg sync.WaitGroup
	numGoroutines := 100
	incrementsPerGoroutine := 100

	// Half increment ModelArts, half increment DecodeServer
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				if idx%2 == 0 {
					stats.IncrementModelArts81101()
				} else {
					stats.IncrementDecodeServerOverloaded()
				}
			}
		}(i)
	}

	wg.Wait()

	_, modelArts, decodeServer := stats.GetStats()
	expectedPerCounter := int64(numGoroutines / 2 * incrementsPerGoroutine)

	if modelArts != expectedPerCounter {
		t.Errorf("modelArts count = %d, want %d", modelArts, expectedPerCounter)
	}
	if decodeServer != expectedPerCounter {
		t.Errorf("decodeServer count = %d, want %d", decodeServer, expectedPerCounter)
	}
}
