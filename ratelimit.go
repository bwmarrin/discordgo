package discordgo

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Ratelimiter holds all ratelimit buckets
type RateLimiter struct {
	sync.Mutex
	global          *Bucket
	buckets         map[string]*Bucket
	globalRateLimit time.Duration
}

// New returns a new RateLimiter
func NewRatelimiter() *RateLimiter {

	return &RateLimiter{
		buckets: make(map[string]*Bucket),
		global:  &Bucket{Key: "global"},
	}
}

// getBucket retrieves or creates a bucket
func (r *RateLimiter) getBucket(key string) *Bucket {
	r.Lock()
	defer r.Unlock()

	if bucket, ok := r.buckets[key]; ok {
		return bucket
	}

	b := &Bucket{
		remaining: 1,
		Key:       key,
		global:    r.global,
	}

	r.buckets[key] = b
	return b
}

// LockBucket Locks until a request can be made
func (r *RateLimiter) LockBucket(bucketID string) *Bucket {

	b := r.getBucket(bucketID)

	b.mu.Lock()

	// If we ran out of calls and the reset time is still ahead of us
	// then we need to take it easy and relax a little
	if b.remaining < 1 && b.reset.After(time.Now()) {
		time.Sleep(b.reset.Sub(time.Now()))

	}

	// Check for global ratelimits
	r.global.mu.Lock()
	r.global.mu.Unlock()

	b.remaining--
	return b
}

// Bucket represents a ratelimit bucket, each bucket gets ratelimited individually (-global ratelimits)
type Bucket struct {
	Key string

	mu        sync.Mutex
	remaining int
	limit     int
	reset     time.Time
	global    *Bucket
}

// Release unlocks the bucket and reads the headers to update the buckets ratelimit info
// and locks up the whole thing in case if there's a global ratelimit.
func (b *Bucket) Release(headers http.Header) error {

	defer b.mu.Unlock()
	if headers == nil {
		return nil
	}

	remaining := headers.Get("X-RateLimit-Remaining")
	reset := headers.Get("X-RateLimit-Reset")
	global := headers.Get("X-RateLimit-Global")
	retryAfter := headers.Get("Retry-After")

	// If it's global just keep the main ratelimit mutex locked
	if global != "" {
		parsedAfter, err := strconv.Atoi(retryAfter)
		if err != nil {
			return err
		}

		// Lock it in a new goroutine so that this isn't a blocking call
		go func() {
			// Make sure if several requests were waiting we don't sleep for n * retry-after
			// where n is the amount of requests that were going on
			sleepTo := time.Now().Add(time.Duration(parsedAfter) * time.Millisecond)

			b.global.mu.Lock()

			sleepDuration := sleepTo.Sub(time.Now())
			if sleepDuration > 0 {
				time.Sleep(sleepDuration)
			}

			b.global.mu.Unlock()
		}()

		return nil
	}

	// Update reset time if either retry after or reset headers are present
	// Prefer retryafter because it's more accurate with time sync and whatnot
	if retryAfter != "" {
		parsedAfter, err := strconv.ParseInt(retryAfter, 10, 64)
		if err != nil {
			return err
		}
		b.reset = time.Now().Add(time.Duration(parsedAfter) * time.Millisecond)

	} else if reset != "" {
		unix, err := strconv.ParseInt(reset, 10, 64)
		if err != nil {
			return err
		}

		// Add a second to account for time desync and such
		b.reset = time.Unix(unix, 0).Add(time.Second)
	}

	// Udpate remaining if header is present
	if remaining != "" {
		parsedRemaining, err := strconv.ParseInt(remaining, 10, 32)
		if err != nil {
			return err
		}
		b.remaining = int(parsedRemaining)
	}

	return nil
}
