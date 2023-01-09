package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/LightningDev1/discordgo"
)

// Flags
var (
	GuildID        = flag.String("guild", "", "Test guild ID")
	StageChannelID = flag.String("stage", "", "Test stage channel ID")
	BotToken       = flag.String("token", "", "Bot token")
)

func init() { flag.Parse() }

// To be correctly used, the bot needs to be in a guild.
// All actions must be done on a stage channel event
func main() {
	s, _ := discordgo.New("Bot " + *BotToken)
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot is ready")
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	// Create a new Stage instance on the previous channel
	si, err := s.StageInstanceCreate(&discordgo.StageInstanceParams{
		ChannelID:             *StageChannelID,
		Topic:                 "Amazing topic",
		PrivacyLevel:          discordgo.StageInstancePrivacyLevelGuildOnly,
		SendStartNotification: true,
	})
	if err != nil {
		log.Fatalf("Cannot create stage instance: %v", err)
	}
	log.Printf("Stage Instance %s has been successfully created", si.Topic)

	// Edit the stage instance with a new Topic
	si, err = s.StageInstanceEdit(*StageChannelID, &discordgo.StageInstanceParams{
		Topic: "New amazing topic",
	})
	if err != nil {
		log.Fatalf("Cannot edit stage instance: %v", err)
	}
	log.Printf("Stage Instance %s has been successfully edited", si.Topic)

	time.Sleep(5 * time.Second)
	if err = s.StageInstanceDelete(*StageChannelID); err != nil {
		log.Fatalf("Cannot delete stage instance: %v", err)
	}
	log.Printf("Stage Instance %s has been successfully deleted", si.Topic)
}
