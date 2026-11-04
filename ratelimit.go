package discordgo

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// customRateLimit holds information for defining a custom rate limit
type customRateLimit struct {
	suffix   string
	requests int
	reset    time.Duration
}

// RateLimiter holds all ratelimit buckets
type RateLimiter struct {
	sync.Mutex
	global           *int64
	buckets          map[string]*Bucket
	globalRateLimit  time.Duration
	customRateLimits []*customRateLimit
}

// NewRatelimiter returns a new RateLimiter
func NewRatelimiter() *RateLimiter {

	return &RateLimiter{
		buckets: make(map[string]*Bucket),
		global:  new(int64),
		customRateLimits: []*customRateLimit{
			{
				suffix:   "//reactions//",
				requests: 1,
				reset:    200 * time.Millisecond,
			},
		},
	}
}

// GetBucket retrieves or creates a bucket
func (r *RateLimiter) GetBucket(key string) *Bucket {
	r.Lock()
	defer r.Unlock()

	if bucket, ok := r.buckets[key]; ok {
		return bucket
	}

	b := &Bucket{
		Remaining: 1,
		Key:       key,
		global:    r.global,
	}

	// Check if there is a custom ratelimit set for this bucket ID.
	for _, rl := range r.customRateLimits {
		if strings.HasSuffix(b.Key, rl.suffix) {
			b.customRateLimit = rl
			break
		}
	}

	r.buckets[key] = b
	return b
}

// GetWaitTime returns the duration you should wait for a Bucket
func (r *RateLimiter) GetWaitTime(b *Bucket, minRemaining int) time.Duration {
	// If we ran out of calls and the reset time is still ahead of us
	// then we need to take it easy and relax a little
	if b.Remaining < minRemaining && b.reset.After(time.Now()) {
		return b.reset.Sub(time.Now())
	}

	// Check for global ratelimits
	sleepTo := time.Unix(0, atomic.LoadInt64(r.global))
	if now := time.Now(); now.Before(sleepTo) {
		return sleepTo.Sub(now)
	}

	return 0
}

// LockBucket Locks until a request can be made
func (r *RateLimiter) LockBucket(bucketID string) *Bucket {
	return r.LockBucketObject(r.GetBucket(bucketID))
}

// LockBucketObject Locks an already resolved bucket until a request can be made
func (r *RateLimiter) LockBucketObject(b *Bucket) *Bucket {
	b, _ = r.LockBucketObjectContext(context.Background(), b)
	return b
}

// LockBucketContext is LockBucket, but stops waiting for the rate limit when ctx is done.
func (r *RateLimiter) LockBucketContext(ctx context.Context, bucketID string) (*Bucket, error) {
	return r.LockBucketObjectContext(ctx, r.GetBucket(bucketID))
}

// LockBucketObjectContext is LockBucketObject, but stops waiting when ctx is done: for the bucket while another
// request holds it, and for the rate limit. It then returns ctx.Err() with the bucket unlocked and none of its
// remaining requests taken.
func (r *RateLimiter) LockBucketObjectContext(ctx context.Context, b *Bucket) (*Bucket, error) {
	if err := lockContext(ctx, b); err != nil {
		return nil, err
	}

	if wait := r.GetWaitTime(b, 1); wait > 0 {
		if err := sleepContext(ctx, wait); err != nil {
			b.Unlock()
			return nil, err
		}
	}
	// A context that ended without a wait to notice it takes no request either
	if err := ctx.Err(); err != nil {
		b.Unlock()
		return nil, err
	}

	b.Remaining--
	return b, nil
}

// lockContext locks b, or returns ctx.Err() without it if ctx is done first. A context that can never be done
// locks b directly; any other waits in a goroutine, which, if ctx gives up first, unlocks b as soon as it has it.
func lockContext(ctx context.Context, b *Bucket) error {
	if ctx.Done() == nil {
		b.Lock()
		return nil
	}

	locked := make(chan struct{})
	go func() {
		b.Lock()
		close(locked)
	}()

	select {
	case <-locked:
		return nil
	case <-ctx.Done():
		go func() {
			<-locked
			b.Unlock()
		}()
		return ctx.Err()
	}
}

// sleepContext waits for d, or until ctx is done, whichever comes first. It returns ctx.Err() if ctx ended the wait.
func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Bucket represents a ratelimit bucket, each bucket gets ratelimited individually (-global ratelimits)
type Bucket struct {
	sync.Mutex
	Key       string
	Remaining int
	limit     int
	reset     time.Time
	global    *int64

	lastReset       time.Time
	customRateLimit *customRateLimit
	Userdata        interface{}
}

// Release unlocks the bucket and reads the headers to update the buckets ratelimit info
// and locks up the whole thing in case if there's a global ratelimit.
func (b *Bucket) Release(headers http.Header) error {
	defer b.Unlock()

	// Check if the bucket uses a custom ratelimiter
	if rl := b.customRateLimit; rl != nil {
		if time.Now().Sub(b.lastReset) >= rl.reset {
			b.Remaining = rl.requests - 1
			b.lastReset = time.Now()
		}
		if b.Remaining < 1 {
			b.reset = time.Now().Add(rl.reset)
		}
		return nil
	}

	if headers == nil {
		return nil
	}

	remaining := headers.Get("X-RateLimit-Remaining")
	reset := headers.Get("X-RateLimit-Reset")
	global := headers.Get("X-RateLimit-Global")
	resetAfter := headers.Get("X-RateLimit-Reset-After")

	// Update global and per bucket reset time if the proper headers are available
	// If global is set, then it will block all buckets until after Retry-After
	// If Retry-After without global is provided it will use that for the new reset
	// time since it's more accurate than X-RateLimit-Reset.
	// If Retry-After after is not proided, it will update the reset time from X-RateLimit-Reset
	if resetAfter != "" {
		parsedAfter, err := strconv.ParseFloat(resetAfter, 64)
		if err != nil {
			return err
		}

		whole, frac := math.Modf(parsedAfter)
		resetAt := time.Now().Add(time.Duration(whole) * time.Second).Add(time.Duration(frac*1000) * time.Millisecond)

		// Lock either this single bucket or all buckets
		if global != "" {
			atomic.StoreInt64(b.global, resetAt.UnixNano())
		} else {
			b.reset = resetAt
		}
	} else if reset != "" {
		// Calculate the reset time by using the date header returned from discord
		discordTime, err := http.ParseTime(headers.Get("Date"))
		if err != nil {
			return err
		}

		unix, err := strconv.ParseFloat(reset, 64)
		if err != nil {
			return err
		}

		// Calculate the time until reset and add it to the current local time
		// some extra time is added because without it i still encountered 429's.
		// The added amount is the lowest amount that gave no 429's
		// in 1k requests
		whole, frac := math.Modf(unix)
		delta := time.Unix(int64(whole), 0).Add(time.Duration(frac*1000)*time.Millisecond).Sub(discordTime) + time.Millisecond*250
		b.reset = time.Now().Add(delta)
	}

	// Udpate remaining if header is present
	if remaining != "" {
		parsedRemaining, err := strconv.ParseInt(remaining, 10, 32)
		if err != nil {
			return err
		}
		b.Remaining = int(parsedRemaining)
	}

	return nil
}
