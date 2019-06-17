package discordgo

import "testing"

func TestMember_GetTopRole(t *testing.T) {
	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
	}

	if envAdmin == "" {
		t.Skip("Skipping, DG_ADMIN not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	g, err := dg.State.Guild(envGuild)
	if err != nil {
		t.Fatalf("Guild not found, id: %s; %s", envGuild, err)
	}

	if g.Unavailable {
		t.Fatalf("Guild %s is still unavailable", envGuild)
	}

	m, err := g.GetMember(envAdmin)
	if err != nil {
		t.Fatalf("User %s is not in Guild", envAdmin)
	}

	c, err := m.GetTopRole()
	if err != nil {
		t.Fatalf("Failed at getting the top role; %s", err)
	}
	if c.ID == envGuild {
		t.Fatalf("Got lowest role instead of top role")
	}
}
