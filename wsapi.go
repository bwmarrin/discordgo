// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains low level functions for interacting with the Discord
// data websocket interface.

package discordgo

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

// Open opens a websocket connection to Discord.
func (s *Session) Open() (err error) {

	// Get the gateway to use for the Websocket connection
	g, err := s.Gateway()
	if err != nil {
		return
	}

	// TODO: See if there's a use for the http response.
	// conn, response, err := websocket.DefaultDialer.Dial(session.Gateway, nil)
	s.wsConn, _, err = websocket.DefaultDialer.Dial(g, nil)
	return
}

type handshakeProperties struct {
	OS              string `json:"$os"`
	Browser         string `json:"$browser"`
	Device          string `json:"$device"`
	Referer         string `json:"$referer"`
	ReferringDomain string `json:"$referring_domain"`
}

type handshakeData struct {
	Version    int                 `json:"v"`
	Token      string              `json:"token"`
	Properties handshakeProperties `json:"properties"`
}

type handshakeOp struct {
	Op   int           `json:"op"`
	Data handshakeData `json:"d"`
}

// Handshake sends the client data to Discord during websocket initial connection.
func (s *Session) Handshake() (err error) {
	// maybe this is SendOrigin? not sure the right name here

	data := handshakeOp{2, handshakeData{3, s.Token, handshakeProperties{runtime.GOOS, "Discordgo v" + VERSION, "", "", ""}}}
	err = s.wsConn.WriteJSON(data)
	return
}

type updateStatusGame struct {
	Name string `json:"name"`
}

type updateStatusData struct {
	IdleSince *int              `json:"idle_since"`
	Game      *updateStatusGame `json:"game"`
}

type updateStatusOp struct {
	Op   int              `json:"op"`
	Data updateStatusData `json:"d"`
}

// UpdateStatus is used to update the authenticated user's status.
// If idle>0 then set status to idle.  If game>0 then set game.
// if otherwise, set status to active, and no game.
func (s *Session) UpdateStatus(idle int, game string) (err error) {

	var usd updateStatusData
	if idle > 0 {
		usd.IdleSince = &idle
	}
	if game != "" {
		usd.Game = &updateStatusGame{game}
	}

	data := updateStatusOp{3, usd}
	err = s.wsConn.WriteJSON(data)

	return
}

// Listen starts listening to the websocket connection for events.
func (s *Session) Listen() (err error) {
	// TODO: need a channel or something to communicate
	// to this so I can tell it to stop listening

	if s.wsConn == nil {
		fmt.Println("No websocket connection exists.")
		return // TODO need to return an error.
	}

	// Make sure Listen is not already running
	s.listenLock.Lock()
	if s.listenChan != nil {
		s.listenLock.Unlock()
		return
	}
	s.listenChan = make(chan struct{})
	s.listenLock.Unlock()

	// this is ugly.
	defer func() {
		if s.listenChan == nil {
			return
		}
		select {
		case <-s.listenChan:
			break
		default:
			close(s.listenChan)
		}
		s.listenChan = nil
	}()

	// this is ugly.
	defer func() {
		if s.heartbeatChan == nil {
			return
		}
		select {
		case <-s.heartbeatChan:
			break
		default:
			close(s.heartbeatChan)
		}
		s.listenChan = nil
	}()

	for {
		messageType, message, err := s.wsConn.ReadMessage()
		if err != nil {
			fmt.Println("Websocket Listen Error", err)
			// TODO Log error
			break
		}
		go s.event(messageType, message)

		// If our chan gets closed, exit out of this loop.
		// TODO: Can we make this smarter, using select
		// and some other trickery?  http://www.goinggo.net/2013/10/my-channel-select-bug.html
		if s.listenChan == nil {
			return nil
		}
	}

	return
}

// Not sure how needed this is and where it would be best to call it.
// somewhere.

func unmarshalEvent(event *Event, i interface{}) (err error) {
	if err = unmarshal(event.RawData, i); err != nil {
		fmt.Println(event.Type, err)
		printJSON(event.RawData) // TODO: Better error loggingEvent.
	}
	return
}

