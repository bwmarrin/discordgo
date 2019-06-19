package discordgo

import (
	"testing"
)

func TestContentWithMoreMentionsReplaced(t *testing.T) {
	s := &Session{StateEnabled: true, state: NewState()}

	user := &User{
		ID:       "user",
		Username: "User Name",
	}

	s.state.GuildAdd(&Guild{ID: "guild"})
	s.state.RoleAdd("guild", &Role{
		ID:          "role",
		Name:        "Role Name",
		Mentionable: true,
	})
	s.state.MemberAdd(&Member{
		User:    user,
		Nick:    "User Nick",
		GuildID: "guild",
	})
	s.state.ChannelAdd(&Channel{
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
