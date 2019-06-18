package discordgo

import "testing"

func prepareMemberTests(t *testing.T) (g *Guild, m *Member) {
	var err error
	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
	}

	if envAdmin == "" {
		t.Skip("Skipping, DG_ADMIN not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	g, err = dg.State.Guild(envGuild)
	if err != nil {
		t.Fatalf("Guild not found, id: %s; %s", envGuild, err)
	}

	if g.Unavailable {
		t.Fatalf("Guild %s is still unavailable", envGuild)
	}

	m, err = g.GetMember(envAdmin)
	if err != nil {
		t.Fatalf("User %s is not in Guild", envAdmin)
	}

	return
}

func TestMember_GetTopRole(t *testing.T) {
	_, m := prepareMemberTests(t)

	c, err := m.GetTopRole()
	if err != nil {
		t.Fatalf("Failed at getting the top role; %s", err)
	}
	if c.ID == envGuild {
		t.Fatalf("Got lowest role instead of top role")
	}
}

/*
// hehe, oops, can't edit the nickname of someone higher in roles than you ofc
func TestMember_EditNickname(t *testing.T) {
	_, m := prepareMemberTests(t)

	err := m.EditNickname("is potentially a wolf 👀")
	if err != nil {
		t.Fatalf("Changing nickname failed: %s", err)
	}
}
*/

func TestMember_AddRole(t *testing.T) {
	g, m := prepareMemberTests(t)

	r, err := g.GetRole(envRole)
	if err != nil {
		t.Skip("Skipping, DG_ROLE has not been set or is not in DG_GUILD")
	}

	err = m.AddRole(r)
	if err != nil {
		t.Fatalf("Adding role to member failed, because: %s", err)
	}
}

func TestMember_RemoveRole(t *testing.T) {
	g, m := prepareMemberTests(t)

	r, err := g.GetRole(envRole)
	if err != nil {
		t.Skip("Skipping, DG_ROLE has not been set or is not in DG_GUILD")
	}

	err = m.RemoveRole(r)
	if err != nil {
		t.Fatalf("Removing role from member failed, because: %s", err)
	}
}
