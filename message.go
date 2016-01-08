// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains code related to the Message struct

package discordgo

import (
	"fmt"
	"strings"
)

// A Message stores all data related to a specific Discord message.
type Message struct {
	ID              string        `json:"id"`
	Author          *User         `json:"author"`
	Content         string        `json:"content"`
	Attachments     []*Attachment `json:"attachments"`
	Tts             bool          `json:"tts"`
	Embeds          []*Embed      `json:"embeds"`
	Timestamp       string        `json:"timestamp"`
	MentionEveryone bool          `json:"mention_everyone"`
	EditedTimestamp string        `json:"edited_timestamp"`
	Mentions        []*User       `json:"mentions"`
	ChannelID       string        `json:"channel_id"`
}

// An Attachment stores data for message attachments.
type Attachment struct { //TODO figure this out
}

// An Embed stores data for message embeds.
type Embed struct { // TODO figure this out
}

// ContentWithMentionsReplaced will replace all @<id> mentions with the
// username of the mention.
func (m *Message) ContentWithMentionsReplaced() string {
	content := m.Content
	for _, user := range m.Mentions {
		content = strings.Replace(content, fmt.Sprintf("<@%s>", user.ID),
			fmt.Sprintf("@%s", user.Username), -1)
	}
	return content
}
