# Discordgo

Discord API for Golang

This is my first Golang project and it is ~~probably~~ ~~not~~ ~~maybe even~~ 
barely suitable for use :)

Everything here so far is likely to change as I learn Golang better and refine 
the API names and such. Because of that I do not yet recommend this for use 
with anything super important :)

Initially my goal was to write a chatbot and I started working with 
https://github.com/Xackery/discord as my API.  But that code didn't work 100% 
at the time. So I started slowly making changes to it and eventually ended up 
with something entirely different. Anyhow, credit goes to 
https://github.com/Xackery/discord for getting me started.

If you're looking for a more functional Discord API for Golang check out 
https://github.com/gdraynz/go-discord which I recently found.  It's much more 
complete and will likely help me learn how to improve what I have here.

# What Works

Low level functions exist for the core REST API and Websocket API.

* Login to Discord
* Get User information for a given user.
* Get Private Channels (used for Private Messages) for a given user.
* Get Servers for a given user.
* Get Members of a given Server
* Get Channels for a given Server
* Get Messages for a given Channel
* Send Messages to a given Channel
* Start a Websocket connection and listen for and handle events.
* Logout from Discord.


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
