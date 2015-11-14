/******************************************************************************
 * A Discord API for Golang.
 * See discord.go for more information.
 *
 * This file contains low level functions for interacting
 * with the Discord Websocket interface.
 */

package discordgo

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

// Basic struct for all Websocket Event messages
type Event struct {
	Type      string `json:"t"`
	State     int    `json:"s"`
	Operation int    `json:"o"`
	Direction int    `json:"dir"`
	//Direction of command, 0-received, 1-sent -- thanks Xackery/discord

	RawData json.RawMessage `json:"d"`
	Session *Session
}

// The Ready Event given after initial connection
type Ready struct {
	Version           int           `json:"v"`
	SessionID         string        `json:"session_id"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	User              User          `json:"user"`
	ReadState         []ReadState
	PrivateChannels   []PrivateChannel
	Guilds            []Guild
}

// ReadState might need to move? Gives me the read status
// of all my channels when first connecting. I think :)
type ReadState struct {
	MentionCount  int
	LastMessageID int `json:"last_message_id,string"`
	ID            int `json:"id,string"`
}

type TypingStart struct {
	UserId    int `json:"user_id,string"`
	ChannelId int `json:"channel_id,string"`
	Timestamp int `json:"timestamp"`
}

type PresenceUpdate struct {
	User    User     `json:"user"`
	Status  string   `json:"status"`
	Roles   []string `json:"roles"` // TODO: Should be ints, see below
	GuildId int      `json:"guild_id,string"`
	GameId  int      `json:"game_id"`
}

//Roles   []string `json:"roles"` // TODO: Should be ints, see below
// Above "Roles" should be an array of ints
// TODO: Figure out how to make it be one.
/*
	{
		"roles": [
			"89544728336416768",
			"110429733396676608"
		],
	}
*/

type MessageAck struct {
	MessageId int `json:"message_id,string"`
	ChannelId int `json:"channel_id,string"`
}

type MessageDelete struct {
	Id        int `json:"id,string"`
	ChannelId int `json:"channel_id,string"`
} // so much like MessageAck..

type GuildIntegrationsUpdate struct {
	GuildId int `json:"guild_id,string"`
}

type GuildRoleUpdate struct {
	Role    Role `json:"role"`
	GuildId int  `json:"guild_id,int"`
}

// Open a websocket connection to Discord
func (s *Session) Open() (err error) {

	// Get the gateway to use for the Websocket connection
	g, err := s.Gateway()

	// TODO: See if there's a use for the http response.
	// conn, response, err := websocket.DefaultDialer.Dial(session.Gateway, nil)
	s.wsConn, _, err = websocket.DefaultDialer.Dial(g, nil)
	return
}

// maybe this is SendOrigin? not sure the right name here
// also bson.M vs string interface map?  Read about
// how to send JSON the right way.
func (s *Session) Handshake() (err error) {

	err = s.wsConn.WriteJSON(map[string]interface{}{
		"op": 2,
		"d": map[string]interface{}{
			"v":     3,
			"token": s.Token,
			"properties": map[string]string{
				"$os":               runtime.GOOS,
				"$browser":          "Discordgo",
				"$device":           "Discordgo",
				"$referer":          "",
				"$referring_domain": "",
			},
		},
	})

	return
}

func (s *Session) UpdateStatus(idleSince, gameId string) (err error) {

	err = s.wsConn.WriteJSON(map[string]interface{}{
		"op": 2,
		"d": map[string]interface{}{
			"idle_since": idleSince,
			"game_id":    gameId,
		},
	})

	return
}

// TODO: need a channel or something to communicate
// to this so I can tell it to stop listening
func (s *Session) Listen() (err error) {

	if s.wsConn == nil {
		fmt.Println("No websocket connection exists.")
		return // need to return an error.
	}

	for { // s.wsConn != nil { // need a cleaner way to exit?  this doesn't acheive anything.
		messageType, message, err := s.wsConn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		go s.event(messageType, message)
	}

	return
}

// Not sure how needed this is and where it would be best to call it.
// somewhere.
func (s *Session) Close() {
	s.wsConn.Close()
}

// Front line handler for all Websocket Events.  Determines the
// event type and passes the message along to the next handler.
func (s *Session) event(messageType int, message []byte) (err error) {

	if s.Debug {
		printJSON(message)
	}

	var e Event
	if err := json.Unmarshal(message, &e); err != nil {
		return err
	}

	switch e.Type {

	case "READY":
		if s.OnReady != nil {
			var st Ready
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logging
				return err
			}
			s.OnReady(s, st)
			return
		}
	case "VOICE_STATE_UPDATE":
		if s.OnVoiceStateUpdate != nil {
			var st VoiceState
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logging
				return err
			}
			s.OnVoiceStateUpdate(s, st)
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
		if s.OnGuildRoleUpdate != nil {
			var st GuildRoleUpdate
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				fmt.Println(e.Type, err)
				printJSON(e.RawData) // TODO: Better error logginEventg
				return err
			}
			s.OnGuildRoleUpdate(s, st)
			return
		}
		/*
			case "GUILD_ROLE_DELETE":
				if s.OnGuildRoleDelete != nil {
				}
		*/
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
		// learn the log package
		// log.print type and JSON data
	}

	// if still here, send to generic OnEvent
	if s.OnEvent != nil {
		s.OnEvent(s, e)
	}

	return
}

// This heartbeat is sent to keep the Websocket conenction
// to Discord alive. If not sent, Discord will close the
// connection.
func (s *Session) Heartbeat(i time.Duration) {

	if s.wsConn == nil {
		fmt.Println("No websocket connection exists.")
		return // need to return an error.
	}

	ticker := time.NewTicker(i * time.Millisecond)
	for range ticker.C {
		timestamp := int(time.Now().Unix())
		err := s.wsConn.WriteJSON(map[string]int{
			"op": 1,
			"d":  timestamp,
		})
		if err != nil {
			return // log error?
		}
	}
}