// Front line handler for all Websocket Events.  Determines the
// event type and passes the message along to the next handler.

// event is the front line handler for all events.  This needs to be
// broken up into smaller functions to be more idiomatic Go.
// Events will be handled by any implemented handler in Session.
// All unhandled events will then be handled by OnEvent.
func (s *Session) event(messageType int, message []byte) (err error) {

	if s.Debug {
		printJSON(message)
	}

	var e *Event
	if err = unmarshal(message, &e); err != nil {
		fmt.Println(err)
		return
	}

	switch e.Type {

	case "READY":
		var st *Ready
		if err = unmarshalEvent(e, &st); err == nil {
			if s.StateEnabled {
				s.State.OnReady(st)
			}
			if s.OnReady != nil {
				s.OnReady(s, st)
			}
			go s.Heartbeat(st.HeartbeatInterval)
		}
		if s.OnReady != nil {
			return
		}
	case "VOICE_SERVER_UPDATE":
		// TEMP CODE FOR TESTING VOICE
		var st *VoiceServerUpdate
		if err = unmarshalEvent(e, &st); err == nil {
			s.onVoiceServerUpdate(st)
		}
		return
	case "VOICE_STATE_UPDATE":
		// TEMP CODE FOR TESTING VOICE
		var st *VoiceState
		if err = unmarshalEvent(e, &st); err == nil {
			s.onVoiceStateUpdate(st)
		}
		return
	case "USER_UPDATE":
		if s.OnUserUpdate != nil {
			var st *User
			if err = unmarshalEvent(e, &st); err == nil {
				s.OnUserUpdate(s, st)
			}
			return
		}
	case "PRESENCE_UPDATE":
		if s.OnPresenceUpdate != nil {
			var st *PresenceUpdate
			if err = unmarshalEvent(e, &st); err == nil {
				s.OnPresenceUpdate(s, st)
			}
			return
		}
	case "TYPING_START":
		if s.OnTypingStart != nil {
			var st *TypingStart
			if err = unmarshalEvent(e, &st); err == nil {
				s.OnTypingStart(s, st)
			}
			return
		}
		/* Never seen this come in but saw it in another Library.
		case "MESSAGE_ACK":
			if s.OnMessageAck != nil {
			}
		*/
	case "MESSAGE_CREATE":
		stateEnabled := s.StateEnabled && s.State.MaxMessageCount > 0
		if !stateEnabled && s.OnMessageCreate == nil {
			break
		}
		var st *Message
		if err = unmarshalEvent(e, &st); err == nil {
			if stateEnabled {
				s.State.MessageAdd(st)
			}
			if s.OnMessageCreate != nil {
				s.OnMessageCreate(s, st)
			}
		}
		if s.OnMessageCreate != nil {
			return
		}
	case "MESSAGE_UPDATE":
		stateEnabled := s.StateEnabled && s.State.MaxMessageCount > 0
		if !stateEnabled && s.OnMessageUpdate == nil {
			break
		}
		var st *Message
		if err = unmarshalEvent(e, &st); err == nil {
			if stateEnabled {
				s.State.MessageAdd(st)
			}
			if s.OnMessageUpdate != nil {
				s.OnMessageUpdate(s, st)
			}
		}
		return
		if s.OnMessageUpdate != nil {
			return
		}
	case "MESSAGE_DELETE":
		stateEnabled := s.StateEnabled && s.State.MaxMessageCount > 0
		if !stateEnabled && s.OnMessageDelete == nil {
			break
		}
		var st *Message
		if err = unmarshalEvent(e, &st); err == nil {
			if stateEnabled {
				s.State.MessageRemove(st)
			}
			if s.OnMessageDelete != nil {
				s.OnMessageDelete(s, st)
			}
		}
		return
		if s.OnMessageDelete != nil {
			return
		}
	case "MESSAGE_ACK":
		if s.OnMessageAck != nil {
			var st *MessageAck
			if err = unmarshalEvent(e, &st); err == nil {
				s.OnMessageAck(s, st)
			}
			return
		}
	case "CHANNEL_CREATE":
		if !s.StateEnabled && s.OnChannelCreate == nil {
			break
		}
		var st *Channel
		if err = unmarshalEvent(e, &st); err == nil {
			if s.StateEnabled {
				s.State.ChannelAdd(st)
			}
			if s.OnChannelCreate != nil {
				s.OnChannelCreate(s, st)
			}
		}
		if s.OnChannelCreate != nil {
			return
		}
	case "CHANNEL_UPDATE":
		if !s.StateEnabled && s.OnChannelUpdate == nil {
			break
		}
		var st *Channel
		if err = unmarshalEvent(e, &st); err == nil {
			if s.StateEnabled {
				s.State.ChannelAdd(st)
			}
			if s.OnChannelUpdate != nil {
				s.OnChannelUpdate(s, st)
			}
		}
		if s.OnChannelUpdate != nil {
			return
		}
	case "CHANNEL_DELETE":
		if !s.StateEnabled && s.OnChannelDelete == nil {
			break
		}
		var st *Channel
		if err = unmarshalEvent(e, &st); err == nil {
			if s.StateEnabled {
				s.State.ChannelRemove(st)
			}
			if s.OnChannelDelete != nil {
				s.OnChannelDelete(s, st)
			}
		}
		if s.OnChannelDelete != nil {
			return
		}
	case "GUILD_CREATE":
		if !s.StateEnabled && s.OnGuildCreate == nil {
			break
		}
		var st *Guild
		if err = unmarshalEvent(e, &st); err == nil {
			if s.StateEnabled {
				s.State.GuildAdd(st)
			}
			if s.OnGuildCreate != nil {
				s.OnGuildCreate(s, st)
			}
		}
		if s.OnGuildCreate != nil {
			return
		}
	case "GUILD_UPDATE":
		if !s.StateEnabled && s.OnGuildUpdate == nil {
			break
		}
		var st *Guild
		if err = unmarshalEvent(e, &st); err == nil {
			if s.StateEnabled {
				s.State.GuildAdd(st)
			}
			if s.OnGuildCreate != nil {
				s.OnGuildUpdate(s, st)
			}
		}
		if s.OnGuildUpdate != nil {
			return
		}
	case "GUILD_DELETE":
		if !s.StateEnabled && s.OnGuildDelete == nil {
			break
		}
		var st *Guild
		if err = unmarshalEvent(e, &st); err == nil {
			if s.StateEnabled {
				s.State.GuildRemove(st)
			}
			if s.OnGuildDelete != nil {
				s.OnGuildDelete(s, st)
			}
		}
		if s.OnGuildDelete != nil {
			return
		}
	case "GUILD_MEMBER_ADD":
		if !s.StateEnabled && s.OnGuildMemberAdd == nil {
			break
		}
		var st *Member
		if err = unmarshalEvent(e, &st); err == nil {
			if s.StateEnabled {
				s.State.MemberAdd(st)
			}
			if s.OnGuildMemberAdd != nil {
				s.OnGuildMemberAdd(s, st)
			}
		}
		if s.OnGuildMemberAdd != nil {
			return
		}
	case "GUILD_MEMBER_REMOVE":
		if !s.StateEnabled && s.OnGuildMemberRemove == nil {
			break
		}
		var st *Member
		if err = unmarshalEvent(e, &st); err == nil {
			if s.StateEnabled {
				s.State.MemberRemove(st)
			}
			if s.OnGuildMemberRemove != nil {
				s.OnGuildMemberRemove(s, st)
			}
		}
		if s.OnGuildMemberRemove != nil {
			return
		}
	case "GUILD_MEMBER_UPDATE":
		if !s.StateEnabled && s.OnGuildMemberUpdate == nil {
			break
		}
		var st *Member
		if err = unmarshalEvent(e, &st); err == nil {
			if s.StateEnabled {
				s.State.MemberAdd(st)
			}
			if s.OnGuildMemberUpdate != nil {
				s.OnGuildMemberUpdate(s, st)
			}
		}
		if s.OnGuildMemberUpdate != nil {
			return
		}
	case "GUILD_ROLE_CREATE":
		if s.OnGuildRoleCreate != nil {
			var st *GuildRole
			if err = unmarshalEvent(e, &st); err == nil {
				s.OnGuildRoleCreate(s, st)
			}
			return
		}
	case "GUILD_ROLE_UPDATE":
		if s.OnGuildRoleUpdate != nil {
			var st *GuildRole
			if err = unmarshalEvent(e, &st); err == nil {
				s.OnGuildRoleUpdate(s, st)
			}
			return
		}
	case "GUILD_ROLE_DELETE":
		if s.OnGuildRoleDelete != nil {
			var st *GuildRoleDelete
			if err = unmarshalEvent(e, &st); err == nil {
				s.OnGuildRoleDelete(s, st)
			}
			return
		}
	case "GUILD_INTEGRATIONS_UPDATE":
		if s.OnGuildIntegrationsUpdate != nil {
			var st *GuildIntegrationsUpdate
			if err = unmarshalEvent(e, &st); err == nil {
				s.OnGuildIntegrationsUpdate(s, st)
			}
			return
		}
	case "GUILD_BAN_ADD":
		if s.OnGuildBanAdd != nil {
			var st *GuildBan
			if err = unmarshalEvent(e, &st); err == nil {
				s.OnGuildBanAdd(s, st)
			}
			return
		}
	case "GUILD_BAN_REMOVE":
		if s.OnGuildBanRemove != nil {
			var st *GuildBan
			if err = unmarshalEvent(e, &st); err == nil {
				s.OnGuildBanRemove(s, st)
			}
			return
		}
	case "GUILD_EMOJIS_UPDATE":
		if !s.StateEnabled && s.OnGuildEmojisUpdate == nil {
			break
		}
		var st *GuildEmojisUpdate
		if err = unmarshalEvent(e, &st); err == nil {
			if s.StateEnabled {
				s.State.EmojisAdd(st.GuildID, st.Emojis)
			}
			if s.OnGuildEmojisUpdate != nil {
				s.OnGuildEmojisUpdate(s, st)
			}
		}
		if s.OnGuildEmojisUpdate != nil {
			return
		}
	case "USER_SETTINGS_UPDATE":
		if s.OnUserSettingsUpdate != nil {
			var st map[string]interface{}
			if err = unmarshalEvent(e, &st); err == nil {
				s.OnUserSettingsUpdate(s, st)
			}
			return
		}
	default:
		fmt.Println("UNKNOWN EVENT: ", e.Type)
		printJSON(message)
	}

	// if still here, send to generic OnEvent
	if s.OnEvent != nil {
		s.OnEvent(s, e)
		return
	}

	return
}

