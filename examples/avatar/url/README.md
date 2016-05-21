<img align="right" src="http://bwmarrin.github.io/discordgo/img/discordgo.png">
Avatar Url Example
====

This example demonstrates how to utilize DiscordGo to change the account avatar using a remote url provided via a commandline flag.

### Build

This assumes you already have a working Go environment setup and that DiscordGo is correctly installed on your system.
Change directory into the example.

```sh
cd $GOPATH/src/github.com/bwmarrin/discordgo/examples/avatar/url
```

```sh
go build
```

### Usage

Please place the file you wish to use as an avatar inside the directory named as ``avatar.jpg``. The filename is not important if you supply it via the commandline flag ``-l`` when starting the application. If the flag is not specified avatar is set to DiscordGo Logo.

```sh
./url --help
Usage of ./url:
  -e string
        Account Email
  -p string
        Account Password
  -t string
        Account Token
  -l string
  		Link to the avatar image.
```

For example to start application with Token and a non-default avatar:

```sh
./url -t "YOUR_BOT_TOKEN" -l "http://bwmarrin.github.io/discordgo/img/discordgo.png"
```