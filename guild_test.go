package discordgo

import "testing"

func getGuild(t *testing.T) (g *Guild) {
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
	return
}

func TestGuild_GetChannel(t *testing.T) {
	g := getGuild(t)

	_, err := g.GetChannel(envChannel)
	if err != nil {
		t.Fatalf("Channel not found in guild")
	}
}

func TestGuild_GetRole(t *testing.T) {
	g := getGuild(t)

	_, err := g.GetRole(envRole)
	if err != nil {
		t.Fatalf("Role not found in guild")
	}
}

func TestGuild_CreateDeleteRole(t *testing.T) {
	g := getGuild(t)

	r, err := g.CreateRole()
	if err != nil {
		t.Fatalf("Role failed to create in Guild; %s", err)
	}

	editData := &RoleEdit{
		Name:        "OwO a testing role",
		Hoist:       false,
		Color:       0xff00ff,
		Permissions: r.Permissions,
		Mentionable: true,
	}

	r, err = r.EditComplex(editData)
	if err != nil {
		t.Fatalf("Failed at editing role; %s", err)
	}

	err = r.Delete()
	if err != nil {
		t.Fatalf("Failed at deleteing role; %s", err)
	}
}
