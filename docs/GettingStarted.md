# Getting Started

This page is dedicated to helping you get started on your way to making the
next great Discord bot or client with DiscordGo. Once you've done that please
don't forget to submit it to the 
[Awesome DiscordGo](https://github.com/bwmarrin/discordgo/wiki/Awesome-DiscordGo) list :).


**First, lets cover a few topics so you can make the best choices on how to 
move forward from here.**


### Master vs Develop
**When installing DiscordGo you will need to decide if you want to use the current
master branch or the bleeding edge development branch.**

* The **master** branch represents the latest released version of DiscordGo. This
branch will always have a stable and tested version of the library. Each 
release is tagged and you can easily download a specific release and view the 
release notes on the github [releases](https://github.com/bwmarrin/discordgo/releases) 
page.

* The **develop** branch is where all development happens and almost always has
new features over the master branch.  However breaking changes are frequently
added the develop branch and sometimes bugs are introduced.  Bugs get fixed
and the breaking changes get documented before pushing to master.  

*So, what should you use?*

Due to the how frequently the Discord API is changing there is a high chance
that the *master* branch may be lacking important features.  Because of that, if
you can accept the constant changing nature of the *develop* branch and the 
chance that it may occasionally contain bugs then it is the recommended branch 
to use.  Otherwise, if you want to tail behind development slightly and have a 
more stable package with documented releases then please use the *master* 
branch instead.


### Client vs Bot

You probably already know the answer to this but now is a good time to decide
if your goal is to write a client application or a bot.  DiscordGo aims to fully
support both client applications and bots but there are some differences 
between the two that you should understand.

#### Client Application
A client application is a program that is intended to be used by a normal user 
as a replacement for the official clients that Discord provides. An example of
this would be a terminal client used to read and send messages with your normal
user account or possibly a new desktop client that provides a different set of
features than the official desktop client that Discord already provides.

Client applications work with normal user accounts and you can login with an
email address and password or a special authentication token.  However, normal
user accounts are not allowed to perform any type of automation and doing so can
cause the account to be banned from Discord. Also normal user accounts do not 
support multi-server voice connections and some other features that are 
exclusive to Bot accounts only.

To create a new user account (if you have not done so already) visit the 
[Discord](https://discordapp.com/) website and click on the 
**Try Discord Now, It's Free** button then follow the steps to setup your
new account.


#### Bot Application
A bot application is a special program that interacts with the Discord servers
to perform some form of automation or provide some type of service.  Examples 
are things like number trivia games, music streaming, channel moderation, 
sending reminders, playing loud airhorn sounds, comic generators, YouTube 
integration, Twitch integration.. You're *almost* only limited by your imagination.

Bot applications require the use of a special Bot account.  These accounts are
tied to your personal user account. Bot accounts cannot login with the normal
user clients and they cannot join servers the same way a user does. They do not 
have access to some user client specific features however they gain access to
many Bot specific features.

To create a new bot account first create yourself a normal user account on 
Discord then visit the [My Applications](https://discordapp.com/developers/applications/me)
page and click on the **New Application** box.  Follow the prompts from there
to finish creating your account.


**More information about Bots vs Client accounts can be found [here](https://discordapp.com/developers/docs/topics/oauth2#bot-vs-user-accounts)**

# Requirements

DiscordGo requires Go version 1.4 or higher.  It has been tested to compile and
run successfully on Debian Linux 8, FreeBSD 10, and Windows 7.  It is expected 
that it should work anywhere Go 1.4 or higher works. If you run into problems
please let us know :)

You must already have a working Go environment setup to use DiscordGo.  If you 
are new to Go and have not yet installed and tested it on your computer then 
please visit [this page](https://golang.org/doc/install) first then I highly
recommend you walk though [A Tour of Go](https://tour.golang.org/welcome/1) to
help get your familiar with the Go language.  Also checkout the relevent Go plugin 
for your editor - they are hugely helpful when developing Go code.

* Vim - [vim-go](https://github.com/fatih/vim-go)
* Sublime - [GoSublime](https://github.com/DisposaBoy/GoSublime)
* Atom - [go-plus](https://atom.io/packages/go-plus)
* Visual Studio - [vscode-go](https://github.com/Microsoft/vscode-go)


# Install DiscordGo

Like any other Go package the fist step is to `go get` the package.  This will
always pull the latest released version from the master branch. Then run 
`go install` to compile and install the libraries on your system.

#### Linux/BSD

Run go get to download the package to your GOPATH/src folder.

```sh
go get github.com/bwmarrin/discordgo
```

If you want to use the develop branch, follow these steps next.

```sh
cd $GOPATH/src/github.com/bwmarrin/discordgo
git checkout develop
```

Finally, compile and install the package into the GOPATH/pkg folder. This isn't
absolutely required but doing this will allow the Go plugin for your editor to
provide autocomplete for all DiscordGo functions.

```sh
cd $GOPATH/src/github.com/bwmarrin/discordgo
go install
```

#### Windows
Placeholder.


# Next...
More coming soon.
