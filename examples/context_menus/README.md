<img align="right" alt="DiscordGo logo" src="/docs/img/discordgo.svg" width="200">

## DiscordGo Context Menu Commands Example

This example demonstrates how to utilize DiscordGo to create and use context
menu commands. This example heavily relies on `slash_commands` example in
command handling and registration, therefore it is recommended to be read
before proceeding.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the context_menus example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

```
Usage of context_menus:
  -app string
    	Application ID
  -cleanup
    	Cleanup of commands (default true)
  -guild string
    	Test guild ID
  -token string
    	Bot access token
```

The below example shows how to start the bot from the context_menus example folder.

```sh
./context_menus -app YOUR_APPLICATION_ID -guild YOUR_TESTING_GUILD -token YOUR_BOT_TOKEN
```