// Heartbeat sends regular heartbeats to Discord so it knows the client
// is still connected.  If you do not send these heartbeats Discord will
// disconnect the websocket connection after a few seconds.
func (s *Session) Heartbeat(i time.Duration) {

	if s.wsConn == nil {
		fmt.Println("No websocket connection exists.")
		return // TODO need to return/log an error.
	}

	// Make sure Heartbeat is not already running
	s.heartbeatLock.Lock()
	if s.heartbeatChan != nil {
		s.heartbeatLock.Unlock()
		return
	}
	s.heartbeatChan = make(chan struct{})
	s.heartbeatLock.Unlock()

	// this is ugly.
	defer func() {
		if s.heartbeatChan == nil {
			return
		}
		select {
		case <-s.heartbeatChan:
			break
		default:
			close(s.heartbeatChan)
		}
		s.listenChan = nil
	}()

	// send first heartbeat immediately because lag could put the
	// first heartbeat outside the required heartbeat interval window
	ticker := time.NewTicker(i * time.Millisecond)
	for {

		err := s.wsConn.WriteJSON(map[string]int{
			"op": 1,
			"d":  int(time.Now().Unix()),
		})
		if err != nil {
			fmt.Println("error sending data heartbeat:", err)
			s.DataReady = false
			return // TODO log error?
		}
		s.DataReady = true

		select {
		case <-ticker.C:
		case <-s.heartbeatChan:
			return
		}
	}
}

