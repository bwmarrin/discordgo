<img align="right" src="http://bwmarrin.github.io/discordgo/img/discordgo.png">

## DiscordGo Basic Embed Example

This example demonstrates how to utilize DiscordGo to create an embed message.

This Bot will respond to "ping" with a basic embed message.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.


From within the embed example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

This example uses bot tokens for authentication only. While user/password is
supported by DiscordGo, it is not recommended for bots.

```
./pingpong --help
Usage of ./embed:
  -t string
        Bot Token
```

The below example shows how to start the bot

```sh
./embed -t YOUR_BOT_TOKEN
Bot is now running.  Press CTRL-C to exit.
```
