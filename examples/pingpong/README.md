<img align="right" src="http://bwmarrin.github.io/discordgo/img/discordgo.png">
PingPong Example
====

This example demonstrates how to utilize DiscordGo to create a Ping Pong Bot.

This Bot will respond to "ping" with "Pong!" and "pong" with "Ping!".

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

```sh
go build
```

### Usage

You must authenticate using either an Authentication Token or both Email and
Password for an account.  Keep in mind official Bot accounts only support
authenticating via Token.

```
./pingpong --help
Usage of ./pingpong:
  -e string
        Account Email
  -p string
        Account Password
  -t string
        Account Token
```

The below example shows how to start the bot using an Email and Password for
authentication.

```sh
./pingpong -e EmailHere -p PasswordHere
```

The below example shows how to start the bot using the bot user's token

```sh
./pingpong -t "Bot YOUR_BOT_TOKEN"
```
