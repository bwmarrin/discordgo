# Discordgo

A Discord API for Golang

Discordgo provides an almost complete low-level Golang interface to the Discord
REST and Websocket API layers.  The majority of both of these interfaces are
complete and I should have the remaining functions finished soon.

At this point Discordgo is suitable for use with most projects including bots 
or clients.  The function naming conventions and usage style should not change 
in the future.  Function names are based primarily on the naming used by Discord 
within their API calls.  Should Discord change their naming then Discordgo will 
be updated to match it.

Special thanks goes to both the below projects who helped me get started with
this project.  If you're looking for alternative Golang interfaces to Discord
please check both of these out.

* https://github.com/gdraynz/go-discord
* https://github.com/Xackery/discord

# What Works

Low level functions exist for the majority of the REST and Websocket API.

* Login/Logout
* Open/Close Websocket and listen for events.
* Accept/Create/Delete Invites
* Get User details (Name, ID, Settings, etc)
* List/Create User Channels (Private Message Channels)
* List/Create Guilds
* List/Create Guild Channels
* List Guild Members
* Receive/Send Messages to Channels

# What's Left

* Permissions related functions.
* Editing User Profile settings
* Voice Channel support.
* Functions for Maintenance Status

# Other Discord APIs

- [discord.py](https://github.com/Rapptz/discord.py)
- [discord.js](https://github.com/discord-js/discord.js)
- [discord.io](https://github.com/izy521/discord.io)
- [Discord.NET](https://github.com/RogueException/Discord.Net)
- [DiscordSharp](https://github.com/Luigifan/DiscordSharp)
- [Discord4J](https://github.com/knobody/Discord4J)
- [discordrb](https://github.com/meew0/discordrb)
