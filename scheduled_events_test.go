package discordgo

import (
	"testing"
	"time"
)

func TestScheduledsEvents(t *testing.T) {
	if dgBot == nil {
		t.Skip("Skipping, dgBot not set.")
	}

	event, err := dgBot.GuildScheduledEventCreate(envGuild, &GuildScheduledEvent{
		Name:               "Test Event",
		PrivacyLevel:       GuildScheduledEventPrivacyLevelGuildOnly,
		ScheduledStartTime: time.Now().Add(1 * time.Hour),
		ScheduledEndTime:   time.Now().Add(2 * time.Hour),
		Description:        "Awesome Test Event created on livestream",
		EntityType:         GuildScheduledEventEntityTypeExternal,
		EntityMetadata: GuildScheduledEventEntityMetadata{
			Location: "https://discord.com",
		},
	})
	if err != nil || event.Name != "Test Event" {
		t.Fatal(err)
	}

	events, err := dgBot.GuildScheduledEvents(envGuild, true)
	if err != nil {
		t.Fatal(err)
	}

	var foundEvent *GuildScheduledEvent
	for _, e := range events {
		if e.ID == event.ID {
			foundEvent = e
			break
		}
	}
	if foundEvent.Name != event.Name {
		t.Fatal("err on GuildScheduledEvents endpoint. Missing Scheduled Event")
	}

	event.Name = "Test Event Updated"
	eventUpdated, err := dgBot.GuildScheduledEventUpdate(envGuild, event.ID, event)
	if err != nil {
		t.Fatal(err)
	}

	if eventUpdated.Name != event.Name {
		t.Fatal("err on GuildScheduledEventUpdate endpoint. Scheduled Event Name mismatch")
	}

	users, err := dgBot.GuildScheduledEventUsers(envGuild, event.ID, 0, true, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(users) < 1 {
		t.Fatal("err on GuildScheduledEventUsers. No Data")
	}

	err = dgBot.GuildScheduledEventDelete(envGuild, event.ID)
	if err != nil {
		t.Fatal(err)
	}
}
