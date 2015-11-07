# Discordgo
Golang Discord API.

This is my first Golang project and it is ~~probably~~ ~~not~~ maybe even suitable for use :)

Everything here so far is likely to change as I learn Golang better and refine the API names and such.

Initially my goal was to write a chatbot and I started working with https://github.com/Xackery/discord as my API.  But that code doesn't work 100% and so I started slowly making changes to it.  Anyhow, credit goes to https://github.com/Xackery/discord for getting me started.

If you're looking for a functional Discord API for Golang check out https://github.com/gdraynz/go-discord which I recently found.  It's much more complete and will likely help me learn how to improve what I have here.

# What Works
Right now I'm focusing on the REST API and have not done any Websockets work.  You can do the following things using the client.go functions.

* Login to Discord
* Get User information for a given user.
* Get Private Channels (used for Private Messages) for a given user.
* Get Servers for a given user.
* Get Members of a given Server
* Get Channels for a given Server
* Get Messages for a given Channel
* Send Messages to a given Channel
* Logout from Discord.

All the code in the other files such as discord.go, session.go, etc are a playground where I'm working to provide another and easier way to access the API.

You can look at the demo.go example file to see all of the client.go functions in use.



# Other Discord APIs
- [go-discord](https://github.com/gdraynz/go-discord)
- [discord-go](https://github.com/Xackery/discord)
- [discord.py](https://github.com/Rapptz/discord.py)
- [discord.js](https://github.com/discord-js/discord.js)
- [discord.io](https://github.com/izy521/discord.io)
- [Discord.NET](https://github.com/RogueException/Discord.Net)
- [DiscordSharp](https://github.com/Luigifan/DiscordSharp)
- [Discord4J](https://github.com/knobody/Discord4J)
- [discordrb](https://github.com/meew0/discordrb)

