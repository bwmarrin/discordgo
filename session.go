/******************************************************************************
 * A Discord API for Golang.
 */

package discordgo

import "github.com/gorilla/websocket"

// A Session represents a connection to the Discord REST API.
// token : The authentication token returned from Discord
// Debug : If set to ture debug logging will be displayed.
type Session struct {
	Token   string // Authentication token for this session
	Gateway string // Websocket Gateway for this session
	Debug   bool   // Debug for printing JSON request/responses
	Cache   int    // number in X to cache some responses

	// Settable Callback functions for Websocket Events
	OnEvent                   func(*Session, Event) // should Event be *Event?
	OnReady                   func(*Session, Ready)
	OnTypingStart             func(*Session, TypingStart)
	OnMessageCreate           func(*Session, Message)
	OnMessageUpdate           func(*Session, Message)
	OnMessageDelete           func(*Session, MessageDelete)
	OnMessageAck              func(*Session, MessageAck)
	OnPresenceUpdate          func(*Session, PresenceUpdate)
	OnVoiceStateUpdate        func(*Session, VoiceState)
	OnChannelCreate           func(*Session, Channel)
	OnChannelUpdate           func(*Session, Channel)
	OnChannelDelete           func(*Session, Channel)
	OnGuildCreate             func(*Session, Guild)
	OnGuildDelete             func(*Session, Guild)
	OnGuildMemberAdd          func(*Session, Member)
	OnGuildMemberRemove       func(*Session, Member)
	OnGuildMemberDelete       func(*Session, Member) // which is it?
	OnGuildMemberUpdate       func(*Session, Member)
	OnGuildRoleCreate         func(*Session, Role)
	OnGuildRoleDelete         func(*Session, Role)
	OnGuildIntegrationsUpdate func(*Session, GuildIntegrationsUpdate)

	// OnMessageCreate func(Session, Event, Message)
	// ^^ Any value to passing session, event, message?
	// probably just the Message is all one would need.
	// but having the sessin could be handy?

	wsConn *websocket.Conn
	//TODO, add bools for like.
	// are we connnected to websocket?
	// have we authenticated to login?
	// lets put all the general session
	// tracking and infos here.. clearly
}

/******************************************************************************
 * The below functions are "shortcut" methods for functions in restapi.go
 * Reference the client.go file for more documentation.
 */
func (session *Session) Login(email string, password string) (token string, err error) {
	token, err = Login(session, email, password)
	return
}

func (session *Session) Self() (user User, err error) {
	user, err = Users(session, "@me")
	return
}

func (session *Session) PrivateChannels() (channels []Channel, err error) {
	channels, err = PrivateChannels(session, "@me")
	return
}

func (session *Session) Guilds() (servers []Guild, err error) {
	servers, err = Guilds(session, "@me")
	return
}

// Logout ends a session and logs out from the Discord REST API.
func (session *Session) Logout() (err error) {
	err = Logout(session)
	return
}
