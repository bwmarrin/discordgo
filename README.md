<img align="right" src="http://bwmarrin.github.io/discordgo/img/discordgo.png">
Discordgo 
====
[![GoDoc](https://godoc.org/github.com/bwmarrin/discordgo?status.svg)](https://godoc.org/github.com/bwmarrin/discordgo) [![Go report](http://goreportcard.com/badge/bwmarrin/discordgo)](http://goreportcard.com/report/bwmarrin/discordgo) [![Build Status](https://travis-ci.org/bwmarrin/discordgo.svg?branch=master)](https://travis-ci.org/bwmarrin/discordgo)

Discordgo is a [Go](https://golang.org/) package that provides low level 
bindings to the [Discord](https://discordapp.com/) chat client API. Discordgo 
has nearly complete support for all of the Discord JSON-API endpoints, websocket
interface, and voice interface.

* See [dgVoice](https://github.com/bwmarrin/dgvoice) package to extend Discordgo
with additional voice helper functions and features.

* See [dca](https://github.com/bwmarrin/dca) for an **experimental** stand alone
tool that wraps `ffmpeg` to create opus encoded audio appropriate for use with
Discord (and Discordgo)

Join [#go_discordgo](https://discord.gg/0SBTUU1wZTWT6sqd) Discord chat channel 
for support.

## Getting Started

### master vs develop Branch
* The master branch represents the latest released version of Discordgo.  This
branch will always have a stable and tested version of the library. Each release
is tagged and you can easily download a specific release and view release notes
on the github [releases](https://github.com/bwmarrin/discordgo/releases) page.

* The develop branch is where all development happens and almost always has
new features over the master branch.  However breaking changes are frequently
added to develop and even sometimes bugs are introduced.  Bugs get fixed and 
the breaking changes get documented before pushing to master.  

*So, what should you use?*

If you can accept the constant changing nature of *develop* then it is the 
recommended branch to use.  Otherwise, if you want to tail behind development
slightly and have a more stable package with documented releases then use *master*

### Installing

Discordgo has been tested to compile on Debian 8 (Go 1.3.3), 
FreeBSD 10 (Go 1.5.1), and Windows 7 (Go 1.5.2).

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

`go get` *will always pull the latest released version from the master branch.*

```sh
go get github.com/bwmarrin/discordgo
```

If you want to use the develop branch, follow these steps next.

```sh
cd $GOPATH/src/github.com/bwmarrin/discordgo
git checkout develop
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

### Troubleshooting

If your go environment is brand new, you may encounter a 'cannot find package' 
error when building projects that import discordgo. If go can't find it, tell
go to `go get` it.

For example, if you see this in your build output:

```sh
...cannot find package "golang.org/x/crypto/nacl/secretbox"...
```

Run this to fix it:
```sh
go get golang.org/x/crypto/nacl/secretbox
```

See Documentation and Examples below for more detailed information.


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


## Examples

Below is a list of examples and other projects using Discordgo.  Please submit 
an issue if you would like your project added or removed from this list 

- [Basic - New](https://github.com/bwmarrin/discordgo/tree/develop/examples/new_basic) A basic example using the easy New() helper function
- [Basic - API](https://github.com/bwmarrin/discordgo/tree/develop/examples/api_basic) A basic example using the low level API functions.
- [Bruxism](https://github.com/iopred/bruxism) A chat bot for YouTube and Discord
- [GoGerard](https://github.com/GoGerard/GoGerard) A modern bot for Discord
- [Digo](https://github.com/sethdmoore/digo) A pluggable bot for your Discord server

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


## List of Discord APIs

See [this chart](https://abal.moe/Discord/Libraries.html) for a feature 
comparison and list of other Discord API libraries.