// ------------------------------------------------------------------------------------------------
// Code related to voice connections that initiate over the data websocket
// ------------------------------------------------------------------------------------------------

// A VoiceServerUpdate stores the data received during the Voice Server Update
// data websocket event. This data is used during the initial Voice Channel
// join handshaking.
type VoiceServerUpdate struct {
	Token    string `json:"token"`
	GuildID  string `json:"guild_id"`
	Endpoint string `json:"endpoint"`
}

type voiceChannelJoinData struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	SelfMute  bool   `json:"self_mute"`
	SelfDeaf  bool   `json:"self_deaf"`
}

type voiceChannelJoinOp struct {
	Op   int                  `json:"op"`
	Data voiceChannelJoinData `json:"d"`
}

// ChannelVoiceJoin joins the session user to a voice channel. After calling
// this func please monitor the Session.Voice.Ready bool to determine when
// it is ready and able to send/receive audio, that should happen quickly.
//
//    gID   : Guild ID of the channel to join.
//    cID   : Channel ID of the channel to join.
//    mute  : If true, you will be set to muted upon joining.
//    deaf  : If true, you will be set to deafened upon joining.
func (s *Session) ChannelVoiceJoin(gID, cID string, mute, deaf bool) (err error) {

	if s.wsConn == nil {
		return fmt.Errorf("no websocket connection exists.")
	}

	// Create new voice{} struct if one does not exist.
	// If you create this prior to calling this func then you can manually
	// set some variables if needed, such as to enable debugging.
	if s.Voice == nil {
		s.Voice = &Voice{}
	}
	// TODO : Determine how to properly change channels and change guild
	// and channel when you are already connected to an existing channel.

	// Send the request to Discord that we want to join the voice channel
	data := voiceChannelJoinOp{4, voiceChannelJoinData{gID, cID, mute, deaf}}
	err = s.wsConn.WriteJSON(data)
	if err != nil {
		return
	}

	// Store gID and cID for later use
	s.Voice.guildID = gID
	s.Voice.channelID = cID

	return
}

