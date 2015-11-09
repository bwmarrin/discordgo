/******************************************************************************
 * Discordgo by Bruce Marriner <bruce@sqls.net>
 * A Discord API for Golang.
 * See discord.go for more information.
 *
 * This file contains functions low level functions for interacting
 * with the Discord Websocket interface.
 */

package discordgo

import (
	"encoding/json"
	"fmt"
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
func Open(session *Session) (conn *websocket.Conn, err error) {

	// TODO: See if there's a use for the http response.
	//conn, response, err := websocket.DefaultDialer.Dial(session.Gateway, nil)
	conn, _, err = websocket.DefaultDialer.Dial(session.Gateway, nil)
	if err != nil {
		return
	}

	return
}

// maybe this is SendOrigin? not sure the right name here
// also bson.M vs string interface map?  Read about
// how to send JSON the right way.
func Handshake(conn *websocket.Conn, token string) (err error) {

	err = conn.WriteJSON(map[string]interface{}{
		"op": 2,
		"d": map[string]interface{}{
			"v":     3,
			"token": token,
			"properties": map[string]string{
				"$os":               "linux", // get from os package
				"$browser":          "Discordgo",
				"$device":           "Discordgo",
				"$referer":          "",
				"$referring_domain": "",
			},
		},
	})

	return
}

func UpdateStatus(conn *websocket.Conn, idleSince, gameId string) (err error) {

	err = conn.WriteJSON(map[string]interface{}{
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
func Listen(conn *websocket.Conn) (err error) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		go event(conn, messageType, message)
	}

	return
}

// Not sure how needed this is and where it would be best to call it.
// somewhere.
func Close(conn *websocket.Conn) {
	conn.Close()
}

// Front line handler for all Websocket Events.  Determines the
// event type and passes the message along to the next handler.
func event(conn *websocket.Conn, messageType int, message []byte) {

	//printJSON(message) // TODO: wrap in debug if statement

	var event Event
	err := json.Unmarshal(message, &event)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch event.Type {

	case "READY":
		ready(conn, &event)
	case "TYPING_START":
		// do stuff
	case "MESSAGE_CREATE":
		// do stuff
	case "MESSAGE_ACK":
		// do stuff
	case "MESSAGE_UPDATE":
		// do stuff
	case "MESSAGE_DELETE":
		// do stuff
	case "PRESENCE_UPDATE":
		// do stuff
	case "CHANNEL_CREATE":
		// do stuff
	case "CHANNEL_UPDATE":
		// do stuff
	case "CHANNEL_DELETE":
		// do stuff
	case "GUILD_CREATE":
		// do stuff
	case "GUILD_DELETE":
		// do stuff
	case "GUILD_MEMBER_ADD":
		// do stuff
	case "GUILD_MEMBER_REMOVE": // which is it.
		// do stuff
	case "GUILD_MEMBER_DELETE":
		// do stuff
	case "GUILD_MEMBER_UPDATE":
		// do stuff
	case "GUILD_ROLE_CREATE":
		// do stuff
	case "GUILD_ROLE_DELETE":
		// do stuff
	case "GUILD_INTEGRATIONS_UPDATE":
		// do stuff

	default:
		fmt.Println("UNKNOWN EVENT: ", event.Type)
		// learn the log package
		// log.print type and JSON data
	}

}

// handles the READY Websocket Event from Discord
// this is the motherload of detail provided at
// initial connection to the Websocket.
func ready(conn *websocket.Conn, event *Event) {

	var ready Ready
	err := json.Unmarshal(event.RawData, &ready)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ready)

	go heartbeat(conn, ready.HeartbeatInterval)

	// Start KeepAlive based on .

}

// This heartbeat is sent to keep the Websocket conenction
// to Discord alive. If not sent, Discord will close the
// connection.
func heartbeat(conn *websocket.Conn, interval time.Duration) {

	ticker := time.NewTicker(interval * time.Millisecond)
	for range ticker.C {
		timestamp := int(time.Now().Unix())
		conn.WriteJSON(map[string]int{
			"op": 1,
			"d":  timestamp,
		})
	}
}
