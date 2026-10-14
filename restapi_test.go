package discordgo

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

//////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////// START OF TESTS

// TestChannelMessageSend tests the ChannelMessageSend() function. This should not return an error.
func TestChannelMessageSend(t *testing.T) {

	if envChannel == "" {
		t.Skip("Skipping, DG_CHANNEL not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	_, err := dg.ChannelMessageSend(envChannel, "Running REST API Tests!")
	if err != nil {
		t.Errorf("ChannelMessageSend returned error: %+v", err)
	}
}

/*
// removed for now, only works on BOT accounts now
func TestUserAvatar(t *testing.T) {

	if dg == nil {
		t.Skip("Cannot TestUserAvatar, dg not set.")
	}

	u, err := dg.User("@me")
	if err != nil {
		t.Error("error fetching @me user,", err)
	}

	a, err := dg.UserAvatar(u.ID)
	if err != nil {
		if err.Error() == `HTTP 404 NOT FOUND, {"code": 0, "message": "404: Not Found"}` {
			t.Skip("Skipped, @me doesn't have an Avatar")
		}
		t.Errorf(err.Error())
	}

	if a == nil {
		t.Errorf("a == nil, should be image.Image")
	}
}
*/

/* Running this causes an error due to 2/hour rate limit on username changes
func TestUserUpdate(t *testing.T) {
	if dg == nil {
		t.Skip("Cannot test logout, dg not set.")
	}

	u, err := dg.User("@me")
	if err != nil {
		t.Errorf(err.Error())
	}

	s, err := dg.UserUpdate(envEmail, envPassword, "testname", u.Avatar, "")
	if err != nil {
		t.Error(err.Error())
	}
	if s.Username != "testname" {
		t.Error("Username != testname")
	}
	s, err = dg.UserUpdate(envEmail, envPassword, u.Username, u.Avatar, "")
	if err != nil {
		t.Error(err.Error())
	}
	if s.Username != u.Username {
		t.Error("Username != " + u.Username)
	}
}
*/

//func (s *Session) UserChannelCreate(recipientID string) (st *Channel, err error) {

func TestUserChannelCreate(t *testing.T) {
	if dg == nil {
		t.Skip("Cannot TestUserChannelCreate, dg not set.")
	}

	if envAdmin == "" {
		t.Skip("Skipped, DG_ADMIN not set.")
	}

	_, err := dg.UserChannelCreate(envAdmin)
	if err != nil {
		t.Errorf(err.Error())
	}

	// TODO make sure the channel was added
}

func TestUserGuilds(t *testing.T) {
	if dg == nil {
		t.Skip("Cannot TestUserGuilds, dg not set.")
	}

	_, err := dg.UserGuilds(10, "", "", false)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGateway(t *testing.T) {

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}
	_, err := dg.Gateway()
	if err != nil {
		t.Errorf("Gateway() returned error: %+v", err)
	}
}

func TestGatewayBot(t *testing.T) {

	if dgBot == nil {
		t.Skip("Skipping, dgBot not set.")
	}
	_, err := dgBot.GatewayBot()
	if err != nil {
		t.Errorf("GatewayBot() returned error: %+v", err)
	}
}

func TestVoiceRegions(t *testing.T) {

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	_, err := dg.VoiceRegions()
	if err != nil {
		t.Errorf("VoiceRegions() returned error: %+v", err)
	}
}
func TestGuildRoles(t *testing.T) {

	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	_, err := dg.GuildRoles(envGuild)
	if err != nil {
		t.Errorf("GuildRoles(envGuild) returned error: %+v", err)
	}

}

func TestGuildMemberNickname(t *testing.T) {

	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	err := dg.GuildMemberNickname(envGuild, "@me/nick", "B1nzyRocks")
	if err != nil {
		t.Errorf("GuildNickname returned error: %+v", err)
	}
}

// TestChannelMessageSend2 tests the ChannelMessageSend() function. This should not return an error.
func TestChannelMessageSend2(t *testing.T) {

	if envChannel == "" {
		t.Skip("Skipping, DG_CHANNEL not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	_, err := dg.ChannelMessageSend(envChannel, "All done running REST API Tests!")
	if err != nil {
		t.Errorf("ChannelMessageSend returned error: %+v", err)
	}
}

// TestGuildPruneCount tests GuildPruneCount() function. This should not return an error.
func TestGuildPruneCount(t *testing.T) {

	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	_, err := dg.GuildPruneCount(envGuild, 1)
	if err != nil {
		t.Errorf("GuildPruneCount returned error: %+v", err)
	}
}

/*
// TestGuildPrune tests GuildPrune() function. This should not return an error.
func TestGuildPrune(t *testing.T) {

	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	_, err := dg.GuildPrune(envGuild, 1)
	if err != nil {
		t.Errorf("GuildPrune returned error: %+v", err)
	}
}
*/

func Test_unmarshal(t *testing.T) {
	err := unmarshal([]byte{}, &struct{}{})
	if !errors.Is(err, ErrJSONUnmarshal) {
		t.Errorf("Unexpected error type: %T", err)
	}
}

func TestWithContext(t *testing.T) {
	// Set up a test context.
	type key struct{}
	ctx := context.WithValue(context.Background(), key{}, "value")

	// Set up a test client.
	session, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	testErr := errors.New("test")

	// Intercept the request to assert the context.
	session.Client.Transport = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		val, _ := r.Context().Value(key{}).(string)
		if val != "value" {
			t.Errorf("missing value in context (got %q, wanted %q)", val, "value")
		}
		return nil, testErr
	})

	// Run any client method using WithContext.
	_, err = session.User("", WithContext(ctx))

	// Verify that the assertion code was actually run.
	if !errors.Is(err, testErr) {
		t.Errorf("unexpected error %v returned from client", err)
	}
}

// roundTripperFunc implements http.RoundTripper.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestGuildMembers_GuildID(t *testing.T) {
	session, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	testGuildID := "123456789012345678"
	membersJSON := `[{"user": {"id": "111"}, "roles": []}, {"user": {"id": "222"}, "roles": []}]`

	session.Client.Transport = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != EndpointGuildMembers(testGuildID) {
			t.Errorf("unexpected request path: %q", r.URL.Path)
		}
		res := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(membersJSON)),
			Header:     make(http.Header),
		}
		return res, nil
	})

	members, err := session.GuildMembers(testGuildID, "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}

	for _, m := range members {
		if m.GuildID != testGuildID {
			t.Errorf("expected member GuildID to be %q, got %q", testGuildID, m.GuildID)
		}
	}
}
