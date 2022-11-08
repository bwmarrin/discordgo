<img align="right" alt="DiscordGo logo" src="/docs/img/discordgo.svg" width="200">

## DiscordGo Components Example

This example demonstrates how to utilize DiscordGo to create and use message
components, such as buttons and select menus. For usage of the text input
component and modals, please refer to the `modals` example.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the components example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

```
Usage of components:
  -app string
    	Application ID
  -guild string
    	Test guild ID
  -token string
    	Bot access token
```

The below example shows how to start the bot from the components example folder.

```sh
./components -app YOUR_APPLICATION_ID -guild YOUR_TESTING_GUILD -token YOUR_BOT_TOKEN
```
