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

	s.State.GuildAdd(&Guild{ID: "guild"}, s)
	s.State.RoleAdd("guild", &Role{
		ID:          "role",
		Name:        "Role Name",
		Mentionable: true,
	})
	s.State.MemberAdd(&Member{
		User:    user,
		Nick:    "User Nick",
		GuildID: "guild",
	}, s)
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
		Session:      s,
	}
	if result, _ := m.ContentWithMoreMentionsReplaced(); result != "@Role Name @User Nick @User Name #Channel Name" {
		t.Error(result)
	}
}

func TestMessage_Edit(t *testing.T) {
	if envChannel == "" {
		t.Skip("Skipping, DG_CHANNEL not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	c, err := dg.State.Channel(envChannel)
	if err != nil {
		t.Fatalf("Channel %s wasn't cached", envChannel)
	}

	m, err := c.SendMessage("Testing... please wait for edit", nil, nil)
	if err != nil {
		t.Fatalf("Error while sending message: %s", err)
	}

	e := m.NewMessageEdit().SetContent("Testing editing message")
	m, err = m.Edit(e)
	if err != nil {
		t.Fatalf("Error while editing message: %s", err)
	}
}

func TestEmbed(t *testing.T) {
	if envChannel == "" {
		t.Skip("Skipping, DG_CHANNEL not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	c, err := dg.State.Channel(envChannel)
	if err != nil {
		t.Fatalf("Channel %s wasn't cached", envChannel)
	}

	_, err = c.SendMessage(
		"",
		NewEmbed().
			SetDescription("testing the ability to make and send embeds").
			SetAuthorName("Library Tester").
			SetTitle("Embed Test").
			SetFooterText("t-t-testiiing").
			SetImage("https://imgur.com/KAHPV0d.gif"),
		nil,
	)
}
