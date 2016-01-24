// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains utility functions for the discordgo package. These
// functions are not exported and are likely to change substantially in
// the future to match specific needs of the discordgo package itself.

package discordgo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// printEvent prints out a WSAPI event.
func printEvent(e *Event) {
	fmt.Println(fmt.Sprintf("Event. Type: %s, State: %d Operation: %d Direction: %d", e.Type, e.State, e.Operation, e.Direction))
	printJSON(e.RawData)
}

// printJSON is a helper function to display JSON data in a easy to read format.
func printJSON(body []byte) {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		fmt.Print("JSON parse error: ", error)
	}
	fmt.Println(string(prettyJSON.Bytes()))
}
