// This file provides a basic "quick start" example of using the Discordgo
// package to connect to Discord using the low level API functions.
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Email    string
	Password string
	Token    string
	BotID    string
)

func init() {

	flag.StringVar(&Email, "e", "", "Account Email")
	flag.StringVar(&Password, "p", "", "Account Password")
	flag.StringVar(&Token, "t", "", "Account Token")
	flag.Parse()
}

func main() {

	// Create a new Discord Session struct and set a handler for the
	dg := discordgo.Session{}

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// If no Authentication Token was provided login using the
	// provided Email and Password.
	if Token == "" {
		err := dg.Login(Email, Password)
		if err != nil {
			fmt.Println("error logging into Discord,", err)
			return
		}
	} else {
		dg.Token = Token
	}

	// Open websocket connection to Discord
	err := dg.Open()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Print message to stdout.
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
}
