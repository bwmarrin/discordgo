// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains high level helper functions and easy entry points for the
// entire discordgo package.  These functions are beling developed and are very
// experimental at this point.  They will most likley change so please use the
// low level functions if that's a problem.

// package discordgo provides Discord binding for Go
package discordgo

import "fmt"

// DiscordGo Version, follows Symantic Versioning. (http://semver.org/)
const VERSION = "0.8.0-alpha"

/*
type Config struct {
	Debug bool
}
*/
/*
// possible future main struct for discord connection
type Discord struct {
	Debug    bool      // Set to true to enable debug logging
	Token    string    // authentication token
	User     User      // authenticated user info
	Guilds   []Guild   // Cached Guild info
	Channels []Channel // Cached Channel info
	API      api       // all api endpoint functions
	DataWS   dataWS    // data websocket connection
	VoiceWS  voiceWS   // voice websocket/udp connections
	VoiceUDP voiceUDP
}

type api struct {
	Session
}
type dataWS struct {
}

type voiceWS struct {
}

type voiceUDP struct {
}
*/
// NOTICE: This function should be considered unstable because I am still
// exploring the best way to implement my goals here.  So, it is more likely
// to change than than the low level API functions.
//
// New creates a new Discord session interface and will automate some startup
// tasks if given enough information to do so.  Currently you can pass zero
// arguments and it will return an empty Discord session. If you pass a token
// or username and password (in that order), then it will attempt to login to
// Discord and open a websocket connection.
func New(args ...interface{}) (s *Session, err error) {

	// Create an empty Session interface.
	s = &Session{
		State:        NewState(),
		StateEnabled: true,
	}

	// If no arguments are passed return the empty Session interface.
	// Later I will add default values, if appropriate.
	if args == nil {
		return
	}

	// Varibles used below when parsing func arguments
	var auth, pass string

	// Parse passed arguments
	for _, arg := range args {

		switch v := arg.(type) {

		case []string:
			if len(v) > 2 {
				err = fmt.Errorf("Too many string parameters provided.")
				return
			}

			// First string is either token or username
			if len(v) > 0 {
				auth = v[0]
			}

			// If second string exists, it must be a password.
			if len(v) > 1 {
				pass = v[1]
			}

		case string:
			// First string must be either auth token or username.
			// Second string must be a password.
			// Only 2 input strings are supported.

			if auth == "" {
				auth = v
			} else if pass == "" {
				pass = v
			} else {
				err = fmt.Errorf("Too many string parameters provided.")
				return
			}

			//		case Config:
			// TODO: Parse configuration

		default:
			err = fmt.Errorf("Unsupported parameter type provided.")
			return
		}
	}

	// If only one string was provided, assume it is an auth token.
	// Otherwise get auth token from Discord
	if pass == "" {
		s.Token = auth
	} else {
		s.Token, err = s.Login(auth, pass)
		if err != nil || s.Token == "" {
			err = fmt.Errorf("Unable to fetch discord authentication token. %v", err)
			return
		}
	}

	// TODO: Add code here to fetch authenticated user info like settings,
	// avatar, User ID, etc.  If fails, return error.

	// Open websocket connection
	err = s.Open()
	if err != nil {
		fmt.Println(err)
	}

	// Do websocket handshake.
	err = s.Handshake()
	if err != nil {
		fmt.Println(err)
	}

	// Listen for events.
	go s.Listen()

	return
}

// Close closes a Discord session
// TODO: Add support for Voice WS/UDP connections
func (s *Session) Close() {

	s.DataReady = false

	close(s.listenChan)
	close(s.heartbeatChan)

	if s.wsConn != nil {
		s.wsConn.Close()
	}
}
