// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
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
	"time"

	"github.com/gorilla/websocket"
)

// A VEvent is the initial structure for voice websocket events.  I think
// I can reuse the data websocket structure here.
type VEvent struct {
	Type      string          `json:"t"`
	State     int             `json:"s"`
	Operation int             `json:"op"`
	RawData   json.RawMessage `json:"d"`
}

// A VoiceOP2 stores the data for voice operation 2 websocket events
// which is sort of like the voice READY packet
type VoiceOP2 struct {
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

// VoiceOpenWS opens a voice websocket connection.  This should be called
// after VoiceChannelJoin is used and the data VOICE websocket events
// are captured.
func (s *Session) VoiceOpenWS() {

	// Don't open a socket if one is already open
	if s.VwsConn != nil {
		return
	}

	var self User
	var err error

	self, err = s.User("@me") // TODO: Move to @ login and store in session

	// Connect to Voice Websocket
	vg := fmt.Sprintf("wss://%s", strings.TrimSuffix(s.VEndpoint, ":80"))
	s.VwsConn, _, err = websocket.DefaultDialer.Dial(vg, nil)
	if err != nil {
		fmt.Println("VOICE cannot open websocket:", err)
	}

	json := voiceHandshakeOp{0, voiceHandshakeData{s.VGuildID, self.ID, s.VSessionID, s.VToken}}

	err = s.VwsConn.WriteJSON(json)
	if err != nil {
		fmt.Println("VOICE ERROR sending init packet:", err)
	}

	// Start a listening for voice websocket events
	go s.VoiceListen()
}

// Close closes the connection to the voice websocket.
func (s *Session) VoiceCloseWS() {
	s.VwsConn.Close()
}

// VoiceListen listens on the voice websocket for messages and passes them
// to the voice event handler.
func (s *Session) VoiceListen() (err error) {

	for {
		messageType, message, err := s.VwsConn.ReadMessage()
		if err != nil {
			fmt.Println("Voice Listen Error:", err)
			break
		}

		// Pass received message to voice event handler
		go s.VoiceEvent(messageType, message)
	}

	return
}

// VoiceEvent handles any messages received on the voice websocket
func (s *Session) VoiceEvent(messageType int, message []byte) (err error) {

	if s.Debug {
		fmt.Println("VOICE EVENT:", messageType)
		printJSON(message)
	}

	var e VEvent
	if err := json.Unmarshal(message, &e); err != nil {
		return err
	}

	switch e.Operation {

	case 2: // READY packet
		var st VoiceOP2
		if err := json.Unmarshal(e.RawData, &st); err != nil {
			fmt.Println(e.Type, err)
			printJSON(e.RawData) // TODO: Better error logginEventg
			return err
		}

		// Start the voice websocket heartbeat to keep the connection alive
		go s.VoiceHeartbeat(st.HeartbeatInterval)

		// Store all event data into the session
		s.Vop2 = st

		// We now have enough data to start the UDP connection
		s.VoiceOpenUDP()

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

// VoiceOpenUDP opens a UDP connect to the voice server and completes the
// initial required handshake.  This connect is left open in the session
// and can be used to send or receive audio.
func (s *Session) VoiceOpenUDP() {

	udpHost := fmt.Sprintf("%s:%d", strings.TrimSuffix(s.VEndpoint, ":80"), s.Vop2.Port)
	serverAddr, err := net.ResolveUDPAddr("udp", udpHost)
	if err != nil {
		fmt.Println(err)
	}

	s.UDPConn, err = net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println(err)
	}

	// Create a 70 byte array and put the SSRC code from the Op 2 Voice event
	// into it.  Then send that over the UDP connection to Discord
	sb := make([]byte, 70)
	binary.BigEndian.PutUint32(sb, s.Vop2.SSRC)
	s.UDPConn.Write(sb)

	// Create a 70 byte array and listen for the initial handshake response
	// from Discord.  Once we get it parse the IP and PORT information out
	// of the response.  This should be our public IP and PORT as Discord
	// saw us.
	rb := make([]byte, 70)
	rlen, _, err := s.UDPConn.ReadFromUDP(rb)
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
	jsondata := voiceUDPOp{1, voiceUDPD{"udp", voiceUDPData{ip, port, "plain"}}}

	err = s.VwsConn.WriteJSON(jsondata)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	s.UDPReady = true
}

// VoiceCloseUDP closes the voice UDP connection.
func (s *Session) VoiceCloseUDP() {
	s.UDPConn.Close()
}

type voiceSpeakingData struct {
	Speaking bool `json:"speaking"`
	Delay    int  `json:"delay"`
}

type voiceSpeakingOp struct {
	Op   int               `json:"op"` // Always 5
	Data voiceSpeakingData `json:"d"`
}

func (s *Session) VoiceSpeaking() {

	if s.VwsConn == nil {
		// TODO return an error
		fmt.Println("No Voice websocket.")
		return
	}

	json := voiceSpeakingOp{5, voiceSpeakingData{true, 0}}
	err := s.VwsConn.WriteJSON(json)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
}

// VoiceListenUDP is test code to listen for UDP packets
func (s *Session) VoiceListenUDP() {

	// start the udp keep alive too.  Otherwise listening doesn't get much.
	// THIS DOES NOT WORK YET
	// go s.VoiceUDPKeepalive(s.Vop2.HeartbeatInterval) // lets try the ws timer

	for {
		b := make([]byte, 1024) //TODO DO NOT PUT MAKE INSIDE LOOP
		rlen, _, err := s.UDPConn.ReadFromUDP(b)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			//			return
		}

		if rlen < 1 {
			fmt.Println("Empty UDP packet received")
			continue
			// empty packet?
		}
		fmt.Println("READ FROM UDP: ", b)
	}

}

// VoiceUDPKeepalive sends a packet to keep the UDP connection forwarding
// alive for NATed clients.  Without this no audio can be received
// after short periods of silence.
// Not sure how often this is supposed to be sent or even what payload
// I am suppose to be sending.  So this is very.. unfinished :)
func (s *Session) VoiceUDPKeepalive(i time.Duration) {

	// NONE OF THIS WORKS. SO DON'T USE IT.
	//
	// testing with the above 70 byte SSRC packet.
	//
	// Create a 70 byte array and put the SSRC code from the Op 2 Voice event
	// into it.  Then send that over the UDP connection to Discord

	ticker := time.NewTicker(i * time.Millisecond)
	for range ticker.C {
		sb := make([]byte, 8)
		sb[0] = 0x80
		sb[1] = 0xc9
		sb[2] = 0x00
		sb[3] = 0x01

		ssrcBE := make([]byte, 4)
		binary.BigEndian.PutUint32(ssrcBE, s.Vop2.SSRC)

		sb[4] = ssrcBE[0]
		sb[5] = ssrcBE[1]
		sb[6] = ssrcBE[2]
		sb[7] = ssrcBE[3]

		s.UDPConn.Write(ssrcBE)
	}
}

type voiceHeartbeatOp struct {
	Op   int `json:"op"` // Always 3
	Data int `json:"d"`
}

// VoiceHeartbeat sends regular heartbeats to voice Discord so it knows the client
// is still connected.  If you do not send these heartbeats Discord will
// disconnect the websocket connection after a few seconds.
func (s *Session) VoiceHeartbeat(i time.Duration) {

	ticker := time.NewTicker(i * time.Millisecond)
	for {
		err := s.VwsConn.WriteJSON(voiceHeartbeatOp{3, int(time.Now().Unix())})
		if err != nil {
			s.VoiceReady = false
			fmt.Println(err)
			return // TODO LOG ERROR
		}

		s.VoiceReady = true
		<-ticker.C
	}
}
