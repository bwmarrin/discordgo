// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains high level helper functions and easy entry points for the
// entire discordgo package.  These functions are beling developed and are very
// experimental at this point.  They will most likley change so please use the
// low level functions if that's a problem.

// Package discordgo provides Discord binding for Go
package discordgo

import (
	"fmt"
	"reflect"
)

// VERSION of Discordgo, follows Symantic Versioning. (http://semver.org/)
const VERSION = "0.11.0-alpha"

// New creates a new Discord session and will automate some startup
// tasks if given enough information to do so.  Currently you can pass zero
// arguments and it will return an empty Discord session.
// There are 3 ways to call New:
//     With a single auth token - All requests will use the token blindly,
//         no verification of the token will be done and requests may fail.
//     With an email and password - Discord will sign in with the provided
//         credentials.
//     With an email, password and auth token - Discord will verify the auth
//         token, if it is invalid it will sign in with the provided
//         credentials. This is the Discord recommended way to sign in.
func New(args ...interface{}) (s *Session, err error) {

	// Create an empty Session interface.
	s = &Session{
		State:                  NewState(),
		StateEnabled:           true,
		Compress:               true,
		ShouldReconnectOnError: true,
	}

	// If no arguments are passed return the empty Session interface.
	// Later I will add default values, if appropriate.
	if args == nil {
		return
	}

	// Variables used below when parsing func arguments
	var auth, pass string

	// Parse passed arguments
	for _, arg := range args {

		switch v := arg.(type) {

		case []string:
			if len(v) > 3 {
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

			// If third string exists, it must be an auth token.
			if len(v) > 2 {
				s.Token = v[2]
			}

		case string:
			// First string must be either auth token or username.
			// Second string must be a password.
			// Only 2 input strings are supported.

			if auth == "" {
				auth = v
			} else if pass == "" {
				pass = v
			} else if s.Token == "" {
				s.Token = v
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
	// Otherwise get auth token from Discord, if a token was specified
	// Discord will verify it for free, or log the user in if it is
	// invalid.
	if pass == "" {
		s.Token = auth
	} else {
		err = s.Login(auth, pass)
		if err != nil || s.Token == "" {
			err = fmt.Errorf("Unable to fetch discord authentication token. %v", err)
			return
		}
	}

	// The Session is now able to have RestAPI methods called on it.
	// It is recommended that you now call Open() so that events will trigger.

	return
}

func (s *Session) AddHandler(handler interface{}) {
	s.Lock()
	defer s.Unlock()

	handlerType := reflect.TypeOf(handler)

	if handlerType.NumIn() != 2 {
		panic("Unable to add event handler, handler must be of the type func(*discordgo.Session, *discordgo.EventType).")
	}

	if handlerType.In(0) != reflect.TypeOf(s) {
		panic("Unable to add event handler, first argument must be of type *discordgo.Session.")
	}

	if s.Handlers == nil {
		s.Unlock()
		s.initialize()
		s.Lock()
	}

	eventType := handlerType.In(1)

	handlers := s.Handlers[eventType]
	if handlers == nil {
		handlers = []interface{}{}
	}

	handlers = append(handlers, handler)
	s.Handlers[eventType] = handlers
}

func (s *Session) Handle(event interface{}) (handled bool) {
	s.RLock()
	defer s.RUnlock()

	eventType := reflect.TypeOf(event)

	handlers, ok := s.Handlers[eventType]
	if !ok {
		return
	}

	for _, handler := range handlers {
		reflect.ValueOf(handler).Call([]reflect.Value{reflect.ValueOf(s), reflect.ValueOf(event)})
		handled = true
	}

	return
}

// initialize adds all internal handlers and state tracking handlers.
func (s *Session) initialize() {
	s.Lock()
	defer s.Unlock()

	s.Handlers = map[interface{}][]interface{}{}
	s.AddHandler(s.onReady)
	s.AddHandler(s.onVoiceServerUpdate)
	s.AddHandler(s.onVoiceStateUpdate)

	s.AddHandler(s.State.onReady)
	s.AddHandler(s.State.onMessageCreate)
	s.AddHandler(s.State.onMessageUpdate)
	s.AddHandler(s.State.onMessageDelete)
}

// onReady handles the ready event.
func (s *Session) onReady(se *Session, r *Ready) {
	go s.heartbeat(s.wsConn, s.listening, r.HeartbeatInterval)
}
