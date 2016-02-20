// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains code related to Discord voice suppport

package discordgo

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ------------------------------------------------------------------------------------------------
// Code related to both Voice Websocket and UDP connections.
// ------------------------------------------------------------------------------------------------

// A Voice struct holds all data and functions related to Discord Voice support.
type Voice struct {
	sync.Mutex              // future use
	Ready      bool         // If true, voice is ready to send/receive audio
	Debug      bool         // If true, print extra logging
	OP2        *voiceOP2    // exported for dgvoice, may change.
	OpusSend   chan []byte  // Chan for sending opus audio
	OpusRecv   chan *Packet // Chan for receiving opus audio
	//	FrameRate  int         // This can be used to set the FrameRate of Opus data
	//	FrameSize  int         // This can be used to set the FrameSize of Opus data

	wsConn  *websocket.Conn
	UDPConn *net.UDPConn // this will become unexported soon.

	sessionID string
	token     string
	endpoint  string
	guildID   string
	channelID string
	userID    string

	// Used to send a close signal to goroutines
	close chan struct{}
}

// ------------------------------------------------------------------------------------------------
// Code related to the Voice websocket connection
// ------------------------------------------------------------------------------------------------

// A voiceOP2 stores the data for the voice operation 2 websocket event
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
func (v *Voice) Open() (err error) {

	v.Lock()
	defer v.Unlock()

	// Don't open a websocket if one is already open
	if v.wsConn != nil {
		return
	}

	// Connect to Voice Websocket
	vg := fmt.Sprintf("wss://%s", strings.TrimSuffix(v.endpoint, ":80"))
	v.wsConn, _, err = websocket.DefaultDialer.Dial(vg, nil)
	if err != nil {
		fmt.Println("VOICE error opening websocket:", err)
		return
	}

	data := voiceHandshakeOp{0, voiceHandshakeData{v.guildID, v.userID, v.sessionID, v.token}}

	err = v.wsConn.WriteJSON(data)
	if err != nil {
		fmt.Println("VOICE error sending init packet:", err)
		return
	}

	// Start a listening for voice websocket events
	// TODO add a check here to make sure Listen worked by monitoring
	// a chan or bool?
	v.close = make(chan struct{})
	go v.wsListen(v.wsConn, v.close)

	return
}

// wsListen listens on the voice websocket for messages and passes them
// to the voice event handler.  This is automatically called by the Open func
func (v *Voice) wsListen(wsConn *websocket.Conn, close <-chan struct{}) {

	for {
		messageType, message, err := v.wsConn.ReadMessage()
		if err != nil {
			// TODO: add reconnect, matching wsapi.go:listen()
			// TODO: Handle this problem better.
			// TODO: needs proper logging
			fmt.Println("Voice Listen Error:", err)
			return
		}

		// Pass received message to voice event handler
		select {
		case <-close:
			return
		default:
			go v.wsEvent(messageType, message)
		}
	}
}

