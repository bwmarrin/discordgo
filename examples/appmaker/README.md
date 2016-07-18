<img align="right" src="http://bwmarrin.github.io/discordgo/img/discordgo.png">
AppMaker Example
====

This example demonstrates how to utilize DiscordGo to create Bot Applications.

You can create a new bot account, view the list of applications you have, and
delete applications.

### Build

This assumes you already have a working Go environment setup and that 
DiscordGo is correctly installed on your system.

```sh
go build
```

### Usage

```
Usage of ./appmaker:
  -a string
        App/Bot Name
  -d string
        Application ID to delete
  -e string
        Account Email
  -l    List Applications Only
  -p string
        Account Password
  -t string
        Account Token
```

* Account Email and Password or Token are required.  The account provided with
these fields will be the "owner" of any bot applications created.

* If you provide the **-l** flag than appmaker will only display a list of 
applications on the provided account.

* If you provide a **-d** flag with a valid application ID then that application
will be deleted.

Below example will create a new Bot Application under the given Email/Password 
account. The Bot will be named **DiscordGoRocks**

```sh
./appmaker -e Email -p Password -a DiscordGoRocks
```
