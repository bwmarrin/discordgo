<img align="right" src="http://bwmarrin.github.io/discordgo/img/discordgo.png">
AppMaker Example
====

This example demonstrates how to utilize DiscordGo to create Bot Applications.

You can create a new bot account or convert an existing normal user account into
a bot account with this tool.  You can also view the list of applications you 
have on your account and delete applications.

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
  -c string
        Token of account to convert.
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

* If you provide a **-c** flag with a valid user token then than user account
will be converted into a Bot account instead of creating a new Bot account for
an application.


Below example will create a new Bot Application under the given Email/Password 
account. The Bot will be named **DiscordGoRocks**

```sh
./appmaker -e Email -p Password -a DiscordGoRocks
```
