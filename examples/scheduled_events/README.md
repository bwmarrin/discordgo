<img align="right" alt="DiscordGo logo" src="/docs/img/discordgo.svg" width="200">

## DiscordGo Scheduled Events Example

This example demonstrates how to utilize DiscordGo to manage scheduled events
in a guild.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the scheduled_events example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

```
Usage of scheduled_events:
  -guild string
    	Test guild ID
  -token string
    	Bot token
  -voice string
    	Test voice channel ID
```

The below example shows how to start the bot from the scheduled_events example folder.

```sh
./scheduled_events -guild YOUR_TESTING_GUILD -token YOUR_BOT_TOKEN -voice YOUR_TESTING_CHANNEL
```
