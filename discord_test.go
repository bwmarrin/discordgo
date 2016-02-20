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

	envToken    = os.Getenv("DG_TOKEN")    // Token to use when authenticating
	envEmail    = os.Getenv("DG_EMAIL")    // Email to use when authenticating
	envPassword = os.Getenv("DG_PASSWORD") // Password to use when authenticating
	//	envGuild    = os.Getenv("DG_GUILD")    // Guild ID to use for tests
	envChannel = os.Getenv("DG_CHANNEL") // Channel ID to use for tests
	//	envUser     = os.Getenv("DG_USER")     // User ID to use for tests
	envAdmin = os.Getenv("DG_ADMIN") // User ID of admin user to use for tests
)

func init() {
	if envEmail == "" || envPassword == "" || envToken == "" {
		return
	}

	if d, err := New(envEmail, envPassword, envToken); err == nil {
		dg = d
	}
}

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
	d, err := New("asjkldhflkjasdh")
	if err != nil {
		t.Fatalf("New(InvalidToken) returned error: %+v", err)
	}

	// New with just a token does not do any communication, so attempt an api call.
	_, err = d.UserSettings()
	if err == nil {
		t.Errorf("New(InvalidToken), d.UserSettings returned nil error.")
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
		t.Skip("Skipping New(username,InvalidPass), DG_EMAIL not set")
		return
	}
	_, err := New(envEmail, "invalidpassword")
	if err == nil {
		t.Errorf("New(Email, InvalidPass) returned nil error.")
	}
}

// TestNewUserPass tests the New() function with a username and password.
// This should return a valid Session{}, a valid Session.Token.
func TestNewUserPass(t *testing.T) {

	if envEmail == "" || envPassword == "" {
		t.Skip("Skipping New(username,password), DG_EMAIL or DG_PASSWORD not set")
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
}

// TestNewUserPassToken tests the New() function with a username, password and token.
// This should return the same as the TestNewUserPass function.
func TestNewUserPassToken(t *testing.T) {

	if envEmail == "" || envPassword == "" || envToken == "" {
		t.Skip("Skipping New(username,password,token), DG_EMAIL, DG_PASSWORD or DG_TOKEN not set")
		return
	}

	d, err := New(envEmail, envPassword, envToken)
	if err != nil {
		t.Fatalf("New(user,pass,token) returned error: %+v", err)
	}

	if d == nil {
		t.Fatal("New(user,pass,token), d is nil, should be Session{}")
	}

	if d.Token == "" {
		t.Fatal("New(user,pass,token), d.Token is empty, should be a valid Token.")
	}
}

func TestOpenClose(t *testing.T) {
	if envToken == "" {
		t.Skip("Skipping TestClose, DG_TOKEN not set")
	}

	d, err := New(envToken)
	if err != nil {
		t.Fatalf("TestClose, New(envToken) returned error: %+v", err)
	}

	if err = d.Open(); err != nil {
		t.Fatalf("TestClose, d.Open failed: %+v", err)
	}

	if !waitBoolEqual(10*time.Second, &d.DataReady, true) {
		t.Fatal("DataReady never became true.")
	}

	// TODO find a better way
	// Add a small sleep here to make sure heartbeat and other events
	// have enough time to get fired.  Need a way to actually check
	// those events.
	time.Sleep(2 * time.Second)

	// UpdateStatus - maybe we move this into wsapi_test.go but the websocket
	// created here is needed.  This helps tests that the websocket was setup
	// and it is working.
	if err = d.UpdateStatus(0, time.Now().String()); err != nil {
		t.Errorf("UpdateStatus error: %+v", err)
	}

	if err = d.Close(); err != nil {
		t.Fatalf("TestClose, d.Close failed: %+v", err)
	}
}

func TestAddHandler(t *testing.T) {
	testHandlerCalled := 0
	testHandler := func(s *Session, m *MessageCreate) {
		testHandlerCalled++
	}

	interfaceHandlerCalled := 0
	interfaceHandler := func(s *Session, i interface{}) {
		interfaceHandlerCalled++
	}

	bogusHandlerCalled := false
	bogusHandler := func(s *Session, se *Session) {
		bogusHandlerCalled = true
	}

	d := Session{}
	d.AddHandler(testHandler)
	d.AddHandler(testHandler)

	d.AddHandler(interfaceHandler)
	d.AddHandler(bogusHandler)

	d.handle(&MessageCreate{})
	d.handle(&MessageDelete{})

	// testHandler will be called twice because it was added twice.
	if testHandlerCalled != 2 {
		t.Fatalf("testHandler was not called twice.")
	}

	// interfaceHandler will be called twice, once for each event.
	if interfaceHandlerCalled != 2 {
		t.Fatalf("interfaceHandler was not called twice.")
	}

	if bogusHandlerCalled {
		t.Fatalf("bogusHandler was called.")
	}
}

func TestRemoveHandler(t *testing.T) {
	testHandlerCalled := 0
	testHandler := func(s *Session, m *MessageCreate) {
		testHandlerCalled++
	}

	d := Session{}
	r := d.AddHandler(testHandler)

	d.handle(&MessageCreate{})

	r()

	d.handle(&MessageCreate{})

	// testHandler will be called once, as it was removed in between calls.
	if testHandlerCalled != 1 {
		t.Fatalf("testHandler was not called once.")
	}
}
