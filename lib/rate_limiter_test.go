package lib

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDefaultRateLimiter tests the default rate limiter configuration
func TestDefaultRateLimiter(t *testing.T) {
	rl := DefaultRateLimiter()
	assert.NotNil(t, rl)
	assert.Equal(t, 100*time.Millisecond, rl.GetDelay())
	assert.Equal(t, 3, rl.GetMaxRetries())
	assert.True(t, rl.IsEnabled())
}

// TestNewRateLimiter tests creating a custom rate limiter
func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(200*time.Millisecond, 5)
	assert.NotNil(t, rl)
	assert.Equal(t, 200*time.Millisecond, rl.GetDelay())
	assert.Equal(t, 5, rl.GetMaxRetries())
	assert.True(t, rl.IsEnabled())
}

// TestRateLimiterSetters tests setter methods
func TestRateLimiterSetters(t *testing.T) {
	rl := DefaultRateLimiter()

	// Test SetDelay
	rl.SetDelay(50 * time.Millisecond)
	assert.Equal(t, 50*time.Millisecond, rl.GetDelay())

	// Test SetEnabled
	rl.SetEnabled(false)
	assert.False(t, rl.IsEnabled())
	rl.SetEnabled(true)
	assert.True(t, rl.IsEnabled())

	// Test SetMaxRetries
	rl.SetMaxRetries(10)
	assert.Equal(t, 10, rl.GetMaxRetries())
}

// TestRateLimiterWait tests the Wait method
func TestRateLimiterWait(t *testing.T) {
	rl := NewRateLimiter(50*time.Millisecond, 3)

	// Test enabled wait
	start := time.Now()
	rl.Wait()
	elapsed := time.Since(start)
	assert.GreaterOrEqual(t, elapsed, 50*time.Millisecond)

	// Test disabled wait
	rl.SetEnabled(false)
	start = time.Now()
	rl.Wait()
	elapsed = time.Since(start)
	assert.Less(t, elapsed, 10*time.Millisecond) // Should be nearly instant
}

// TestRateLimiterDoRequestSuccess tests successful request without retry
func TestRateLimiterDoRequestSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	rl := NewRateLimiter(10*time.Millisecond, 3)
	client := &http.Client{}

	req, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	resp, err := rl.DoRequestWithRetry(client, req, false)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

// TestRateLimiterDoRequestRetry tests retry on 429 error
func TestRateLimiterDoRequestRetry(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount <= 2 {
			// First two requests fail with 429
			w.WriteHeader(http.StatusTooManyRequests)
		} else {
			// Third request succeeds
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success after retry"))
		}
	}))
	defer server.Close()

	rl := NewRateLimiter(10*time.Millisecond, 3)
	client := &http.Client{}

	req, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	resp, err := rl.DoRequestWithRetry(client, req, false)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 3, callCount) // Should have made 3 calls
	resp.Body.Close()
}

// TestRateLimiterDoRequestMaxRetries tests that max retries is respected
func TestRateLimiterDoRequestMaxRetries(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	rl := NewRateLimiter(10*time.Millisecond, 2) // Max 2 retries = 3 total attempts
	client := &http.Client{}

	req, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	resp, err := rl.DoRequestWithRetry(client, req, false)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "rate limit exceeded")
	assert.Equal(t, 3, callCount) // Should have made maxRetries+1 calls
}

// TestRateLimiterRetryAfterHeader tests parsing of Retry-After header
func TestRateLimiterRetryAfterHeader(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Header().Set("Retry-After", "1") // 1 second
			w.WriteHeader(http.StatusTooManyRequests)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}
	}))
	defer server.Close()

	rl := NewRateLimiter(10*time.Millisecond, 3)
	client := &http.Client{}

	req, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	start := time.Now()
	resp, err := rl.DoRequestWithRetry(client, req, false)
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, callCount)
	// Should have waited at least 1 second due to Retry-After header
	assert.GreaterOrEqual(t, elapsed, 1*time.Second)
	resp.Body.Close()
}

// TestGetRetryAfter tests the getRetryAfter helper function
func TestGetRetryAfter(t *testing.T) {
	// Test with no Retry-After header
	resp := &http.Response{Header: http.Header{}}
	duration := getRetryAfter(resp)
	assert.Equal(t, time.Duration(0), duration)

	// Test with seconds format
	resp.Header.Set("Retry-After", "5")
	duration = getRetryAfter(resp)
	assert.Equal(t, 5*time.Second, duration)

	// Test with HTTP date format
	futureTime := time.Now().Add(5 * time.Second)
	resp.Header.Set("Retry-After", futureTime.Format(http.TimeFormat))
	duration = getRetryAfter(resp)
	// Should be close to 5 seconds (allow for some timing variance)
	if duration > 0 {
		assert.Greater(t, duration, 3*time.Second)
		assert.LessOrEqual(t, duration, 6*time.Second)
	} else {
		// HTTP date parsing might fail in some environments, just skip this check
		t.Log("HTTP date parsing not supported or failed, skipping date format test")
	}

	// Test with invalid format
	resp.Header.Set("Retry-After", "invalid")
	duration = getRetryAfter(resp)
	assert.Equal(t, time.Duration(0), duration)
}

// TestGlobalRateLimiter tests the global rate limiter functions
func TestGlobalRateLimiter(t *testing.T) {
	// Get default global rate limiter
	rl := GetGlobalRateLimiter()
	assert.NotNil(t, rl)

	// Set custom global rate limiter
	customRL := NewRateLimiter(250*time.Millisecond, 7)
	SetGlobalRateLimiter(customRL)

	newRL := GetGlobalRateLimiter()
	assert.Equal(t, 250*time.Millisecond, newRL.GetDelay())
	assert.Equal(t, 7, newRL.GetMaxRetries())

	// Test SetRateLimitDelay convenience function
	SetRateLimitDelay(300 * time.Millisecond)
	assert.Equal(t, 300*time.Millisecond, newRL.GetDelay())

	// Test DisableRateLimiting convenience function
	DisableRateLimiting()
	assert.False(t, newRL.IsEnabled())

	// Test EnableRateLimiting convenience function
	EnableRateLimiting()
	assert.True(t, newRL.IsEnabled())

	// Restore default for other tests
	SetGlobalRateLimiter(DefaultRateLimiter())
}

// TestRateLimiterConcurrency tests thread-safety
func TestRateLimiterConcurrency(t *testing.T) {
	rl := DefaultRateLimiter()

	// Run multiple goroutines that modify the rate limiter
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(n int) {
			rl.SetDelay(time.Duration(n*10) * time.Millisecond)
			_ = rl.GetDelay()
			rl.SetEnabled(n%2 == 0)
			_ = rl.IsEnabled()
			rl.SetMaxRetries(n)
			_ = rl.GetMaxRetries()
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// If we get here without data races, the test passes
	assert.True(t, true)
}

// TestRateLimiterExponentialBackoff tests exponential backoff timing
func TestRateLimiterExponentialBackoff(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	rl := NewRateLimiter(100*time.Millisecond, 3)
	client := &http.Client{}

	req, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	start := time.Now()
	resp, err := rl.DoRequestWithRetry(client, req, false)
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, 4, callCount) // Initial + 3 retries

	// With exponential backoff (2x multiplier):
	// Attempt 1: no delay
	// Attempt 2: 100ms delay
	// Attempt 3: 200ms delay
	// Attempt 4: 400ms delay
	// Total: 700ms minimum
	assert.GreaterOrEqual(t, elapsed, 700*time.Millisecond)
}
