<img align="right" alt="DiscordGo logo" src="/docs/img/discordgo.svg" width="200">

## DiscordGo Airhorn Example

This example demonstrates how to utilize DiscordGo to listen for an !airhorn
command in a channel and then play a sound to that user's current voice channel.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the airhorn example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

```
Usage of ./airhorn:
  -t string
        Bot Token
```

The below example shows how to start the bot from the airhorn example folder.

```sh
./airhorn -t YOUR_BOT_TOKEN
```

### Creating sounds

Airhorn bot uses [DCA](https://github.com/bwmarrin/dca) files, which are 
pre-computed files that are easy to send to Discord.


See the below example of creating a DCA file from a MP3 file.  This also works
with WAV, FLAC, and many other file formats. Of course, you will need to [install](https://github.com/bwmarrin/dca/tree/master/cmd/dca#Getting-Started)
FFmpeg and the DCA CLI first.

```sh
ffmpeg -i test.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca > test.dca
```
