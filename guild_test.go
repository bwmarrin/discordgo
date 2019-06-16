package discordgo

import "testing"

func TestGuild_GetChannel(t *testing.T) {
	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
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

	_, err = g.GetChannel(envChannel)
	if err != nil {
		t.Fatalf("Channel not found in guild")
	}
}

func TestGuild_GetRole(t *testing.T) {
	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
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

	_, err = g.GetRole(envRole)
	if err != nil {
		t.Fatalf("Role not found in guild")
	}
}
