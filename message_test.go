package discordgo

import (
	"testing"
)

func TestContentWithMoreMentionsReplaced(t *testing.T) {
	s := &Session{StateEnabled: true, State: NewState()}

	user := &User{
		ID:       "user",
		Username: "User Name",
	}

	s.State.GuildAdd(&Guild{ID: "guild"})
	s.State.RoleAdd("guild", &Role{
		ID:          "role",
		Name:        "Role Name",
		Mentionable: true,
	})
	s.State.MemberAdd(&Member{
		User:    user,
		Nick:    "User Nick",
		GuildID: "guild",
	})
	s.State.ChannelAdd(&Channel{
		Name:    "Channel Name",
		GuildID: "guild",
		ID:      "channel",
	})
	m := &Message{
		Content:      "<@&role> <@!user> <@user> <#channel>",
		ChannelID:    "channel",
		MentionRoles: []string{"role"},
		Mentions:     []*User{user},
	}
	if result, _ := m.ContentWithMoreMentionsReplaced(s); result != "@Role Name @User Nick @User Name #Channel Name" {
		t.Error(result)
	}
}
func TestGettingEmojisFromMessage(t *testing.T) {
	msg := "test test <:kitty14:811736565172011058> <:kitty4:811736468812595260>"
	m := &Message{
		Content: msg,
	}
	emojis := m.GetCustomEmojis()
	if len(emojis) < 1 {
		t.Error("No emojis found.")
		return
	}

}

func TestMessage_DisplayName(t *testing.T) {
	user := &User{
		GlobalName: "Global",
	}
	t.Run("no server nickname set", func(t *testing.T) {
		m := &Message{
			Member: &Member{
				Nick: "",
			},
			Author: user,
		}
		if dn := m.DisplayName(); dn != user.GlobalName {
			t.Errorf("Member.DisplayName() = %v, want %v", dn, user.GlobalName)
		}
	})

	t.Run("server nickname set", func(t *testing.T) {
		m := &Message{
			Member: &Member{
				Nick: "Server",
			},
			Author: user,
		}
		if dn := m.DisplayName(); dn != m.Member.Nick {
			t.Errorf("Member.DisplayName() = %v, want %v", dn, m.Member.Nick)
		}
	})

	bot := &User{
		Username: "Bot",
		Bot:      true,
	}

	t.Run("bot no server nickname set", func(t *testing.T) {
		m := &Message{
			Member: &Member{
				Nick: "",
			},
			Author: bot,
		}
		if dn := m.DisplayName(); dn != m.Author.Username {
			t.Errorf("Member.DisplayName() = %v, want %v", dn, m.Author.Username)
		}
	})

	t.Run("bot server nickname set", func(t *testing.T) {
		m := &Message{
			Member: &Member{
				Nick: "Server",
			},
			Author: bot,
		}
		if dn := m.DisplayName(); dn != m.Member.Nick {
			t.Errorf("Member.DisplayName() = %v, want %v", dn, m.Member.Nick)
		}
	})
}
