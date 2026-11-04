package discordgo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// This test takes ~2 seconds to run
func TestRatelimitReset(t *testing.T) {
	rl := NewRatelimiter()

	sendReq := func(endpoint string) {
		bucket := rl.LockBucket(endpoint)

		headers := http.Header(make(map[string][]string))

		headers.Set("X-RateLimit-Remaining", "0")
		// Reset for approx 2 seconds from now
		headers.Set("X-RateLimit-Reset", fmt.Sprint(float64(time.Now().Add(time.Second*2).UnixNano())/1e9))
		headers.Set("Date", time.Now().Format(time.RFC850))

		err := bucket.Release(headers)
		if err != nil {
			t.Errorf("Release returned error: %v", err)
		}
	}

	sent := time.Now()
	sendReq("/guilds/99/channels")
	sendReq("/guilds/55/channels")
	sendReq("/guilds/66/channels")

	sendReq("/guilds/99/channels")
	sendReq("/guilds/55/channels")
	sendReq("/guilds/66/channels")

	// We hit the same endpoint 2 times, so we should only be ratelimited 2 second
	// And always less than 4 seconds (unless you're on a stoneage computer or using swap or something...)
	if time.Since(sent) >= time.Second && time.Since(sent) < time.Second*4 {
		t.Log("OK", time.Since(sent))
	} else {
		t.Error("Did not ratelimit correctly, got:", time.Since(sent))
	}
}

// This test takes ~1 seconds to run
func TestRatelimitGlobal(t *testing.T) {
	rl := NewRatelimiter()

	sendReq := func(endpoint string) {
		bucket := rl.LockBucket(endpoint)

		headers := http.Header(make(map[string][]string))

		headers.Set("X-RateLimit-Global", "1")
		// Reset for approx 1 seconds from now
		headers.Set("X-RateLimit-Reset-After", "1")

		err := bucket.Release(headers)
		if err != nil {
			t.Errorf("Release returned error: %v", err)
		}
	}

	sent := time.Now()

	// This should trigger a global ratelimit
	sendReq("/guilds/99/channels")
	time.Sleep(time.Millisecond * 100)

	// This shouldn't go through in less than 1 second
	sendReq("/guilds/55/channels")

	if time.Since(sent) >= time.Second && time.Since(sent) < time.Second*2 {
		t.Log("OK", time.Since(sent))
	} else {
		t.Error("Did not ratelimit correctly, got:", time.Since(sent))
	}
}

func BenchmarkRatelimitSingleEndpoint(b *testing.B) {
	rl := NewRatelimiter()
	for i := 0; i < b.N; i++ {
		sendBenchReq("/guilds/99/channels", rl)
	}
}

func BenchmarkRatelimitParallelMultiEndpoints(b *testing.B) {
	rl := NewRatelimiter()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sendBenchReq("/guilds/"+strconv.Itoa(i)+"/channels", rl)
			i++
		}
	})
}

// Does not actually send requests, but locks the bucket and releases it with made-up headers
func sendBenchReq(endpoint string, rl *RateLimiter) {
	bucket := rl.LockBucket(endpoint)

	headers := http.Header(make(map[string][]string))

	headers.Set("X-RateLimit-Remaining", "10")
	headers.Set("X-RateLimit-Reset", fmt.Sprint(float64(time.Now().UnixNano())/1e9))
	headers.Set("Date", time.Now().Format(time.RFC850))

	bucket.Release(headers)
}

// A bucket wait stops when the context is done, and leaves the bucket unlocked and its requests untaken.
func TestLockBucketContextStopsWaiting(t *testing.T) {
	rl := NewRatelimiter()
	b := rl.GetBucket("/guilds/99/members")
	b.Remaining = 0
	b.reset = time.Now().Add(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	if _, err := rl.LockBucketObjectContext(ctx, b); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got error %v, want %v", err, context.DeadlineExceeded)
	}
	if took := time.Since(start); took > time.Second {
		t.Errorf("waited %v for a context of 50ms", took)
	}
	if b.Remaining != 0 {
		t.Errorf("a request was taken from the bucket: remaining %d", b.Remaining)
	}

	locked := make(chan struct{})
	go func() {
		b.Lock()
		close(locked)
		b.Unlock()
	}()
	select {
	case <-locked:
	case <-time.After(time.Second):
		t.Fatal("the bucket was left locked")
	}
}

// Waiting for a bucket another request holds stops when the context is done; the bucket is free again once
// that request lets go, with none of its requests taken.
func TestLockBucketContextStopsWaitingForTheLock(t *testing.T) {
	rl := NewRatelimiter()
	b := rl.GetBucket("/guilds/99/members")
	b.Remaining = 5
	b.Lock() // another request holds it

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	if _, err := rl.LockBucketObjectContext(ctx, b); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got error %v, want %v", err, context.DeadlineExceeded)
	}
	if took := time.Since(start); took > time.Second {
		t.Errorf("waited %v for a context of 50ms", took)
	}

	b.Unlock()
	got, err := rl.LockBucketObjectContext(context.Background(), b)
	if err != nil {
		t.Fatal(err)
	}
	if got.Remaining != 4 {
		t.Errorf("remaining %d after one request, want 4", got.Remaining)
	}
	got.Unlock()
}

// A context already done takes no request from a bucket that has some left.
func TestLockBucketContextDoneTakesNoRequest(t *testing.T) {
	rl := NewRatelimiter()
	b := rl.GetBucket("/guilds/99/members")
	b.Remaining = 5

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := rl.LockBucketObjectContext(ctx, b); !errors.Is(err, context.Canceled) {
		t.Fatalf("got error %v, want %v", err, context.Canceled)
	}
	if b.Remaining != 5 {
		t.Errorf("remaining %d, want 5", b.Remaining)
	}
	b.Lock() // not left locked
	b.Unlock()
}
