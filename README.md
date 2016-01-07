Discordgo [![GoDoc](https://godoc.org/github.com/bwmarrin/discordgo?status.svg)](https://godoc.org/github.com/bwmarrin/discordgo) [![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/bwmarrin/discordgo) [![Go report](http://goreportcard.com/badge/bwmarrin/discordgo)](http://goreportcard.com/report/bwmarrin/discordgo) [![Build Status](https://travis-ci.org/bwmarrin/discordgo.svg?branch=master)](https://travis-ci.org/bwmarrin/discordgo)
====

Discordgo is a [Go](https://golang.org/) package that provides low level 
bindings to the [Discord](https://discordapp.com/) chat client API.

* See [dgVoice](https://github.com/bwmarrin/dgvoice) for **experimental** voice
support.

Join [#go_discordgo](https://discord.gg/0SBTUU1wZTWT6sqd) Discord chat channel 
for support.

## Getting Started

### Installing

Discordgo has been tested to compile on Debian 8 (Go 1.3.3), 
FreeBSD 10 (Go 1.5.1), and Windows 7 (Go 1.5.2).

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

```sh
$ go get github.com/bwmarrin/discordgo
```

### Usage

Import the package into your project.

```go
import "github.com/bwmarrin/discordgo"
```

Construct a new Discord client which can be used to access the variety of 
Discord API functions and to set callback functions for Discord events.

```go
discord, err := discordgo.New("username", "password")
```

See Documentation and Examples below for more detailed information.


## Contributing
Contributions are very welcomed, however please follow the below guidelines.

- First open an issue describing the bug or enhancement so it can be
discussed.  
- Fork the develop branch and make your changes.  
- Try to match current naming conventions as closely as possible.  
- This package is intended to be a low level direct mapping of the Discord API 
so please avoid adding enhancements outside of that scope without first 
discussing it.
- Create a Pull Request with your changes against the develop branch.


## Documentation

**NOTICE** : This library and the Discord API are unfinished.
Because of that there may be major changes to library functions, constants,
and structures.

The Discordgo code is fairly well documented at this point and is currently
the only documentation available.  Both GoDoc and GoWalker (below) present
that information in a nice format.

- [![GoDoc](https://godoc.org/github.com/bwmarrin/discordgo?status.svg)](https://godoc.org/github.com/bwmarrin/discordgo) 
- [![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/bwmarrin/discordgo) 
- [Unofficial Discord API Documentation](https://discordapi.readthedocs.org/en/latest/) 
- Hand crafted documentation coming eventually.


## Examples / Projects using Discordgo

Below is a list of examples and other projects using Discordgo.  Please submit 
an issue if you would like your project added or removed from this list 

- [Basic - New](https://github.com/bwmarrin/discordgo/tree/develop/example/new_basic) A basic example using the easy New() helper function
- [Basic - API](https://github.com/bwmarrin/discordgo/tree/develop/example/api_basic) A basic example using the low level API functions.
- [Bruxism](https://github.com/iopred/bruxism) A chat bot for YouTube and Discord
- [GoGerard](https://github.com/GoGerard/GoGerard) A modern bot for Discord
- [Digo](https://github.com/sethdmoore/digo) A pluggable bot for your Discord server

## List & Comparison of Discord APIs

See [this chart](https://abal.moe/Discord/Libraries.html) for a feature 
comparison and list of other Discord API libraries.
