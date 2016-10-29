<img align="right" src="http://bwmarrin.github.io/discordgo/img/discordgo.png">
PingPong Example
====

This example demonstrates how to utilize DiscordGo to create a Ping Pong Bot.

This Bot will respond to "ping" with "Pong!" and "pong" with "Ping!".

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

```sh
go build
```

### Usage

This example uses bot tokens for authentication only.
While user/password is supported by DiscordGo, it is not recommended.

```
./pingpong --help
Usage of ./pingpong:
  -t string
        Bot Token
```

The below example shows how to start the bot

```sh
./pingpong -t YOUR_BOT_TOKEN
```
