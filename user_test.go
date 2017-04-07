package discordgo

import "testing"

func TestUser(t *testing.T) {
	t.Parallel()

	user := &User{
		Username:      "bob",
		Discriminator: "8192",
	}

	if user.String() != "bob#8192" {
		t.Errorf("user.String() == %v", user.String())
	}
}
