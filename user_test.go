package discordgo

import "testing"

func TestUser(t *testing.T) {
	t.Parallel()

	user := &User{
		ID:            "123456789",
		Username:      "bob",
		Discriminator: "8192",
	}

	if user.String() != "bob#8192" {
		t.Errorf("user.String() == %v", user.String())
	}

	if user.Mention() != "<@123456789>" {
		t.Errorf("user.Mention() == %v", user.Mention())
	}
}

func TestUser_SendMessage(t *testing.T) {
	if envAdmin == "" {
		t.Skip("Skipping, DG_ADMIN not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	user, err := dg.User(envAdmin)
	if err != nil {
		t.Fatalf("failed to retrieve user %s", err)
	}

	_, err = user.SendMessage("Hi boss, just testing!", nil, nil)
	if err != nil {
		t.Fatalf("Error while sending message to %s: %s", user, err)
	}
}
