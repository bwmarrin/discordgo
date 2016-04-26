<img align="right" src="http://bwmarrin.github.io/discordgo/img/discordgo.png">
Airhorn Example
====

This example demonstrates how to utilize DiscordGo to listen to an !airhorn
command in a channel and play a sound to that users current voice channel.

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

```sh
go install github.com/bwmarrin/discordgo/examples/airhorn
cd $GOPATH/bin
cp ../src/github.com/bwmarrin/discordgo/examples/airhorn/airhorn.dca .
```

### Usage

```
Usage of ./airhorn:
  -t string
        Account Token
```

The below example shows how to start the bot.

```sh
./airhorn -t <bot token>
```

### Creating sounds

Airhorn bot uses DCA files that are pre-computed files that are easy to send to Discord.

If you would like to create your own DCA files, please use either:
* [https://github.com/nstafie/dca-rs](dca-rs)
* [https://github.com/bwmarrin/dca/tree/master/cmd/dca](dca).

```sh
./dca-rs -i <input wav file> --raw > <output file>
```
