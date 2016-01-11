package discordgo

import (
	"os"
	"runtime"
	"testing"
	"time"
)

//////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////// VARS NEEDED FOR TESTING
var (
	dg *Session // Stores global discordgo session

	envToken    string = os.Getenv("DG_TOKEN")    // Token to use when authenticating
	envEmail    string = os.Getenv("DG_EMAIL")    // Email to use when authenticating
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

// TestInvalidToken tests the New() function with an invalid token
func TestInvalidToken(t *testing.T) {

	_, err := New("asjkldhflkjasdh")
	if err == nil {
		t.Errorf("New(InvalidToken) returned nil error.")
	}

}

// TestInvalidUserPass tests the New() function with an invalid Email and Pass
func TestInvalidEmailPass(t *testing.T) {

	_, err := New("invalidemail", "invalidpassword")
	if err == nil {
		t.Errorf("New(InvalidEmail, InvalidPass) returned nil error.")
	}

}

// TestInvalidPass tests the New() function with an invalid Password
func TestInvalidPass(t *testing.T) {

	if envEmail == "" {
		t.Skip("Skipping New(username,InvalidPass), DG_USERNAME not set")
		return
	}
	_, err := New(envEmail, "invalidpassword")
	if err == nil {
		t.Errorf("New(Email, InvalidPass) returned nil error.")
	}
}

// TestNewUserPass tests the New() function with a username and password.
// This should return a valid Session{}, a valid Session.Token, and open
// a websocket connection to Discord.
func TestNewUserPass(t *testing.T) {

	if envEmail == "" || envPassword == "" {
		t.Skip("Skipping New(username,password), DG_USERNAME or DG_PASSWORD not set")
		return
	}

	d, err := New(envEmail, envPassword)
	if err != nil {
		t.Fatalf("New(user,pass) returned error: %+v", err)
	}

	if d == nil {
		t.Fatal("New(user,pass), d is nil, should be Session{}")
	}

	if d.Token == "" {
		t.Fatal("New(user,pass), d.Token is empty, should be a valid Token.")
	}

	if !waitBoolEqual(10*time.Second, &d.DataReady, true) {
		t.Fatal("New(user,pass), d.DataReady is false after 10 seconds.  Should be true.")
	}

	t.Log("Successfully connected to Discord via New(user,pass).")
	dg = d

	// Not testing yet.
}

// TestNewToken tests the New() function with a Token.  This should return
// the same as the TestNewUserPass function.
func TestNewToken(t *testing.T) {

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

	t.Log("Successfully connected to Discord via New(token).")
	dg = d

}
