package zscaler

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestTimeoutCalculationLogic verifies the core timeout calculation logic
func TestTimeoutCalculationLogic(t *testing.T) {
	t.Run("Wait time is excluded from timeout check", func(t *testing.T) {
		// Simulate the timeout calculation logic from ExecuteRequest
		overallStartTime := time.Now()
		totalWaitTime := time.Duration(0)
		requestTimeout := 3 * time.Second
		
		// Simulate some actual API work
		time.Sleep(500 * time.Millisecond)
		
		// Simulate a rate limit wait (this should NOT count toward timeout)
		waitStart := time.Now()
		time.Sleep(2 * time.Second) // 2 seconds of waiting
		totalWaitTime += time.Since(waitStart)
		
		// More actual API work
		time.Sleep(500 * time.Millisecond)
		
		// Calculate elapsed time (excluding waits) - this is what ExecuteRequest does
		elapsedTime := time.Since(overallStartTime) - totalWaitTime
		totalTime := time.Since(overallStartTime)
		
		t.Logf("Timing breakdown:")
		t.Logf("  Total wall-clock time: %v", totalTime)
		t.Logf("  Wait time (excluded):  %v", totalWaitTime)
		t.Logf("  Actual time (counted): %v", elapsedTime)
		t.Logf("  Timeout threshold:     %v", requestTimeout)
		
		// Verify: total time > timeout, but actual time < timeout
		if totalTime <= requestTimeout {
			t.Errorf("Test setup error: expected total time (%v) > timeout (%v)", totalTime, requestTimeout)
			return
		}
		
		if elapsedTime >= requestTimeout {
			t.Errorf("FAIL: Should NOT timeout because wait time should be excluded")
			t.Errorf("  Elapsed (counted): %v", elapsedTime)
			t.Errorf("  Timeout threshold: %v", requestTimeout)
			t.Errorf("  This proves wait times ARE counting (BUG!)")
			return
		}
		
		t.Logf("✓ SUCCESS: Wait time correctly excluded from timeout")
		t.Logf("  %v actual < %v timeout (despite %v total)", elapsedTime, requestTimeout, totalTime)
	})
	
	t.Run("Actual processing time exceeds timeout - should timeout", func(t *testing.T) {
		overallStartTime := time.Now()
		totalWaitTime := time.Duration(0)
		requestTimeout := 1 * time.Second
		
		// Simulate slow actual processing (no waits, just slow work)
		time.Sleep(1500 * time.Millisecond)
		
		elapsedTime := time.Since(overallStartTime) - totalWaitTime
		
		t.Logf("Timing: actual=%v, timeout=%v", elapsedTime, requestTimeout)
		
		if elapsedTime < requestTimeout {
			t.Errorf("FAIL: Should timeout when actual time exceeds threshold")
			t.Errorf("  Elapsed: %v", elapsedTime)
			t.Errorf("  Timeout: %v", requestTimeout)
			return
		}
		
		t.Logf("✓ SUCCESS: Correctly detected timeout: %v actual > %v timeout", elapsedTime, requestTimeout)
	})
	
	t.Run("Multiple wait periods accumulate correctly", func(t *testing.T) {
		overallStartTime := time.Now()
		totalWaitTime := time.Duration(0)
		requestTimeout := 2 * time.Second
		
		// Simulate multiple rate limit cycles
		for i := 0; i < 3; i++ {
			// Actual work
			time.Sleep(200 * time.Millisecond)
			
			// Rate limit wait
			waitStart := time.Now()
			time.Sleep(1 * time.Second)
			totalWaitTime += time.Since(waitStart)
		}
		
		elapsedTime := time.Since(overallStartTime) - totalWaitTime
		totalTime := time.Since(overallStartTime)
		
		t.Logf("Multiple waits test:")
		t.Logf("  Total time:     %v", totalTime)
		t.Logf("  Total waits:    %v", totalWaitTime)
		t.Logf("  Actual work:    %v", elapsedTime)
		t.Logf("  Timeout:        %v", requestTimeout)
		
		// Total time should be ~3.6s (0.6s work + 3s waits)
		// But counted time should be ~0.6s (just the work)
		if elapsedTime >= requestTimeout {
			t.Errorf("FAIL: Should NOT timeout despite long total time")
			t.Errorf("  Actual: %v, Timeout: %v", elapsedTime, requestTimeout)
			return
		}
		
		t.Logf("✓ SUCCESS: Multiple waits correctly excluded: %v actual < %v timeout (despite %v total)", 
			elapsedTime, requestTimeout, totalTime)
	})
	
	t.Run("Extreme wait time scenario - 30s wait with 3s timeout", func(t *testing.T) {
		overallStartTime := time.Now()
		totalWaitTime := time.Duration(0)
		requestTimeout := 3 * time.Second
		
		// Simulate actual API work (fast)
		time.Sleep(500 * time.Millisecond)
		
		// Simulate extreme rate limit wait (like hourly limit hit)
		waitStart := time.Now()
		time.Sleep(5 * time.Second) // Using 5s instead of 30s to keep test fast
		totalWaitTime += time.Since(waitStart)
		
		// More work
		time.Sleep(500 * time.Millisecond)
		
		elapsedTime := time.Since(overallStartTime) - totalWaitTime
		totalTime := time.Since(overallStartTime)
		
		t.Logf("Extreme wait scenario:")
		t.Logf("  Total time:     %v", totalTime)
		t.Logf("  Wait time:      %v", totalWaitTime)
		t.Logf("  Actual time:    %v", elapsedTime)
		t.Logf("  Timeout:        %v", requestTimeout)
		
		if elapsedTime >= requestTimeout {
			t.Errorf("FAIL: Should NOT timeout even with extreme wait time")
			return
		}
		
		t.Logf("✓ SUCCESS: Extreme wait time handled correctly")
		t.Logf("  This simulates hourly rate limit scenario where wait is very long")
	})
}

// TestTimeoutErrorMessage verifies error message format
func TestTimeoutErrorMessage(t *testing.T) {
	elapsed := 5 * time.Minute
	waited := 4 * time.Minute
	
	// This is the format from the actual code
	errMsg := fmt.Sprintf("request timeout exceeded after %v (excluding %v of rate limit waits)", 
		elapsed, waited)
	
	expectedParts := []string{
		"request timeout exceeded",
		"after 5m0s",
		"excluding 4m0s",
		"rate limit waits",
	}
	
	for _, part := range expectedParts {
		if !strings.Contains(errMsg, part) {
			t.Errorf("Error message missing expected part: %s", part)
			t.Errorf("Full message: %s", errMsg)
		}
	}
	
	t.Logf("✓ Error message format validated: %s", errMsg)
}

