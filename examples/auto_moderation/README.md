<img align="right" alt="DiscordGo logo" src="/docs/img/discordgo.svg" width="200">

## DiscordGo Auto Moderation Example

This example demonstrates how to utilize DiscordGo to manage auto moderation
rules and triggers.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the auto_moderation example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

```
Usage of auto_moderation:
  -channel string
    	ID of the testing channel
  -guild string
    	ID of the testing guild
  -token string
    	Bot authorization token
```

The below example shows how to start the bot from the auto_moderation example folder.

```sh
./auto_moderation -channel YOUR_TESTING_CHANNEL -guild YOUR_TESTING_GUILD -token YOUR_BOT_TOKEN
```
