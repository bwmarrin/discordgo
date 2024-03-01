// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package discordgo

import (
	"testing"
)

func TestMember_DisplayName(t *testing.T) {
	user := &User{
		GlobalName: "Global",
	}
	t.Run("no server nickname set", func(t *testing.T) {
		m := &Member{
			Nick: "",
			User: user,
		}
		if dn := m.DisplayName(); dn != user.GlobalName {
			t.Errorf("Member.DisplayName() = %v, want %v", dn, user.GlobalName)
		}
	})
	t.Run("server nickname set", func(t *testing.T) {
		m := &Member{
			Nick: "Server",
			User: user,
		}
		if dn := m.DisplayName(); dn != m.Nick {
			t.Errorf("Member.DisplayName() = %v, want %v", dn, m.Nick)
		}
	})
}
