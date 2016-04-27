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
	LogError int = iota
	LogWarning
	LogNotice
	LogDebug
)

// logs messages to stderr
func msglog(cfgL, msgL int, format string, a ...interface{}) {

	if msgL > cfgL {
		return
	}

	pc, file, line, _ := runtime.Caller(1)

	files := strings.Split(file, "/")
	file = files[len(files)-1]

	name := runtime.FuncForPC(pc).Name()
	fns := strings.Split(name, ".")
	name = fns[len(fns)-1]

	msg := fmt.Sprintf(format, a...)

	log.Printf("%s:%d:%s %s\n", file, line, name, msg)
}

// helper function that wraps msglog for the Session struct
func (s *Session) log(msgL int, format string, a ...interface{}) {

	if s.Debug { // Deprecated
		s.LogLevel = LogDebug
	}
	msglog(s.LogLevel, msgL, format, a...)
}

// helper function that wraps msglog for the VoiceConnection struct
func (v *VoiceConnection) log(msgL int, format string, a ...interface{}) {

	if v.Debug { // Deprecated
		v.LogLevel = LogDebug
	}

	msglog(v.LogLevel, msgL, format, a...)
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
