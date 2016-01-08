package discordgo_test

import (
	"os"
	"runtime"
	"testing"
	"time"

	. "github.com/bwmarrin/discordgo"
)

//////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////// VARS NEEDED FOR TESTING
var (
	dg *Session // Stores global discordgo session

	envToken    string = os.Getenv("DG_TOKEN")    // Token to use when authenticating
	envUsername string = os.Getenv("DG_USERNAME") // Username to use when authenticating
	envPassword string = os.Getenv("DG_PASSWORD") // Password to use when authenticating
	envGuild    string = os.Getenv("DG_GUILD")    // Guild ID to use for tests
	envChannel  string = os.Getenv("DG_CHANNEL")  // Channel ID to use for tests
	envUser     string = os.Getenv("DG_USER")     // User ID to use for tests
	envAdmin    string = os.Getenv("DG_ADMIN")    // User ID of admin user to use for tests
)

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////// HELPER FUNCTIONS USED FOR TESTING

// This waits x time for the check bool to be the want bool
func waitBoolEqual(timeout time.Duration, check *bool, want bool) bool {

	start := time.Now()
	for {
		if *check == want {
			return true
		}

		if time.Since(start) > timeout {
			return false
		}

		runtime.Gosched()
	}
}

// Checks if we're connected to Discord
func isConnected() bool {

	if dg == nil {
		return false
	}

	if dg.Token == "" {
		return false
	}

	// Need a way to see if the ws connection is nil

	if !waitBoolEqual(10*time.Second, &dg.DataReady, true) {
		return false
	}

	return true
}

//////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////// START OF TESTS

// TestNew tests the New() function without any arguments.  This should return
// a valid Session{} struct and no errors.
func TestNew(t *testing.T) {

	_, err := New()
	if err != nil {
		t.Errorf("New() returned error: %+v", err)
	}
}

// TestNewUserPass tests the New() function with a username and password.
// This should return a valid Session{}, a valid Session.Token, and open
// a websocket connection to Discord.
func TestNewUserPass(t *testing.T) {

	if isConnected() {
		t.Skip("Skipping New(username,password), already connected.")
	}

	if envUsername == "" || envPassword == "" {
		t.Skip("Skipping New(username,password), DG_USERNAME or DG_PASSWORD not set")
		return
	}
	// Not testing yet.
}

// TestNewToken tests the New() function with a Token.  This should return
// the same as the TestNewUserPass function.
func TestNewToken(t *testing.T) {

	if isConnected() {
		t.Skip("Skipping New(token), already connected.")
	}

	if envToken == "" {
		t.Skip("Skipping New(token), DG_TOKEN not set")
	}

	d, err := New(envToken)
	if err != nil {
		t.Fatalf("New(envToken) returned error: %+v", err)
	}

	if d == nil {
		t.Fatal("New(envToken), d is nil, should be Session{}")
	}

	if d.Token == "" {
		t.Fatal("New(envToken), d.Token is empty, should be a valid Token.")
	}

	if !waitBoolEqual(10*time.Second, &d.DataReady, true) {
		t.Fatal("New(envToken), d.DataReady is false after 10 seconds.  Should be true.")
	}

	t.Log("Successfully connected to Discord.")
	dg = d

}
