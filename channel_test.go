package discordgo

import "testing"

func TestChannel_SendMessage(t *testing.T) {
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

	_, err = c.SendMessage("Testing Channel.SendMessage", nil, nil)
	if err != nil {
		t.Fatalf("Error while sending message: %s", err)
	}
}
