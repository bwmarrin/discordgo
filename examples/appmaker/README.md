<img align="right" src="http://bwmarrin.github.io/discordgo/img/discordgo.png">

## DiscordGo AppMaker Example

This example demonstrates how to utilize DiscordGo to create, view, and delete
Bot Applications on your account.

These tasks are normally accomplished from the 
[Discord Developers](https://discordapp.com/developers/applications/me) site.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the appmaker example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

This example only uses authentication tokens for authentication. While 
user email/password is supported by DiscordGo, it is not recommended.

```
./appmaker --help
Usage of ./appmaker:
  -d string
        Application ID to delete
  -l    List Applications Only
  -n string
        Name to give App/Bot
  -t string
        Owner Account Token
```

* Account Token is required.  The account will be the "owner" of any bot 
applications created.

* If you provide the **-l** flag than appmaker will only display a list of 
applications on the provided account.

* If you provide a **-d** flag with a valid application ID then that application
will be deleted.

Below example will create a new Bot Application under the given account.
The Bot will be named **DiscordGoRocks**

```sh
./appmaker -t YOUR_USER_TOKEN -n DiscordGoRocks
```
