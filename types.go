// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains custom types, currently only a timestamp wrapper.

package discordgo

import (
	"time"
)

// Timestamp stores a timestamp, as sent by the Discord API.
type Timestamp string

// Parse parses a timestamp string into a time.Time object.
// The only time this can fail is if Discord changes their timestamp format.
func (t Timestamp) Parse() (time.Time, error) {
	return time.Parse(time.RFC3339, string(t))
}
