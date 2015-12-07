# Discordgo

This package provides low level bindings for the [Discord](https://discordapp.com/) 
REST & Websocket API in the [Go](https://golang.org/)  Programming Language (Golang).

* See out [dgVoice](https://github.com/bwmarrin/dgvoice) for **experimental** 
Discord voice support.

* See out [dgTest](https://github.com/bwmarrin/dgTest) for more examples and test code.

----

[![GoDoc](https://godoc.org/github.com/bwmarrin/discordgo?status.svg)](https://godoc.org/github.com/bwmarrin/discordgo) 
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/bwmarrin/discordgo) 
[![Go report](http://goreportcard.com/badge/bwmarrin/discordgo)](http://goreportcard.com/report/bwmarrin/discordgo) 
[![Build Status](https://travis-ci.org/bwmarrin/discordgo.svg?branch=master)](https://travis-ci.org/bwmarrin/discordgo)

# Usage Examples
See the examples sub-folder for examples.  Each example accepts a username and 
password as a CLI argument when run.

# Documentation

**NOTICE** : This library and the Discord API are unfinished.
Because of that there may be major changes to library functions, constants,
and structures.

- [![GoDoc](https://godoc.org/github.com/bwmarrin/discordgo?status.svg)](https://godoc.org/github.com/bwmarrin/discordgo) 
- [![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/bwmarrin/discordgo) 
- Hand crafted documentation coming eventually.

# What Works

Current package provides a **low level direct mapping** to the majority of Discord 
REST and Websock API.

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
* Finish Voice support.
* Add a higher level interface with user friendly helper functions.

# Other Discord APIs

**Go**:
[gdraynz/**go-discord**](https://github.com/gdraynz/go-discord),
[Xackery/**discord**](https://github.com/Xackery/discord),
[Nerketur/**discordapi**](https://github.com/Nerketur/discordapi)

**.NET**:
[RogueException/**Discord.Net**](https://github.com/RogueException/Discord.Net),
[Luigifan/**DiscordSharp**](https://github.com/Luigifan/DiscordSharp)

**Java**:
[nerd/**Discord4J**](https://github.com/nerd/Discord4J)

**Node.js**:
[izy521/**discord.io**](https://github.com/izy521/discord.io),
[hydrabolt/**discord.js**](https://github.com/hydrabolt/discord.js),
[qeled/**discordie**](https://github.com/qeled/discordie),

**PHP**:
[Cleanse/**discord-hypertext**](https://github.com/Cleanse/discord-hypertext),
[teamreflex/**DiscordPHP**](https://github.com/teamreflex/DiscordPHP)

**Python**:
[Rapptz/**discord.py**](https://github.com/Rapptz/discord.py)

**Ruby**:
[meew0/**discordrb**](https://github.com/meew0/discordrb)

**Scala**:
[eaceaser/**discord-akka**](https://github.com/eaceaser/discord-akka)

**Rust**:
[SpaceManiac/**discord-rs**](https://github.com/SpaceManiac/discord-rs)
