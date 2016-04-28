// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains code related to discordgo package logging

package discordgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
)

const (

	// Critical Errors that could lead to data loss or panic
	// Only errors that would not be returned to a calling function
	LogError int = iota

	// Very abnormal events.
	// Errors that are also returend to a calling function.
	LogWarning

	// Normal non-error activity
	// Generally, not overly spammy events
	LogInformational

	// Detailed activity
	// All HTTP/Websocket packets.
	// Very spammy and will impact performance
	LogDebug
)

// msglog provides package wide logging consistancy for discordgo
// the format, a...  portion this command follows that of fmt.Printf
//   msgL   : LogLevel of the message
//   caller : 1 + the number of callers away from the message source
//   format : Printf style message format
//   a ...  : comma seperated list of values to pass
func msglog(msgL, caller int, format string, a ...interface{}) {

	pc, file, line, _ := runtime.Caller(caller)

	files := strings.Split(file, "/")
	file = files[len(files)-1]

	name := runtime.FuncForPC(pc).Name()
	fns := strings.Split(name, ".")
	name = fns[len(fns)-1]

	msg := fmt.Sprintf(format, a...)

	log.Printf("[DG%d] %s:%d %s %s\n", msgL, file, line, name, msg)
}

// helper function that wraps msglog for the Session struct
// This adds a check to insure the message is only logged
// if the session log level is equal or higher than the
// message log level
func (s *Session) log(msgL int, format string, a ...interface{}) {

	if s.Debug { // Deprecated
		s.LogLevel = LogDebug
	}

	if msgL > s.LogLevel {
		return
	}

	msglog(msgL, 2, format, a...)
}

// helper function that wraps msglog for the VoiceConnection struct
// This adds a check to insure the message is only logged
// if the voice connection log level is equal or higher than the
// message log level
func (v *VoiceConnection) log(msgL int, format string, a ...interface{}) {

	if v.Debug { // Deprecated
		v.LogLevel = LogDebug
	}

	if msgL > v.LogLevel {
		return
	}

	msglog(msgL, 2, format, a...)
}

// printEvent prints out a WSAPI event.
func printEvent(e *Event) {
	log.Println(fmt.Sprintf("Event. Type: %s, State: %d Operation: %d Direction: %d", e.Type, e.State, e.Operation, e.Direction))
	printJSON(e.RawData)
}

// printJSON is a helper function to display JSON data in a easy to read format.
func printJSON(body []byte) {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		log.Print("JSON parse error: ", error)
	}
	log.Println(string(prettyJSON.Bytes()))
}
