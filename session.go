/******************************************************************************
 * A Discord API for Golang.
 *
 * This file has structs and functions specific to a session.
 *
 * A session is a single connection to Discord for a given
 * user and all REST and Websock API functions exist within
 * a session.
 *
 * See the restapi.go and wsapi.go for more information.
 */

package discordgo

import (
	"net"

	"github.com/gorilla/websocket"
)

// A Session represents a connection to the Discord REST API.
// token : The authentication token returned from Discord
// Debug : If set to ture debug logging will be displayed.
type Session struct {
	Token     string // Authentication token for this session
	Debug     bool   // Debug for printing JSON request/responses
	Cache     int    // number in X to cache some responses
	SessionID string // from websocket READY packet

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
	OnGuildUpdate             func(*Session, Guild)
	OnGuildDelete             func(*Session, Guild)
	OnGuildMemberAdd          func(*Session, Member)
	OnGuildMemberRemove       func(*Session, Member)
	OnGuildMemberDelete       func(*Session, Member) // which is it?
	OnGuildMemberUpdate       func(*Session, Member)
	OnGuildRoleCreate         func(*Session, GuildRole)
	OnGuildRoleUpdate         func(*Session, GuildRole)
	OnGuildRoleDelete         func(*Session, GuildRoleDelete)
	OnGuildIntegrationsUpdate func(*Session, GuildIntegrationsUpdate)

	wsConn *websocket.Conn
	//TODO, add bools for like.
	// are we connnected to websocket?
	// have we authenticated to login?
	// lets put all the general session
	// tracking and infos here.. clearly

	// Everything below here is used for Voice testing.
	// This stuff is almost guarenteed to change a lot
	// and is even a bit hackish right now.
	VwsConn    *websocket.Conn // new for voice
	VSessionID string
	VToken     string
	VEndpoint  string
	VGuildID   string
	VChannelID string
	Vop2       VoiceOP2
	UDPConn    *net.UDPConn
}
