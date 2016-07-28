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

type Timestamp string

func (t Timestamp) Parse() (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000000-07:00", string(t))
}
