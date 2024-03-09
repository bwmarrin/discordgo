<img align="right" alt="DiscordGo logo" src="/docs/img/discordgo.svg" width="200">

## DiscordGo Echo Example

This example demonstrates how to utilize DiscordGo to create a simple, 
slash commands based bot, that will echo your messages. 

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

```
Usage of echo:
  -app string
        Application ID
  -guild string
        Guild ID
  -token string
        Bot authentication token

```

Run the command below to start the bot.

```sh
./echo -guild YOUR_TESTING_GUILD -app YOUR_TESTING_APP -token YOUR_BOT_TOKEN
```
