package lib

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimiter manages API request rate limiting with retry logic
type RateLimiter struct {
	delay           time.Duration
	maxRetries      int
	backoffMultiple float64
	enabled         bool
	mu              sync.RWMutex
}

// DefaultRateLimiter returns a rate limiter with sensible defaults
// - 100ms delay between requests
// - 3 retries on rate limit errors
// - 2x backoff multiplier
func DefaultRateLimiter() *RateLimiter {
	return &RateLimiter{
		delay:           100 * time.Millisecond,
		maxRetries:      3,
		backoffMultiple: 2.0,
		enabled:         true,
	}
}

// NewRateLimiter creates a custom rate limiter
func NewRateLimiter(delay time.Duration, maxRetries int) *RateLimiter {
	return &RateLimiter{
		delay:           delay,
		maxRetries:      maxRetries,
		backoffMultiple: 2.0,
		enabled:         true,
	}
}

// Wait applies the rate limiting delay
func (rl *RateLimiter) Wait() {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	if rl.enabled && rl.delay > 0 {
		time.Sleep(rl.delay)
	}
}

// SetDelay updates the delay duration
func (rl *RateLimiter) SetDelay(delay time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.delay = delay
}

// GetDelay returns the current delay
func (rl *RateLimiter) GetDelay() time.Duration {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.delay
}

// SetEnabled enables or disables rate limiting
func (rl *RateLimiter) SetEnabled(enabled bool) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.enabled = enabled
}

// IsEnabled returns whether rate limiting is enabled
func (rl *RateLimiter) IsEnabled() bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.enabled
}

// SetMaxRetries sets the maximum number of retries for 429 errors
func (rl *RateLimiter) SetMaxRetries(maxRetries int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.maxRetries = maxRetries
}

// GetMaxRetries returns the maximum number of retries
func (rl *RateLimiter) GetMaxRetries() int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.maxRetries
}

// DoRequestWithRetry executes an HTTP request with automatic retry on rate limit errors (429)
// It implements exponential backoff and respects the Retry-After header if present
func (rl *RateLimiter) DoRequestWithRetry(httpClient *http.Client, req *http.Request, verbose bool) (*http.Response, error) {
	rl.mu.RLock()
	maxRetries := rl.maxRetries
	baseDelay := rl.delay
	backoffMultiple := rl.backoffMultiple
	rl.mu.RUnlock()

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err = httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("HTTP request failed: %w", err)
		}

		// Success - return response
		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		// Rate limited (429)
		resp.Body.Close()

		// If this was the last attempt, return error
		if attempt >= maxRetries {
			return nil, fmt.Errorf("rate limit exceeded after %d retries", maxRetries)
		}

		// Calculate backoff delay
		var backoffDelay time.Duration

		// Try to get Retry-After header
		retryAfter := getRetryAfter(resp)
		if retryAfter > 0 {
			backoffDelay = retryAfter
			if verbose {
				fmt.Printf("Rate limited (429), retrying after %v (from Retry-After header, attempt %d/%d)...\n",
					backoffDelay, attempt+1, maxRetries)
			}
		} else {
			// Use exponential backoff
			backoffDelay = time.Duration(float64(baseDelay) * math.Pow(backoffMultiple, float64(attempt)))
			if verbose {
				fmt.Printf("Rate limited (429), retrying in %v with exponential backoff (attempt %d/%d)...\n",
					backoffDelay, attempt+1, maxRetries)
			}
		}

		time.Sleep(backoffDelay)
	}

	return nil, fmt.Errorf("rate limit exceeded after %d retries", maxRetries)
}

// getRetryAfter parses the Retry-After header from an HTTP response
// It supports both delay-seconds and HTTP-date formats
func getRetryAfter(resp *http.Response) time.Duration {
	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter == "" {
		return 0
	}

	// Try parsing as seconds
	if seconds, err := strconv.Atoi(retryAfter); err == nil {
		return time.Duration(seconds) * time.Second
	}

	// Try parsing as HTTP date
	if t, err := http.ParseTime(retryAfter); err == nil {
		duration := time.Until(t)
		if duration > 0 {
			return duration
		}
	}

	return 0
}

// Global rate limiter instance
var (
	globalRateLimiter = DefaultRateLimiter()
	globalRLMutex     sync.RWMutex
)

// SetGlobalRateLimiter sets the global rate limiter instance
func SetGlobalRateLimiter(rl *RateLimiter) {
	globalRLMutex.Lock()
	defer globalRLMutex.Unlock()
	globalRateLimiter = rl
}

// GetGlobalRateLimiter returns the global rate limiter instance
func GetGlobalRateLimiter() *RateLimiter {
	globalRLMutex.RLock()
	defer globalRLMutex.RUnlock()
	return globalRateLimiter
}

// SetRateLimitDelay is a convenience function to set the global rate limit delay
func SetRateLimitDelay(delay time.Duration) {
	globalRLMutex.RLock()
	rl := globalRateLimiter
	globalRLMutex.RUnlock()
	rl.SetDelay(delay)
}

// DisableRateLimiting is a convenience function to disable global rate limiting (useful for testing)
func DisableRateLimiting() {
	globalRLMutex.RLock()
	rl := globalRateLimiter
	globalRLMutex.RUnlock()
	rl.SetEnabled(false)
}

// EnableRateLimiting is a convenience function to enable global rate limiting
func EnableRateLimiting() {
	globalRLMutex.RLock()
	rl := globalRateLimiter
	globalRLMutex.RUnlock()
	rl.SetEnabled(true)
}
