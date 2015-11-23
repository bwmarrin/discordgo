// Discordgo - Go bindings for Discord
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
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

// printJSON is a helper function to display JSON data in a easy to read format.
func printJSON(body []byte) {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		fmt.Print("JSON parse error: ", error)
	}
	fmt.Println(string(prettyJSON.Bytes()))
}
