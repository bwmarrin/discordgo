// This program provides a simple Ping/Pong bot example
// using the DiscordGo API package.
package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

// Will use this to store the Bot's account ID
var BotID string

func main() {

	// Check for Token
	if len(os.Args) != 2 {
		fmt.Println("You must provide the authentication token for your bot account.")
		fmt.Println(os.Args[0], " [token]")
		return
	}

	dg, err := discordgo.New(os.Args[1])
	if err != nil {
		fmt.Println("error creating Discord session: ", err)
		return
	}

	// Get the Bot account information.
	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining bot account details: ", err)
	}

	// Store the account ID for later use.
	BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	dg.Open()

	// Simple way to keep program running until any key press.
	var input string
	fmt.Scanln(&input)
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated user has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	// If the message is "ping"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
