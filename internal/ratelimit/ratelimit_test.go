package ratelimit

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestParseURL(t *testing.T) {
	tests := map[string]string{
		// Input		Expected Output
		"/gateway":       "/gateway",
		"/channels/5050": "/channels/5050",

		"/channels/5050/messages":             "/channels/5050/messages",
		"/channels/5050/messages/1337?asdasd": "/channels/5050/messages/",
		"/channels/5050/messages/bulk_delete": "/channels/5050/messages/bulk_delete",

		"/channels/5050/permissions/1337": "/channels/5050/permissions/",
		"/channels/5050/invites":          "/channels/5050/invites",

		"/channels/5050/pins":      "/channels/5050/pins",
		"/channels/5050/pins/1337": "/channels/5050/pins/",

		"/guilds/99":                        "/guilds/99",
		"/guilds/99/channels":               "/guilds/99/channels",
		"/guilds/99/members":                "/guilds/99/members",
		"/guilds/99/members/1337":           "/guilds/99/members/",
		"/guilds/99/bans":                   "/guilds/99/bans",
		"/guilds/99/bans/1337":              "/guilds/99/bans/",
		"/guilds/99/roles":                  "/guilds/99/roles",
		"/guilds/99/roles/1337":             "/guilds/99/roles/",
		"/guilds/99/prune":                  "/guilds/99/prune",
		"/guilds/99/regions":                "/guilds/99/regions",
		"/guilds/99/invites":                "/guilds/99/invites",
		"/guilds/99/integrations":           "/guilds/99/integrations",
		"/guilds/99/integrations/1337":      "/guilds/99/integrations/",
		"/guilds/99/integrations/1337/sync": "/guilds/99/integrations//sync",
		"/guilds/99/embed":                  "/guilds/99/embed",

		"/users/@me":             "/users/@me",
		"/users/@me/guilds":      "/users/@me/guilds",
		"/users/@me/guilds/99":   "/users/@me/guilds/99",
		"/users/@me/channels":    "/users/@me/channels",
		"/users/@me/connections": "/users/@me/connections",
		"/users/1337":            "/users/",

		"/invites/code":             "/invites/",
		"/icons/1ffa3.jpg":          "/icons/.jpg",
		"/splashes/1ffa3.jpg":       "/splashes/.jpg",
		"/emojis/1ffa3.png":         "/emojis/.png",
		"/application/123123":       "/application/",
		"/application/123123/bot":   "/application//bot",
		"/users/123/avatars/1231af": "/users//avatars/",
	}

	for input, correct := range tests {
		out := ParseURL(input)
		if out != correct {
			t.Errorf("Incorrect parsed url, input: %q, got: %q, expected: %q", input, out, correct)
		}
	}
}

// This test takes ~2 seconds to run
func TestRatelimitReset(t *testing.T) {
	rl := New()

	sendReq := func(endpoint string) {
		bucket := rl.LockBucket(endpoint)

		headers := http.Header(make(map[string][]string))

		headers.Set("X-RateLimit-Remaining", "0")
		// Reset for approx 2 seconds from now
		headers.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Second*2).Unix(), 10))

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
		t.Error("Did not ratelimit correctly")
	}
}

func BenchmarkRatelimitSingleEndpoint(b *testing.B) {
	rl := New()
	for i := 0; i < b.N; i++ {
		sendBenchReq("/guilds/99/channels", rl)
	}
}

func BenchmarkRatelimitParallelMultiEndpoints(b *testing.B) {
	rl := New()
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
	headers.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Unix(), 10))

	bucket.Release(headers)
}
