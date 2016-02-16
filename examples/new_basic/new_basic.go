// This file provides a basic "quick start" example of using the Discordgo
// package to connect to Discord using the New() helper function.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {

	// Check for Username and Password CLI arguments.
	if len(os.Args) != 3 {
		fmt.Println("You must provide username and password as arguments. See below example.")
		fmt.Println(os.Args[0], " [username] [password]")
		return
	}

	// Call the helper function New() passing username and password command
	// line arguments. This returns a new Discord session, authenticates,
	// connects to the Discord data websocket, and listens for events.
	dg, err := discordgo.New(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

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

	// Print message to stdout.
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
}
