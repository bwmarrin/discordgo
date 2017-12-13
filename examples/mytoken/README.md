<img align="right" src="http://bwmarrin.github.io/discordgo/img/discordgo.png">

## DiscordGo MyToken Example

This example demonstrates how to utilize DiscordGo to login with an email and
password then to print out the Authentication Token for that user's account.

Everytime this application is run a new authentication token is generated 
for your account.  Logging you in via email and password then creating a new
token is a cpu/mem expensive task for Discord.  Because of that, it is highly
recommended to avoid doing this very often.  Please only use this once to get a 
token for your use and then always just your token.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the mytoken example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

You must authenticate using both Email and Password for an account.

```
./mytoken --help
Usage of ./mytoken:
  -e string
        Account Email
  -p string
        Account Password
```

The below example shows how to start the program using an Email and Password for
authentication.

```sh
./mytoken -e youremail@here.com -p MySecretPassword
```
