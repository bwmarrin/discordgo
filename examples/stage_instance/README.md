<img align="right" alt="DiscordGo logo" src="/docs/img/discordgo.svg" width="200">

## DiscordGo Stage Instance Example

This example demonstrates how to utilize DiscordGo to manage stage instances.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the stage_instance example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

```
Usage of stage_instance:
  -guild string
    	Test guild ID
  -stage string
    	Test stage channel ID
  -token string
    	Bot token
```

The below example shows how to start the bot from the stage_instance example folder.

```sh
./stage_instance -guild YOUR_TESTING_GUILD -stage STAGE_CHANNEL_ID -token YOUR_BOT_TOKEN
```
```
