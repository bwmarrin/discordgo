// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains low level functions for interacting with the Discord
// data websocket interface.

package discordgo

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

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
	Compress   bool                `json:"compress"`
}

type handshakeOp struct {
	Op   int           `json:"op"`
	Data handshakeData `json:"d"`
}

// Open opens a websocket connection to Discord.
func (s *Session) Open() (err error) {
	s.Lock()
	defer func() {
		if err != nil {
			s.Unlock()
		}
	}()

	if s.wsConn != nil {
		err = errors.New("Web socket already opened.")
		return
	}

	// Get the gateway to use for the Websocket connection
	g, err := s.Gateway()
	if err != nil {
		return
	}

	header := http.Header{}
	header.Add("accept-encoding", "zlib")

	// TODO: See if there's a use for the http response.
	// conn, response, err := websocket.DefaultDialer.Dial(session.Gateway, nil)
	s.wsConn, _, err = websocket.DefaultDialer.Dial(g, header)
	if err != nil {
		return
	}

	err = s.wsConn.WriteJSON(handshakeOp{2, handshakeData{3, s.Token, handshakeProperties{runtime.GOOS, "Discordgo v" + VERSION, "", "", ""}, s.Compress}})
	if err != nil {
		return
	}

	// Create listening outside of listen, as it needs to happen inside the mutex
	// lock.
	s.listening = make(chan interface{})
	go s.listen(s.wsConn, s.listening)

	s.Unlock()

	if s.OnConnect != nil {
		s.OnConnect(s)
	}

	return
}

// Close closes a websocket and stops all listening/heartbeat goroutines.
// TODO: Add support for Voice WS/UDP connections
func (s *Session) Close() (err error) {
	s.Lock()

	s.DataReady = false

	if s.listening != nil {
		close(s.listening)
		s.listening = nil
	}

	if s.wsConn != nil {
		err = s.wsConn.Close()
		s.wsConn = nil
	}

	s.Unlock()

	if s.OnDisconnect != nil {
		s.OnDisconnect(s)
	}

	return
}

// listen polls the websocket connection for events, it will stop when
// the listening channel is closed, or an error occurs.
func (s *Session) listen(wsConn *websocket.Conn, listening <-chan interface{}) {
	for {
		messageType, message, err := wsConn.ReadMessage()
		if err != nil {
			// Detect if we have been closed manually. If a Close() has already
			// happened, the websocket we are listening on will be different to the
			// current session.
			s.RLock()
			sameConnection := s.wsConn == wsConn
			s.RUnlock()
			if sameConnection {
				// There has been an error reading, Close() the websocket so that
				// OnDisconnect is fired.
				err := s.Close()
				if err != nil {
					fmt.Println("error closing session connection: ", err)
				}

				// Attempt to reconnect, with expenonential backoff up to 10 minutes.
				if s.ShouldReconnectOnError {
					wait := time.Duration(1)
					for {
						if s.Open() == nil {
							return
						}
						<-time.After(wait * time.Second)
						wait *= 2
						if wait > 600 {
							wait = 600
						}
					}
				}
			}
			return
		}

		select {
		case <-listening:
			return
		default:
			go s.event(messageType, message)
		}
	}
}

type heartbeatOp struct {
	Op   int `json:"op"`
	Data int `json:"d"`
}

func (s *Session) sendHeartbeat(wsConn *websocket.Conn) error {
	return wsConn.WriteJSON(heartbeatOp{1, int(time.Now().Unix())})
}

// heartbeat sends regular heartbeats to Discord so it knows the client
// is still connected.  If you do not send these heartbeats Discord will
// disconnect the websocket connection after a few seconds.
func (s *Session) heartbeat(wsConn *websocket.Conn, listening <-chan interface{}, i time.Duration) {
	if listening == nil || wsConn == nil {
		return
	}

	s.Lock()
	s.DataReady = true
	s.Unlock()

	// Send first heartbeat immediately because lag could put the
	// first heartbeat outside the required heartbeat interval window.
	err := s.sendHeartbeat(wsConn)
	if err != nil {
		fmt.Println("Error sending initial heartbeat:", err)
		return
	}

	ticker := time.NewTicker(i * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			err := s.sendHeartbeat(wsConn)
			if err != nil {
				fmt.Println("Error sending heartbeat:", err)
				return
			}
		case <-listening:
			return
		}
	}
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
	s.RLock()
	defer s.RUnlock()
	if s.wsConn == nil {
		return errors.New("No websocket connection exists.")
	}

	var usd updateStatusData
	if idle > 0 {
		usd.IdleSince = &idle
	}
	if game != "" {
		usd.Game = &updateStatusGame{game}
	}

	err = s.wsConn.WriteJSON(updateStatusOp{3, usd})

	return
}

