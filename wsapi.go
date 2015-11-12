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
	Session Session
}

// The Ready Event given after initial connection
type Ready struct {
	Version           int           `json:"v"`
	SessionID         string        `json:"session_id"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	User              User          `json:"user"`
	ReadState         []ReadState
	PrivateChannels   []PrivateChannel
	Servers           []Server
}

// ReadState might need to move? Gives me the read status
// of all my channels when first connecting. I think :)
type ReadState struct {
	MentionCount  int
	LastMessageID int `json:"last_message_id,string"`
	ID            int `json:"id,string"`
}

// Open a websocket connection to Discord
func Open(s *Session) (err error) {

	// TODO: See if there's a use for the http response.
	// conn, response, err := websocket.DefaultDialer.Dial(session.Gateway, nil)
	s.wsConn, _, err = websocket.DefaultDialer.Dial(s.Gateway, nil)
	return
}

// maybe this is SendOrigin? not sure the right name here
// also bson.M vs string interface map?  Read about
// how to send JSON the right way.
func Handshake(s *Session) (err error) {

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

func UpdateStatus(s *Session, idleSince, gameId string) (err error) {

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
func Listen(s *Session) (err error) {

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
		go event(s, messageType, message)
	}

	return
}

// Not sure how needed this is and where it would be best to call it.
// somewhere.
func Close(s *Session) {
	s.wsConn.Close()
}

// Front line handler for all Websocket Events.  Determines the
// event type and passes the message along to the next handler.
func event(s *Session, messageType int, message []byte) (err error) {

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
				return err
			}
			s.OnReady(s, st)
			return
		}
	case "TYPING_START":
		if s.OnTypingStart != nil {
		}
	case "MESSAGE_CREATE":
		if s.OnMessageCreate != nil {
			var st Message
			if err := json.Unmarshal(e.RawData, &st); err != nil {
				return err
			}
			s.OnMessageCreate(s, st)
			return
		}
	case "MESSAGE_ACK":
		if s.OnMessageAck != nil {
		}
	case "MESSAGE_UPDATE":
		if s.OnMessageUpdate != nil {
		}
	case "MESSAGE_DELETE":
		if s.OnMessageDelete != nil {
		}
	case "PRESENCE_UPDATE":
		if s.OnPresenceUpdate != nil {
		}
	case "CHANNEL_CREATE":
		if s.OnChannelCreate != nil {
		}
	case "CHANNEL_UPDATE":
		if s.OnChannelUpdate != nil {
		}
	case "CHANNEL_DELETE":
		if s.OnChannelDelete != nil {
		}
	case "GUILD_CREATE":
		if s.OnGuildCreate != nil {
		}
	case "GUILD_DELETE":
		if s.OnGuildDelete != nil {
		}
	case "GUILD_MEMBER_ADD":
		if s.OnGuildMemberAdd != nil {
		}
	case "GUILD_MEMBER_REMOVE": // which is it.
		if s.OnGuildMemberRemove != nil {
		}
	case "GUILD_MEMBER_DELETE":
		if s.OnGuildMemberDelete != nil {
		}
	case "GUILD_MEMBER_UPDATE":
		if s.OnGuildMemberUpdate != nil {
		}
	case "GUILD_ROLE_CREATE":
		if s.OnGuildRoleCreate != nil {
		}
	case "GUILD_ROLE_DELETE":
		if s.OnGuildRoleDelete != nil {
		}
	case "GUILD_INTEGRATIONS_UPDATE":
		if s.OnGuildIntegrationsUpdate != nil {
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
func Heartbeat(s *Session, i time.Duration) {

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
