// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains low level functions for interacting with the Discord
// data websocket interface.

package discordgo

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

// Open opens a websocket connection to Discord.
func (s *Session) Open() (err error) {

	// Get the gateway to use for the Websocket connection
	g, err := s.Gateway()

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

	data := handshakeOp{2, handshakeData{3, s.Token, handshakeProperties{runtime.GOOS, "DiscordGo v" + VERSION, "", "", ""}}}
	err = s.wsConn.WriteJSON(data)
	return
}

type updateStatusGame struct {
	Name *string `json:"name"`
}

type updateStatusData struct {
	IdleSince *int             `json:"idle_since"`
	Game      updateStatusGame `json:"game"`
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
		usd.Game.Name = &game
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

	for {
		messageType, message, err := s.wsConn.ReadMessage()
		if err != nil {
			fmt.Println("Websocket Listen Error", err)
			// TODO Log error
			break
		}
		go s.event(messageType, message)
	}

	return
}

// Not sure how needed this is and where it would be best to call it.
// somewhere.

// Close closes the connection to the websocket.
func (s *Session) Close() {
	s.wsConn.Close()
}

// Front line handler for all Websocket Events.  Determines the
// event type and passes the message along to the next handler.

// event is the front line handler for all events.  This needs to be
// broken up into smaller functions to be more idiomatic Go.
func (s *Session) event(messageType int, message []byte) (err error) {

	if s.Debug {
		printJSON(message)
	}

	var e Event
	if err := json.Unmarshal(message, &e); err != nil {
		fmt.Println(err)
		return err
	}

	switch e.Type {

	case "READY":
		var st Ready
		if err := json.Unmarshal(e.RawData, &st); err != nil {
			fmt.Println(e.Type, err)
			printJSON(e.RawData) // TODO: Better error logging
			return err
		}
		if s.OnReady != nil {
			s.OnReady(s, st)
			return
		}
		go s.Heartbeat(st.HeartbeatInterval)
	case "VOICE_SERVER_UPDATE":
		// TEMP CODE FOR TESTING VOICE
		var st VoiceServerUpdate
		if err := json.Unmarshal(e.RawData, &st); err != nil {
			fmt.Println(e.Type, err)
			printJSON(e.RawData) // TODO: Better error logging
			return err
		}
		s.onVoiceServerUpdate(st)
		return
	case "VOICE_STATE_UPDATE":
		// TEMP CODE FOR TESTING VOICE
		var st VoiceState
		if err := json.Unmarshal(e.RawData, &st); err != nil {
			fmt.Println(e.Type, err)
			printJSON(e.RawData) // TODO: Better error logging
			return err
		}
		s.onVoiceStateUpdate(st)
		return
	case "USER_UPDATE":
		if s.OnUserUpdate != nil {
			var st User
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logging
				return err
			}
			s.OnUserUpdate(s, st)
			return
		}
	case "PRESENCE_UPDATE":
		if s.OnPresenceUpdate != nil {
			var st PresenceUpdate
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logging
				return err
			}
			s.OnPresenceUpdate(s, st)
			return
		}
	case "TYPING_START":
		if s.OnTypingStart != nil {
			var st TypingStart
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logging
				return err
			}
			s.OnTypingStart(s, st)
			return
		}
		/* // Never seen this come in but saw it in another Library.
		case "MESSAGE_ACK":
			if s.OnMessageAck != nil {
			}
		*/
	case "MESSAGE_CREATE":
		if s.OnMessageCreate != nil {
			var st Message
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logging
				return err
			}
			s.OnMessageCreate(s, st)
			return
		}
	case "MESSAGE_UPDATE":
		if s.OnMessageUpdate != nil {
			var st Message
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logging
				return err
			}
			s.OnMessageUpdate(s, st)
			return
		}
	case "MESSAGE_DELETE":
		if s.OnMessageDelete != nil {
			var st MessageDelete
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logging
				return err
			}
			s.OnMessageDelete(s, st)
			return
		}
	case "MESSAGE_ACK":
		if s.OnMessageAck != nil {
			var st MessageAck
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logging
				return err
			}
			s.OnMessageAck(s, st)
			return
		}
	case "CHANNEL_CREATE":
		if s.OnChannelCreate != nil {
			var st Channel
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnChannelCreate(s, st)
			return
		}
	case "CHANNEL_UPDATE":
		if s.OnChannelUpdate != nil {
			var st Channel
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnChannelUpdate(s, st)
			return
		}
	case "CHANNEL_DELETE":
		if s.OnChannelDelete != nil {
			var st Channel
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnChannelDelete(s, st)
			return
		}
	case "GUILD_CREATE":
		if s.OnGuildCreate != nil {
			var st Guild
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnGuildCreate(s, st)
			return
		}
	case "GUILD_UPDATE":
		if s.OnGuildCreate != nil {
			var st Guild
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnGuildUpdate(s, st)
			return
		}
	case "GUILD_DELETE":
		if s.OnGuildDelete != nil {
			var st Guild
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnGuildDelete(s, st)
			return
		}
	case "GUILD_MEMBER_ADD":
		if s.OnGuildMemberAdd != nil {
			var st Member
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnGuildMemberAdd(s, st)
			return
		}
	case "GUILD_MEMBER_REMOVE":
		if s.OnGuildMemberRemove != nil {
			var st Member
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnGuildMemberRemove(s, st)
			return
		}
	case "GUILD_MEMBER_UPDATE":
		if s.OnGuildMemberUpdate != nil {
			var st Member
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnGuildMemberUpdate(s, st)
			return
		}
	case "GUILD_ROLE_CREATE":
		if s.OnGuildRoleCreate != nil {
			var st GuildRole
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnGuildRoleCreate(s, st)
			return
		}
	case "GUILD_ROLE_UPDATE":
		if s.OnGuildRoleUpdate != nil {
			var st GuildRole
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnGuildRoleUpdate(s, st)
			return
		}
	case "GUILD_ROLE_DELETE":
		if s.OnGuildRoleDelete != nil {
			var st GuildRoleDelete
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnGuildRoleDelete(s, st)
			return
		}
	case "GUILD_INTEGRATIONS_UPDATE":
		if s.OnGuildIntegrationsUpdate != nil {
			var st GuildIntegrationsUpdate
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnGuildIntegrationsUpdate(s, st)
			return
		}
	default:
		fmt.Println("UNKNOWN EVENT: ", e.Type)
		// TODO learn the log package
		// log.print type and JSON data
	}

	// if still here, send to generic OnEvent
	if s.OnEvent != nil {
		s.OnEvent(s, e)
		return
	}

	return
}

