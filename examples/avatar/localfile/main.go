package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Email       string
	Password    string
	Token       string
	Avatar      string
	BotID       string
	BotUsername string
)

func init() {

	flag.StringVar(&Email, "e", "", "Account Email")
	flag.StringVar(&Password, "p", "", "Account Password")
	flag.StringVar(&Token, "t", "", "Account Token")
	flag.StringVar(&Avatar, "f", "./avatar.jpg", "Avatar File Name")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided login information.
	// Use discordgo.New(Token) to just use a token for login.
	dg, err := discordgo.New(Email, Password, Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	dg.Open()

	bot, err := dg.User("@me")
	if err != nil {
		fmt.Println("error fetching the bot details,", err)
		return
	}

	BotID = bot.ID
	BotUsername = bot.Username
	changeAvatar(dg)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

// Helper function to change the avatar
func changeAvatar(s *discordgo.Session) {
	img, err := ioutil.ReadFile(Avatar)
	if err != nil {
		fmt.Println(err)
	}

	base64 := base64.StdEncoding.EncodeToString(img)

	avatar := fmt.Sprintf("data:%s;base64,%s", http.DetectContentType(img), string(base64))

	_, err = s.UserUpdate("", "", BotUsername, avatar, "")
	if err != nil {
		fmt.Println(err)
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Print message to stdout.
	fmt.Printf("%-20s %-20s\n %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
}