// Not sure how needed this is and where it would be best to call it.
// somewhere.

func unmarshalEvent(event *Event, i interface{}) (err error) {
	if err = unmarshal(event.RawData, i); err != nil {
		fmt.Println("Unable to unmarshal event data.")
		printEvent(event)
	}
	return
}

// Front line handler for all Websocket Events.  Determines the
// event type and passes the message along to the next handler.

// event is the front line handler for all events.  This needs to be
// broken up into smaller functions to be more idiomatic Go.
// Events will be handled by any implemented handler in Session.
// All unhandled events will then be handled by OnEvent.
func (s *Session) event(messageType int, message []byte) {
	var err error
	var reader io.Reader
	reader = bytes.NewBuffer(message)

	if messageType == 2 {
		z, err1 := zlib.NewReader(reader)
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		defer func() {
			err := z.Close()
			if err != nil {
				fmt.Println("error closing zlib:", err)
			}
		}()
		reader = z
	}

	var e *Event
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&e); err != nil {
		fmt.Println(err)
		return
	}

	if s.Debug {
		printEvent(e)
	}

	switch e.Type {
	case "READY":
		var st *Ready
		if err = unmarshalEvent(e, &st); err == nil {
			go s.heartbeat(s.wsConn, s.listening, st.HeartbeatInterval)
			if s.StateEnabled {
				err := s.State.OnReady(st)
				if err != nil {
					fmt.Println("error: ", err)
				}

			}
			if s.OnReady != nil {
				s.OnReady(s, st)
			}
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
				err := s.State.MessageAdd(st)
				if err != nil {
					fmt.Println("error :", err)
				}
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
				err := s.State.MessageAdd(st)
				if err != nil {
					fmt.Println("error :", err)
				}
			}
			if s.OnMessageUpdate != nil {
				s.OnMessageUpdate(s, st)
			}
		}
		return
	case "MESSAGE_DELETE":
		stateEnabled := s.StateEnabled && s.State.MaxMessageCount > 0
		if !stateEnabled && s.OnMessageDelete == nil {
			break
		}
		var st *Message
		if err = unmarshalEvent(e, &st); err == nil {
			if stateEnabled {
				err := s.State.MessageRemove(st)
				if err != nil {
					fmt.Println("error :", err)
				}
			}
			if s.OnMessageDelete != nil {
				s.OnMessageDelete(s, st)
			}
		}
		return
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
				err := s.State.ChannelAdd(st)
				if err != nil {
					fmt.Println("error :", err)
				}
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
				err := s.State.ChannelAdd(st)
				if err != nil {
					fmt.Println("error :", err)
				}
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
				err := s.State.ChannelRemove(st)
				if err != nil {
					fmt.Println("error :", err)
				}
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
				err := s.State.GuildAdd(st)
				if err != nil {
					fmt.Println("error :", err)
				}
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
				err := s.State.GuildAdd(st)
				if err != nil {
					fmt.Println("error :", err)
				}
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
				err := s.State.GuildRemove(st)
				if err != nil {
					fmt.Println("error :", err)
				}
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
				err := s.State.MemberAdd(st)
				if err != nil {
					fmt.Println("error :", err)
				}
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
				err := s.State.MemberRemove(st)
				if err != nil {
					fmt.Println("error :", err)
				}
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
				err := s.State.MemberAdd(st)
				if err != nil {
					fmt.Println("error :", err)
				}
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
				err := s.State.EmojisAdd(st.GuildID, st.Emojis)
				if err != nil {
					fmt.Println("error :", err)
				}
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
		fmt.Println("Unknown Event.")
		printEvent(e)
	}

	// if still here, send to generic OnEvent
	if s.OnEvent != nil {
		s.OnEvent(s, e)
		return
	}

	return
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
		return fmt.Errorf("no websocket connection exists")
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