// wsEvent handles any voice websocket events. This is only called by the
// wsListen() function.
func (v *Voice) wsEvent(messageType int, message []byte) {

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
		go v.wsHeartbeat(v.wsConn, v.close, v.OP2.HeartbeatInterval)
		// TODO monitor a chan/bool to verify this was successful

		// Start the UDP connection
		err := v.udpOpen()
		if err != nil {
			fmt.Println("Error opening udp connection: ", err)
			return
		}

		// Start the opusSender.
		// TODO: Should we allow 48000/960 values to be user defined?
		if v.OpusSend == nil {
			v.OpusSend = make(chan []byte, 2)
		}
		go v.opusSender(v.UDPConn, v.close, v.OpusSend, 48000, 960)

		// Start the opusReceiver
		if v.OpusRecv == nil {
			v.OpusRecv = make(chan *Packet, 2)
		}
		go v.opusReceiver(v.UDPConn, v.close, v.OpusRecv)
		return

	case 3: // HEARTBEAT response
		// add code to use this to track latency?
		return

	case 4:
		// TODO

	case 5:
		// SPEAKING TRUE/FALSE NOTIFICATION
		/*
			{
				"user_id": "1238921738912",
				"ssrc": 2,
				"speaking": false
			}
		*/

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

// NOTE :: When a guild voice server changes how do we shut this down
// properly, so a new connection can be setup without fuss?
//
// wsHeartbeat sends regular heartbeats to voice Discord so it knows the client
// is still connected.  If you do not send these heartbeats Discord will
// disconnect the websocket connection after a few seconds.
func (v *Voice) wsHeartbeat(wsConn *websocket.Conn, close <-chan struct{}, i time.Duration) {

	if close == nil || wsConn == nil {
		return
	}

	var err error
	ticker := time.NewTicker(i * time.Millisecond)
	for {
		err = wsConn.WriteJSON(voiceHeartbeatOp{3, int(time.Now().Unix())})
		if err != nil {
			fmt.Println("wsHeartbeat send error: ", err)
			return
		}

		select {
		case <-ticker.C:
			// continue loop and send heartbeat
		case <-close:
			return
		}
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
func (v *Voice) Speaking(b bool) (err error) {

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

// udpOpen opens a UDP connection to the voice server and completes the
// initial required handshake.  This connection is left open in the session
// and can be used to send or receive audio.  This should only be called
// from voice.wsEvent OP2
func (v *Voice) udpOpen() (err error) {

	v.Lock()
	defer v.Unlock()

	if v.wsConn == nil {
		return fmt.Errorf("nil voice websocket")
	}

	if v.UDPConn != nil {
		return fmt.Errorf("udp connection already open")
	}

	if v.close == nil {
		return fmt.Errorf("nil close channel")
	}

	if v.endpoint == "" {
		return fmt.Errorf("empty endpoint")
	}

	host := fmt.Sprintf("%s:%d", strings.TrimSuffix(v.endpoint, ":80"), v.OP2.Port)
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		fmt.Println("udpOpen resolve addr error: ", err)
		// TODO better logging
		return
	}

	v.UDPConn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("udpOpen dial udp error: ", err)
		// TODO better logging
		return
	}

	// Create a 70 byte array and put the SSRC code from the Op 2 Voice event
	// into it.  Then send that over the UDP connection to Discord
	sb := make([]byte, 70)
	binary.BigEndian.PutUint32(sb, v.OP2.SSRC)
	_, err = v.UDPConn.Write(sb)
	if err != nil {
		fmt.Println("udpOpen udp write error : ", err)
		// TODO better logging
		return
	}

	// Create a 70 byte array and listen for the initial handshake response
	// from Discord.  Once we get it parse the IP and PORT information out
	// of the response.  This should be our public IP and PORT as Discord
	// saw us.
	rb := make([]byte, 70)
	rlen, _, err := v.UDPConn.ReadFromUDP(rb)
	if err != nil {
		fmt.Println("udpOpen udp read error : ", err)
		// TODO better logging
		return
	}
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

	// Grab port from position 68 and 69
	port := binary.LittleEndian.Uint16(rb[68:70])

	// Take the data from above and send it back to Discord to finalize
	// the UDP connection handshake.
	data := voiceUDPOp{1, voiceUDPD{"udp", voiceUDPData{ip, port, "plain"}}}

	err = v.wsConn.WriteJSON(data)
	if err != nil {
		fmt.Println("udpOpen write json error:", err)
		return
	}

	// start udpKeepAlive
	go v.udpKeepAlive(v.UDPConn, v.close, 5*time.Second)
	// TODO: find a way to check that it fired off okay

	return
}

// udpKeepAlive sends a udp packet to keep the udp connection open
// This is still a bit of a "proof of concept"
func (v *Voice) udpKeepAlive(UDPConn *net.UDPConn, close <-chan struct{}, i time.Duration) {

	if UDPConn == nil || close == nil {
		return
	}

	var err error
	var sequence uint64

	packet := make([]byte, 8)

	ticker := time.NewTicker(i)
	for {

		binary.LittleEndian.PutUint64(packet, sequence)
		sequence++

		_, err = UDPConn.Write(packet)
		if err != nil {
			fmt.Println("udpKeepAlive udp write error : ", err)
			return
		}

		select {
		case <-ticker.C:
			// continue loop and send keepalive
		case <-close:
			return
		}
	}
}

// opusSender will listen on the given channel and send any
// pre-encoded opus audio to Discord.  Supposedly.
func (v *Voice) opusSender(UDPConn *net.UDPConn, close <-chan struct{}, opus <-chan []byte, rate, size int) {

	if UDPConn == nil || close == nil {
		return
	}

	runtime.LockOSThread()

	// Voice is now ready to receive audio packets
	// TODO: this needs reviewed as I think there must be a better way.
	v.Ready = true
	defer func() { v.Ready = false }()

	var sequence uint16
	var timestamp uint32
	var recvbuf []byte
	var ok bool
	udpHeader := make([]byte, 12)

	// build the parts that don't change in the udpHeader
	udpHeader[0] = 0x80
	udpHeader[1] = 0x78
	binary.BigEndian.PutUint32(udpHeader[8:], v.OP2.SSRC)

	// start a send loop that loops until buf chan is closed
	ticker := time.NewTicker(time.Millisecond * time.Duration(size/(rate/1000)))
	for {

		// Get data from chan.  If chan is closed, return.
		select {
		case <-close:
			return
		case recvbuf, ok = <-opus:
			if !ok {
				return
			}
			// else, continue loop
		}

		// Add sequence and timestamp to udpPacket
		binary.BigEndian.PutUint16(udpHeader[2:], sequence)
		binary.BigEndian.PutUint32(udpHeader[4:], timestamp)

		// Combine the UDP Header and the opus data
		sendbuf := append(udpHeader, recvbuf...)

		// block here until we're exactly at the right time :)
		// Then send rtp audio packet to Discord over UDP
		select {
		case <-close:
			return
		case <-ticker.C:
			// continue
		}
		_, err := UDPConn.Write(sendbuf)

		if err != nil {
			fmt.Println("error writing to udp connection: ", err)
			return
		}

		if (sequence) == 0xFFFF {
			sequence = 0
		} else {
			sequence++
		}

		if (timestamp + uint32(size)) >= 0xFFFFFFFF {
			timestamp = 0
		} else {
			timestamp += uint32(size)
		}
	}
}

// A Packet contains the headers and content of a received voice packet.
type Packet struct {
	SSRC      uint32
	Sequence  uint16
	Timestamp uint32
	Type      []byte
	Opus      []byte
	PCM       []int16
}

// opusReceiver listens on the UDP socket for incoming packets
// and sends them across the given channel
// NOTE :: This function may change names later.
func (v *Voice) opusReceiver(UDPConn *net.UDPConn, close <-chan struct{}, c chan *Packet) {

	if UDPConn == nil || close == nil {
		return
	}

	p := Packet{}
	recvbuf := make([]byte, 1024)

	for {
		rlen, err := UDPConn.Read(recvbuf)
		if err != nil {
			fmt.Println("opusReceiver UDP Read error:", err)
			return
		}

		select {
		case <-close:
			return
		default:
			// continue loop
		}

		// For now, skip anything except audio.
		if rlen < 12 || recvbuf[0] != 0x80 {
			continue
		}

		p.Type = recvbuf[0:2]
		p.Sequence = binary.BigEndian.Uint16(recvbuf[2:4])
		p.Timestamp = binary.BigEndian.Uint32(recvbuf[4:8])
		p.SSRC = binary.BigEndian.Uint32(recvbuf[8:12])
		p.Opus = recvbuf[12:rlen]

		if c != nil {
			c <- &p
		}
	}
}

// Close closes the voice ws and udp connections
func (v *Voice) Close() {

	v.Lock()
	defer v.Unlock()

	v.Ready = false

	if v.close != nil {
		close(v.close)
		v.close = nil
	}

	if v.UDPConn != nil {
		err := v.UDPConn.Close()
		if err != nil {
			fmt.Println("error closing udp connection: ", err)
		}
		v.UDPConn = nil
	}

	if v.wsConn != nil {
		err := v.wsConn.Close()
		if err != nil {
			fmt.Println("error closing websocket connection: ", err)
		}
		v.wsConn = nil
	}
}
