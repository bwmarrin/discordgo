package discordgo

import (
	"testing"
)

func TestContentWithMoreMentionsReplaced(t *testing.T) {
	s := &Session{StateEnabled: true, State: NewState()}

	user := &User{
		ID:       1,
		Username: "User Name",
	}

	s.StateEnabled = true
	s.State.GuildAdd(&Guild{ID: 10})
	s.State.RoleAdd(10, &Role{
		ID:          20,
		Name:        "Role Name",
		Mentionable: true,
	})
	s.State.MemberAdd(&Member{
		User:    user,
		Nick:    "User Nick",
		GuildID: 10,
	})
	s.State.ChannelAdd(&Channel{
		Name:    "Channel Name",
		GuildID: 10,
		ID:      30,
	})

	m := &Message{
		Content:      "<@&20> <@!1> <@1> <#30>",
		ChannelID:    30,
		MentionRoles: []int64{20},
		Mentions:     []*User{user},
	}

	if result, _ := m.ContentWithMoreMentionsReplaced(s); result != "@Role Name @User Nick @User Name #Channel Name" {
		t.Error(result)
	}
}
