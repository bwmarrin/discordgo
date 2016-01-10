// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains experimental functions for interacting with the Discord
// Voice websocket and UDP connections.
//
// EVERYTHING in this file is very experimental and will change.

package discordgo

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ------------------------------------------------------------------------------------------------
// Code related to both Voice Websocket and UDP connections.
// ------------------------------------------------------------------------------------------------

// A Voice struct holds all data and functions related to Discord Voice support.
type voice struct {
	sync.Mutex
	Ready   bool
	Debug   bool
	Chan    chan struct{}
	UDPConn *net.UDPConn // exported for dgvoice, may change.

	wsConn *websocket.Conn

	sessionID string
	token     string
	endpoint  string
	guildID   string
	channelID string
	userID    string
	OP2       *voiceOP2 // exported for dgvoice, may change.
}

// ------------------------------------------------------------------------------------------------
// Code related to the Voice websocket connection
// ------------------------------------------------------------------------------------------------

// A voiceOP2 stores the data for voice operation 2 websocket events
// which is sort of like the voice READY packet
type voiceOP2 struct {
	SSRC              uint32        `json:"ssrc"`
	Port              int           `json:"port"`
	Modes             []string      `json:"modes"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
}

type voiceHandshakeData struct {
	ServerID  string `json:"server_id"`
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Token     string `json:"token"`
}

type voiceHandshakeOp struct {
	Op   int                `json:"op"` // Always 0
	Data voiceHandshakeData `json:"d"`
}

// Open opens a voice connection.  This should be called
// after VoiceChannelJoin is used and the data VOICE websocket events
// are captured.
func (v *voice) Open() (err error) {

	// TODO: How do we handle changing channels?
	// Don't open a websocket if one is already open
	if v.wsConn != nil {
		return
	}

	// Connect to Voice Websocket
	vg := fmt.Sprintf("wss://%s", strings.TrimSuffix(v.endpoint, ":80"))
	v.wsConn, _, err = websocket.DefaultDialer.Dial(vg, nil)
	if err != nil {
		fmt.Println("VOICE cannot open websocket:", err)
		return
	}

	data := voiceHandshakeOp{0, voiceHandshakeData{v.guildID, v.userID, v.sessionID, v.token}}

	err = v.wsConn.WriteJSON(data)
	if err != nil {
		fmt.Println("VOICE ERROR sending init packet:", err)
		return
	}

	// Start a listening for voice websocket events
	// TODO add a check here to make sure Listen worked by monitoring
	// a chan or bool?
	//	go vws.Listen()
	go v.wsListen()
	return
}

// Close closes the voice connection
func (v *voice) Close() {

	if v.UDPConn != nil {
		v.UDPConn.Close()
	}

	if v.wsConn != nil {
		v.wsConn.Close()
	}
}

// wsListen listens on the voice websocket for messages and passes them
// to the voice event handler.  This is automaticly called by the WS.Open
// func when needed.
func (v *voice) wsListen() {

	for {
		messageType, message, err := v.wsConn.ReadMessage()
		if err != nil {
			// TODO: Handle this problem better.
			// TODO: needs proper logging
			fmt.Println("Voice Listen Error:", err)
			break
		}

		// Pass received message to voice event handler
		go v.wsEvent(messageType, message)
	}

	return
}

// wsEvent handles any voice websocket events. This is only called by the
// wsListen() function.
func (v *voice) wsEvent(messageType int, message []byte) {

	if v.Debug {
		fmt.Println("wsEvent received: ", messageType)
		printJSON(message)
	}

	var e Event
	if err := json.Unmarshal(message, &e); err != nil {
		fmt.Println("wsEvent Unmarshall error: ", err)
		return
	}

	switch e.Operation {

	case 2: // READY

		v.OP2 = &voiceOP2{}
		if err := json.Unmarshal(e.RawData, v.OP2); err != nil {
			fmt.Println("voiceWS.onEvent OP2 Unmarshall error: ", err)
			printJSON(e.RawData) // TODO: Better error logging
			return
		}

		// Start the voice websocket heartbeat to keep the connection alive
		go v.wsHeartbeat(v.OP2.HeartbeatInterval)
		// TODO monitor a chan/bool to verify this was successful

		// We now have enough data to start the UDP connection
		v.udpOpen()

		return
	case 3: // HEARTBEAT response
		// add code to use this to track latency?
		return
	case 4:
		// TODO
	default:
		fmt.Println("UNKNOWN VOICE OP: ", e.Operation)
		printJSON(e.RawData)
	}

	return
}

type voiceHeartbeatOp struct {
	Op   int `json:"op"` // Always 3
	Data int `json:"d"`
}

// wsHeartbeat sends regular heartbeats to voice Discord so it knows the client
// is still connected.  If you do not send these heartbeats Discord will
// disconnect the websocket connection after a few seconds.
func (v *voice) wsHeartbeat(i time.Duration) {

	ticker := time.NewTicker(i * time.Millisecond)
	for {
		err := v.wsConn.WriteJSON(voiceHeartbeatOp{3, int(time.Now().Unix())})
		if err != nil {
			v.Ready = false
			fmt.Println("wsHeartbeat send error: ", err)
			return // TODO better logging
		}
		<-ticker.C
	}
}

type voiceSpeakingData struct {
	Speaking bool `json:"speaking"`
	Delay    int  `json:"delay"`
}

type voiceSpeakingOp struct {
	Op   int               `json:"op"` // Always 5
	Data voiceSpeakingData `json:"d"`
}

// Speaking sends a speaking notification to Discord over the voice websocket.
// This must be sent as true prior to sending audio and should be set to false
// once finished sending audio.
//  b  : Send true if speaking, false if not.
func (v *voice) Speaking(b bool) (err error) {

	if v.wsConn == nil {
		return fmt.Errorf("No Voice websocket.")
	}

	data := voiceSpeakingOp{5, voiceSpeakingData{b, 0}}
	err = v.wsConn.WriteJSON(data)
	if err != nil {
		fmt.Println("Speaking() write json error:", err)
		return
	}

	return
}

// ------------------------------------------------------------------------------------------------
// Code related to the Voice UDP connection
// ------------------------------------------------------------------------------------------------

type voiceUDPData struct {
	Address string `json:"address"` // Public IP of machine running this code
	Port    uint16 `json:"port"`    // UDP Port of machine running this code
	Mode    string `json:"mode"`    // plain or ?  (plain or encrypted)
}

type voiceUDPD struct {
	Protocol string       `json:"protocol"` // Always "udp" ?
	Data     voiceUDPData `json:"data"`
}

type voiceUDPOp struct {
	Op   int       `json:"op"` // Always 1
	Data voiceUDPD `json:"d"`
}

// udpOpen opens a UDP connect to the voice server and completes the
// initial required handshake.  This connect is left open in the session
// and can be used to send or receive audio.  This should only be called
// from voice.wsEvent OP2
func (v *voice) udpOpen() (err error) {

	host := fmt.Sprintf("%s:%d", strings.TrimSuffix(v.endpoint, ":80"), v.OP2.Port)
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		fmt.Println("udpOpen() resolve addr error: ", err)
		// TODO better logging
		return
	}

	v.UDPConn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("udpOpen() dial udp error: ", err)
		// TODO better logging
		return
	}

	// Create a 70 byte array and put the SSRC code from the Op 2 Voice event
	// into it.  Then send that over the UDP connection to Discord
	sb := make([]byte, 70)
	binary.BigEndian.PutUint32(sb, v.OP2.SSRC)
	v.UDPConn.Write(sb)

	// Create a 70 byte array and listen for the initial handshake response
	// from Discord.  Once we get it parse the IP and PORT information out
	// of the response.  This should be our public IP and PORT as Discord
	// saw us.
	rb := make([]byte, 70)
	rlen, _, err := v.UDPConn.ReadFromUDP(rb)
	if rlen < 70 {
		fmt.Println("Voice RLEN should be 70 but isn't")
	}

	// Loop over position 4 though 20 to grab the IP address
	// Should never be beyond position 20.
	var ip string
	for i := 4; i < 20; i++ {
		if rb[i] == 0 {
			break
		}
		ip += string(rb[i])
	}

	// Grab port from postion 68 and 69
	port := binary.LittleEndian.Uint16(rb[68:70])

	// Take the parsed data from above and send it back to Discord
	// to finalize the UDP handshake.
	data := voiceUDPOp{1, voiceUDPD{"udp", voiceUDPData{ip, port, "plain"}}}

	err = v.wsConn.WriteJSON(data)
	if err != nil {
		fmt.Println("udpOpen write json error:", err)
		return
	}

	v.Ready = true
	return
}