// onVoiceStateUpdate handles Voice State Update events on the data
// websocket.  This comes immediately after the call to VoiceChannelJoin
// for the session user.
func (s *Session) onVoiceStateUpdate(st *VoiceState) {

	// If s.Voice is nil, we must not have even requested to join
	// a voice channel yet, so this shouldn't be processed.
	if s.Voice == nil {
		return
	}

	// Need to have this happen at login and store it in the Session
	// TODO : This should be done upon connecting to Discord, or
	// be moved to a small helper function
	self, err := s.User("@me") // TODO: move to Login/New
	if err != nil {
		fmt.Println(err)
		return
	}

	// This event comes for all users, if it's not for the session
	// user just ignore it.
	// TODO Move this IF to the event() func
	if st.UserID != self.ID {
		return
	}

	// Store the SessionID for later use.
	s.Voice.userID = self.ID // TODO: Review
	s.Voice.sessionID = st.SessionID
}

// onVoiceServerUpdate handles the Voice Server Update data websocket event.
// This event tells us the information needed to open a voice websocket
// connection and should happen after the VOICE_STATE event.
func (s *Session) onVoiceServerUpdate(st *VoiceServerUpdate) {

	// This shouldn't ever be the case, I don't think.
	if s.Voice == nil {
		return
	}

	// Store values for later use
	s.Voice.token = st.Token
	s.Voice.endpoint = st.Endpoint
	s.Voice.guildID = st.GuildID

	// We now have enough information to open a voice websocket conenction
	// so, that's what the next call does.
	err := s.Voice.Open()
	if err != nil {
		fmt.Println("onVoiceServerUpdate Voice.Open error: ", err)
		// TODO better logging
	}
}
