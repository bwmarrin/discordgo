<img align="right" alt="DiscordGo logo" src="/docs/img/discordgo.svg" width="200">

## DiscordGo Voice Receive Example

This example experiments with receiving voice data from Discord. It joins
a specified voice channel, listens for 10 seconds and saves .ogg files for each
SSRC that it finds in the channel. An exercise left to the reader is to translate
these SSRCs to user IDs; see speaking update events for this information. :)

This example makes heavy use of the [Pion](https://github.com/pion) family of libraries.
Go check them out for anything to do with voice, video or WebRTC; it's a great
group of people maintaining the project!

Please note that voice receive is **not** officially supported, any may break
at essentially any time (and has in the past). This code works at the time of
its writing, but YMMV in the future.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

To build, make sure that modules are enabled, and run:

```sh
go build
```

### Usage

Three flags are required: the bot's token, the guild ID containing the voice channel to join, and the ID of the voice channel to join.

```sh
./voice_receive -t MY_TOKEN -g 1234123412341234 -c 5678567856785678
```
