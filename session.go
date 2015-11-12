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
	OnMessageCreate           func(*Session, Message)
	OnTypingStart             func(*Session, Event)
	OnMessageAck              func(*Session, Event)
	OnMessageUpdate           func(*Session, Event)
	OnMessageDelete           func(*Session, Event)
	OnPresenceUpdate          func(*Session, Event)
	OnChannelCreate           func(*Session, Event)
	OnChannelUpdate           func(*Session, Event)
	OnChannelDelete           func(*Session, Event)
	OnGuildCreate             func(*Session, Event)
	OnGuildDelete             func(*Session, Event)
	OnGuildMemberAdd          func(*Session, Event)
	OnGuildMemberRemove       func(*Session, Event)
	OnGuildMemberDelete       func(*Session, Event) // which is it?
	OnGuildMemberUpdate       func(*Session, Event)
	OnGuildRoleCreate         func(*Session, Event)
	OnGuildRoleDelete         func(*Session, Event)
	OnGuildIntegrationsUpdate func(*Session, Event)

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

func (session *Session) Servers() (servers []Server, err error) {
	servers, err = Servers(session, "@me")
	return
}

// Logout ends a session and logs out from the Discord REST API.
func (session *Session) Logout() (err error) {
	err = Logout(session)
	return
}
