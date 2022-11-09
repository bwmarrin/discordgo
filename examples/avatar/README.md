<img align="right" alt="DiscordGo logo" src="/docs/img/discordgo.svg" width="200">

## DiscordGo Avatar Example

This example demonstrates how to utilize DiscordGo to change the avatar for
a Discord account.  This example works both with a local file or the URL of
an image.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the avatar example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

This example uses bot tokens for authentication only. While email/password is 
supported by DiscordGo, it is not recommended to use them.

```
./avatar --help
Usage of ./avatar:
  -f string
        Avatar File Name
  -t string
        Bot Token
  -u string
        URL to the avatar image
```

The below example shows how to set your Avatar from a local file.

```sh
./avatar -t TOKEN -f avatar.png
```
The below example shows how to set your Avatar from a URL.

```sh
./avatar -t TOKEN -u https://raw.githubusercontent.com/bwmarrin/discordgo/master/docs/img/discordgo.svg
```
