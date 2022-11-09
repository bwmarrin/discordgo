<img align="right" alt="DiscordGo logo" src="/docs/img/discordgo.svg" width="200">

## DiscordGo Slash Command Autocomplete Option Example

This example demonstrates how to utilize DiscordGo to create and use
autocomplete options in Slash Commands. As this example uses interactions,
slash commands and slash command options, it is recommended to read
`slash_commands` example before proceeding.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the autocomplete example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

```
Usage of autocomplete:
  -guild string
    	Test guild ID. If not passed - bot registers commands globally
  -rmcmd
    	Whether to remove all commands after shutting down (default true)
  -token string
    	Bot access token
```

The below example shows how to start the bot from the autocomplete example folder.

```sh
./autocomplete -guild YOUR_TESTING_GUILD -token YOUR_BOT_TOKEN
```
