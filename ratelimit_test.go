package discordgo

import (
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
