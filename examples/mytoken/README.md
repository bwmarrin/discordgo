<img align="right" src="http://bwmarrin.github.io/discordgo/img/discordgo.png">
MyToken Example
====

This example demonstrates how to utilize DiscordGo to print out the 
Authentication Token for a given user account.

### Build

This assumes you already have a working Go environment setup and that 
DiscordGo is correctly installed on your system.

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
./mytoken -e EmailHere -p PasswordHere
```
