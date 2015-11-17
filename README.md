# Discordgo

Discordgo provides a mostly complete low-level Golang interface to the Discord
REST and Websocket API.

[![GoDoc](https://godoc.org/github.com/bwmarrin/discordgo?status.svg)](https://godoc.org/github.com/bwmarrin/discordgo) [![Go report](http://goreportcard.com/badge/bwmarrin/discordgo)](http://goreportcard.com/report/bwmarrin/discordgo) [![Build Status](https://travis-ci.org/bwmarrin/discordgo.svg?branch=master)](https://travis-ci.org/bwmarrin/discordgo)

# Usage Example
```go
package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {

	var err error

	// Create a new Discord Session and set a handler for the OnMessageCreate
    // event that happens for every new message on any channel
	Session := discordgo.Session{
		OnMessageCreate: messageCreate,
	}

	// Login to the Discord server and store the authentication token
	// inside the Session
	Session.Token, err = Session.Login("coolusername", "cleverpassword")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Open websocket connection
	err = Session.Open()
	if err != nil {
		fmt.Println(err)
	}

	// Do websocket handshake.
	err = Session.Handshake()
	if err != nil {
		fmt.Println(err)
	}

	// Listen for events.
	Session.Listen()
	return
}

func messageCreate(s *discordgo.Session, m discordgo.Message) {
	fmt.Printf("%25d %s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
}
```

# Documentation

**NOTICE** : This library and the Discord API are unfinished.
Because of that there may be major changes to library functions, constants,
and structures.

- [GoDoc](https://godoc.org/github.com/bwmarrin/discordgo)
- Hand crafted documentation coming soon.

# What Works

Low level functions exist for the majority of the REST and Websocket API.

* Login/Logout
* Open/Close Websocket and listen for events.
* Accept/Create/Delete Invites
* Get User details (Name, ID, Settings, etc)
* List/Create User Channels (Private Message Channels)
* List/Create Guilds
* List/Create Guild Channels
* List Guild Members
* Receive/Send Messages to Channels

# What's Unfinished

* Make changes as needed to pass GoLint, GoVet, GoCyclo, etc. (goreportcard.com)
* Editing User Profile settings
* Permissions related functions.
* Functions for Maintenance Status
* Voice Channel support.

# Credits

Special thanks goes to both the below projects who helped me get started with
this project.  If you're looking for alternative Golang interfaces to Discord
please check both of these out.

* https://github.com/gdraynz/go-discord
* https://github.com/Xackery/discord


# Other Discord APIs

- [go-discord](https://github.com/gdraynz/go-discord)
- [discord](https://github.com/Xackery/discord)
- [discord.py](https://github.com/Rapptz/discord.py)
- [discord.js](https://github.com/discord-js/discord.js)
- [discord.io](https://github.com/izy521/discord.io)
- [Discord.NET](https://github.com/RogueException/Discord.Net)
- [DiscordSharp](https://github.com/Luigifan/DiscordSharp)
- [Discord4J](https://github.com/knobody/Discord4J)
- [discordrb](https://github.com/meew0/discordrb)
