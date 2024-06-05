<img align="right" alt="DiscordGo logo" src="/docs/img/discordgo.svg" width="200">

## DiscordGo OAuth2 Example

This example demonstrates how to utilize DiscordGo to request discord tokens 
from the discord oauth2 api using a demo.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the airhorn example folder, run the command below to compile the
example.

```sh
go build
```

### Usage

```
Usage of ./oauth2 
  -id string
        Client id taken from the Settings>OAuth2 section of your application on the Discord Developer's Portal
  -secret string
        Client secret taken from the Settings>OAuth2 section of your application on the Discord Developer's Portal
  -redirect string
        Redirect URI set in Settings>OAuth2 section of your application on the Discord Developers Portal
  -scope string
        Scope to use for all OAuth2 request made with the demo server
  -port int
        Port to bind the demo server must be a number between 1 and 65535
  -cert string
        Full path to an ssl certificate if you wish the server to run in https mode
  -key string
        Full path to an ssl key if you wish the server to run in https mode
  -log bool
        Enable loggin response to stdout
```

### Setup

You must have created an application in the [Discord Developer's Portal](https://discord.com/developers/applications)
and have setup the Settings>OAuth2 section there to use the demo server.

You will need to grab both the client id and client secret from that section 
as well as register a redirect uri to use in that section.

### Http/s Mode

The demo server will run in http mode unless it is provided with an 
ssl certificate and key using the -cert and -key options.

Both should be full path names and the cert file should be the full chain
cert file usually called "fullchain.pem" and not the individual certifcate file.

### Using the demo server

The demo server must be run on the same machine as the redirect uri points to.

Assuming your redirect uri is `https://example.org/signin`:

Go to `https://example.org/signin` and you will be redirected to Discord's
OAuth2 sign in portal. After approving the sign in you will be returned to
the demo server and your access token will be shown.

Go to `https://example.org/signin?cred` and the demo server will obtain the
OAuth2 token for the application's developer. The tokens will be redacted in
the demo server's response but not the stdout log if enabled.

Go to `https://example.org/signin?token=...`, substituting a valid access token
without the `Bearer ` prefix for ..., and the demo server will display the 
current authorization info for that access token.
