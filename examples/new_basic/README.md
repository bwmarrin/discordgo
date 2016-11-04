<img align="right" src="http://bwmarrin.github.io/discordgo/img/discordgo.png">
Basic New Example
====

This example demonstrates how to utilize DiscordGo to connect to Discord
and print out all received chat messages.

This example uses the high level New() helper function to connect to Discord.

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
./new_basic --help
Usage of ./new_basic:
  -t string
        Bot Token
```

The below example shows how to start the bot

```sh
./new_basic -t YOUR_BOT_TOKEN
```
