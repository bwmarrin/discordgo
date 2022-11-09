<img align="right" alt="DiscordGo logo" src="/docs/img/discordgo.svg" width="200">

## DiscordGo Modals Example

This example demonstrates how to utilize DiscordGo to send and process text
inputs in modals. If you have not read `slash_commands` and `components`
examples yet it is recommended to do so before proceesing. As this example
is built using interactions and Slash Commands.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the modals example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

```
Usage of modals:
  -app string
    	Application ID
  -cleanup
    	Cleanup of commands (default true)
  -guild string
    	Test guild ID
  -results string
    	Channel where send survey results to
  -token string
    	Bot access token
```

The below example shows how to start the bot from the modals example folder.

```sh
./modals -app YOUR_APPLICATION_ID -guild YOUR_TESTING_GUILD -results YOUR_TESTING_CHANNEL -token YOUR_BOT_TOKEN
```
