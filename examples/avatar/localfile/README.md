<img align="right" src="http://bwmarrin.github.io/discordgo/img/discordgo.png">
Avatar Local File Example
====

This example demonstrates how to utilize DiscordGo to change the account avatar using a local file inside the current working directory.

### Build

This assumes you already have a working Go environment setup and that DiscordGo is correctly installed on your system.
Change directory into the example.

```sh
cd $GOPATH/src/github.com/bwmarrin/discordgo/examples/avatar/localfile
```

```sh
go build
```

### Usage

Please place the file you wish to use as an avatar inside the directory named as ``avatar.jpg``. The filename is not important if you supply it via the commandline flag ``-f`` when starting the application.

```sh
./localfile --help
Usage of ./ocalfile:
  -e string
        Account Email
  -p string
        Account Password
  -t string
        Account Token
  -f string
  		Avatar File Name.
```

For example to start application with Token and a non-default avatar:

```sh
./localfile -t "YOUR_BOT_TOKEN" -f "./pathtoavatar.jpg"
```