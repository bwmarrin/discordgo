package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	//If the message is "ping" reply with a basic embed.
	if m.Content == "ping" {
		embed := &discordgo.MessageEmbed{
			Title:       "Embed Title",
			URL:         "https://github.com/bwmarrin/discordgo",
			Description: "Embed Description",
			Timestamp:   "2021-05-28",
			Color:       0x78141b,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Inline field 1 title",
					Value:  "value 1",
					Inline: true,
				},
				{
					Name:   "Inline field 2 title",
					Value:  "value 2",
					Inline: true,
				},
				{
					Name:   "Regular field title",
					Value:  "value 3",
					Inline: false,
				},
				{
					Name:   "Regular field 2 title",
					Value:  "value 4",
					Inline: false,
				},
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}
