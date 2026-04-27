package auth

import (
	"sync"
	"testing"
	"time"
)

func TestRateLimitStats_IncrementModelArts81101(t *testing.T) {
	stats := NewRateLimitStats()

	// Initial state
	date, modelArts81101, modelArts81011, decodeServer, error406, ctxLen, tcpTimeout, connectionReset, _ := stats.GetStats()
	if modelArts81101 != 0 {
		t.Errorf("initial modelArts81101 count = %d, want 0", modelArts81101)
	}

	// Increment
	stats.IncrementModelArts81101()
	date, modelArts81101, modelArts81011, decodeServer, error406, ctxLen, tcpTimeout, connectionReset, _ = stats.GetStats()
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
	if tcpTimeout != 0 {
		t.Errorf("tcpTimeout count = %d, want 0", tcpTimeout)
	}
	if connectionReset != 0 {
		t.Errorf("connectionReset count = %d, want 0", connectionReset)
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
	_, modelArts81101, modelArts81011, _, _, _, _, _, _ := stats.GetStats()
	if modelArts81011 != 0 {
		t.Errorf("initial modelArts81011 count = %d, want 0", modelArts81011)
	}

	// Increment
	stats.IncrementModelArts81011()
	_, modelArts81101, modelArts81011, _, _, _, _, _, _ = stats.GetStats()
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
	_, _, _, decodeServer, _, _, _, _, _ := stats.GetStats()
	if decodeServer != 0 {
		t.Errorf("initial decodeServer count = %d, want 0", decodeServer)
	}

	// Increment
	stats.IncrementDecodeServerOverloaded()
	_, modelArts81101, _, decodeServer, _, _, _, _, _ := stats.GetStats()
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
	_, modelArts81101, modelArts81011, decodeServer, error406, ctxLen, tcpTimeout, connectionReset, _ := stats.GetStats()
	if error406 != 1 {
		t.Errorf("after increment, error406 count = %d, want 1", error406)
	}
	if modelArts81101 != 0 || modelArts81011 != 0 || decodeServer != 0 || ctxLen != 0 || tcpTimeout != 0 || connectionReset != 0 {
		t.Errorf("other counters should be 0, got modelArts81101=%d, modelArts81011=%d, decodeServer=%d, ctxLen=%d, tcpTimeout=%d, connectionReset=%d", modelArts81101, modelArts81011, decodeServer, ctxLen, tcpTimeout, connectionReset)
	}
}

func TestRateLimitStats_IncrementContextLengthExceeded(t *testing.T) {
	stats := NewRateLimitStats()

	stats.IncrementContextLengthExceeded()
	_, modelArts81101, modelArts81011, decodeServer, error406, ctxLen, tcpTimeout, connectionReset, _ := stats.GetStats()
	if ctxLen != 1 {
		t.Errorf("after increment, ctxLen count = %d, want 1", ctxLen)
	}
	if modelArts81101 != 0 || modelArts81011 != 0 || decodeServer != 0 || error406 != 0 || tcpTimeout != 0 || connectionReset != 0 {
		t.Errorf("other counters should be 0, got modelArts81101=%d, modelArts81011=%d, decodeServer=%d, error406=%d, tcpTimeout=%d, connectionReset=%d", modelArts81101, modelArts81011, decodeServer, error406, tcpTimeout, connectionReset)
	}
}

func TestRateLimitStats_IncrementTCPTimeout(t *testing.T) {
	stats := NewRateLimitStats()

	stats.IncrementTCPTimeout()
	_, modelArts81101, modelArts81011, decodeServer, error406, ctxLen, tcpTimeout, connectionReset, _ := stats.GetStats()
	if tcpTimeout != 1 {
		t.Errorf("after increment, tcpTimeout count = %d, want 1", tcpTimeout)
	}
	if modelArts81101 != 0 || modelArts81011 != 0 || decodeServer != 0 || error406 != 0 || ctxLen != 0 || connectionReset != 0 {
		t.Errorf("other counters should be 0, got modelArts81101=%d, modelArts81011=%d, decodeServer=%d, error406=%d, ctxLen=%d, connectionReset=%d", modelArts81101, modelArts81011, decodeServer, error406, ctxLen, connectionReset)
	}
}

func TestRateLimitStats_IncrementConnectionReset(t *testing.T) {
	stats := NewRateLimitStats()

	stats.IncrementConnectionReset()
	_, modelArts81101, modelArts81011, decodeServer, error406, ctxLen, tcpTimeout, connectionReset, _ := stats.GetStats()
	if connectionReset != 1 {
		t.Errorf("after increment, connectionReset count = %d, want 1", connectionReset)
	}
	if modelArts81101 != 0 || modelArts81011 != 0 || decodeServer != 0 || error406 != 0 || ctxLen != 0 || tcpTimeout != 0 {
		t.Errorf("other counters should be 0, got modelArts81101=%d, modelArts81011=%d, decodeServer=%d, error406=%d, ctxLen=%d, tcpTimeout=%d", modelArts81101, modelArts81011, decodeServer, error406, ctxLen, tcpTimeout)
	}
}

func TestRateLimitStats_IncrementFailedResponse(t *testing.T) {
	stats := NewRateLimitStats()

	stats.IncrementFailedResponse()
	_, modelArts81101, modelArts81011, decodeServer, error406, ctxLen, tcpTimeout, connectionReset, failedResp := stats.GetStats()
	if failedResp != 1 {
		t.Errorf("after increment, failedResp count = %d, want 1", failedResp)
	}
	if modelArts81101 != 0 || modelArts81011 != 0 || decodeServer != 0 || error406 != 0 || ctxLen != 0 || tcpTimeout != 0 || connectionReset != 0 {
		t.Errorf("other counters should be 0, got modelArts81101=%d, modelArts81011=%d, decodeServer=%d, error406=%d, ctxLen=%d, tcpTimeout=%d, connectionReset=%d", modelArts81101, modelArts81011, decodeServer, error406, ctxLen, tcpTimeout, connectionReset)
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
	stats.IncrementTCPTimeout()
	stats.IncrementConnectionReset()
	stats.IncrementFailedResponse()

	_, modelArts81101, modelArts81011, decodeServer, error406, ctxLen, tcpTimeout, connectionReset, failedResp := stats.GetStats()
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
	if tcpTimeout != 1 {
		t.Errorf("tcpTimeout count = %d, want 1", tcpTimeout)
	}
	if connectionReset != 1 {
		t.Errorf("connectionReset count = %d, want 1", connectionReset)
	}
	if failedResp != 1 {
		t.Errorf("failedResp count = %d, want 1", failedResp)
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
		tcpTimeout:             5,
		connectionReset:        3,
		failedResponse:         2,
	}

	// GetStats should trigger rollover
	date, modelArts81101, modelArts81011, decodeServer, error406, ctxLen, tcpTimeout, connectionReset, failedResp := stats.GetStats()

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
	if tcpTimeout != 0 {
		t.Errorf("after rollover, tcpTimeout count = %d, want 0", tcpTimeout)
	}
	if connectionReset != 0 {
		t.Errorf("after rollover, connectionReset count = %d, want 0", connectionReset)
	}
	if failedResp != 0 {
		t.Errorf("after rollover, failedResp count = %d, want 0", failedResp)
	}
}

func TestRateLimitStats_ConcurrentIncrements(t *testing.T) {
	stats := NewRateLimitStats()

	// Run concurrent increments
	var wg sync.WaitGroup
	// Use 96 goroutines (divisible by 8) for even distribution
	numGoroutines := 96
	incrementsPerGoroutine := 100

	// Distribute across all counter types (now 8 types)
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				switch idx % 8 {
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
				case 5:
					stats.IncrementTCPTimeout()
				case 6:
					stats.IncrementConnectionReset()
				case 7:
					stats.IncrementFailedResponse()
				}
			}
		}(i)
	}

	wg.Wait()

	_, modelArts81101, modelArts81011, decodeServer, error406, ctxLen, tcpTimeout, connectionReset, failedResp := stats.GetStats()
	// 96 goroutines / 8 types = 12 goroutines per type, each doing 100 increments
	expectedPerCounter := int64(numGoroutines / 8 * incrementsPerGoroutine)

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
	if tcpTimeout != expectedPerCounter {
		t.Errorf("tcpTimeout count = %d, want %d", tcpTimeout, expectedPerCounter)
	}
	if connectionReset != expectedPerCounter {
		t.Errorf("connectionReset count = %d, want %d", connectionReset, expectedPerCounter)
	}
	if failedResp != expectedPerCounter {
		t.Errorf("failedResp count = %d, want %d", failedResp, expectedPerCounter)
	}
}
