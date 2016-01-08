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