// This heartbeat is sent to keep the Websocket conenction
// to Discord alive. If not sent, Discord will close the
// connection.

// Heartbeat sends regular heartbeats to Discord so it knows the client
// is still connected.  If you do not send these heartbeats Discord will
// disconnect the websocket connection after a few seconds.
func (s *Session) Heartbeat(i time.Duration) {

	if s.wsConn == nil {
		fmt.Println("No websocket connection exists.")
		return // TODO need to return an error.
	}

	// send first heartbeat immediately because lag could put the
	// first heartbeat outside the required heartbeat interval window
	ticker := time.NewTicker(i * time.Millisecond)
	for {
		timestamp := int(time.Now().Unix())
		err := s.wsConn.WriteJSON(map[string]int{
			"op": 1,
			"d":  timestamp,
		})
		if err != nil {
			fmt.Println("error sending data heartbeat:", err)
			s.DataReady = false
			return // TODO log error?
		}
		s.DataReady = true
		<-ticker.C
	}
}

// Everything below is experimental Voice support code
// all of it will get changed and moved around.

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

// VoiceChannelJoin joins the authenticated session user to
// a voice channel.  All the voice magic starts with this.
func (s *Session) VoiceChannelJoin(guildID, channelID string) (err error) {

	if s.wsConn == nil {
		fmt.Println("error: no websocket connection exists.")
		return // TODO return error
	}

	data := voiceChannelJoinOp{4, voiceChannelJoinData{guildID, channelID, false, false}}
	err = s.wsConn.WriteJSON(data)
	if err != nil {
		return
	}

	// Probably will be removed later.
	s.VGuildID = guildID
	s.VChannelID = channelID

	return
}

// onVoiceStateUpdate handles Voice State Update events on the data
// websocket.  This comes immediately after the call to VoiceChannelJoin
// for the authenticated session user.  This block is experimental
// code and will be chaned in the future.
func (s *Session) onVoiceStateUpdate(st VoiceState) {

	// Need to have this happen at login and store it in the Session
	self, err := s.User("@me") // TODO: move to Login/New
	if err != nil {
		fmt.Println(err)
		return
	}

	// This event comes for all users, if it's not for the session
	// user just ignore it.
	if st.UserID != self.ID {
		return
	}

	// Store the SessionID. Used later.
	s.VSessionID = st.SessionID
}

// onVoiceServerUpdate handles the Voice Server Update data websocket event.
// This will later be exposed but is only for experimental use now.
func (s *Session) onVoiceServerUpdate(st VoiceServerUpdate) {

	// Store all the values.  They are used later.
	// GuildID is probably not needed and may be dropped.
	s.VToken = st.Token
	s.VEndpoint = st.Endpoint
	s.VGuildID = st.GuildID

	// We now have enough information to open a voice websocket conenction
	// so, that's what the next call does.
	s.VoiceOpenWS()
}
