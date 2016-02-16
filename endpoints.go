// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains variables for all known Discord end points.  All functions
// throughout the Discordgo package use these variables for all connections
// to Discord.  These are all exported and you may modify them if needed.

package discordgo

// Known Discord API Endpoints.
var (
	STATUS      = "https://status.discordapp.com/api/v2/"
	SM          = STATUS + "scheduled-maintenances/"
	SM_ACTIVE   = SM + "active.json"
	SM_UPCOMING = SM + "upcoming.json"

	DISCORD  = "https://discordapp.com" // TODO consider removing
	API      = DISCORD + "/api/"
	GUILDS   = API + "guilds/"
	CHANNELS = API + "channels/"
	USERS    = API + "users/"
	GATEWAY  = API + "gateway"

	AUTH            = API + "auth/"
	LOGIN           = AUTH + "login"
	LOGOUT          = AUTH + "logout"
	VERIFY          = AUTH + "verify"
	VERIFY_RESEND   = AUTH + "verify/resend"
	FORGOT_PASSWORD = AUTH + "forgot"
	RESET_PASSWORD  = AUTH + "reset"
	REGISTER        = AUTH + "register"

	VOICE         = API + "/voice/"
	VOICE_REGIONS = VOICE + "regions"
	VOICE_ICE     = VOICE + "ice"

	TUTORIAL            = API + "tutorial/"
	TUTORIAL_INDICATORS = TUTORIAL + "indicators"

	TRACK        = API + "track"
	SSO          = API + "sso"
	REPORT       = API + "report"
	INTEGRATIONS = API + "integrations"

	USER             = func(uID string) string { return USERS + uID }
	USER_AVATAR      = func(uID, aID string) string { return USERS + uID + "/avatars/" + aID + ".jpg" }
	USER_SETTINGS    = func(uID string) string { return USERS + uID + "/settings" }
	USER_GUILDS      = func(uID string) string { return USERS + uID + "/guilds" }
	USER_GUILD       = func(uID, gID string) string { return USERS + uID + "/guilds/" + gID }
	USER_CHANNELS    = func(uID string) string { return USERS + uID + "/channels" }
	USER_DEVICES     = func(uID string) string { return USERS + uID + "/devices" }
	USER_CONNECTIONS = func(uID string) string { return USERS + uID + "/connections" }

	GUILD              = func(gID string) string { return GUILDS + gID }
	GUILD_INIVTES      = func(gID string) string { return GUILDS + gID + "/invites" }
	GUILD_CHANNELS     = func(gID string) string { return GUILDS + gID + "/channels" }
	GUILD_MEMBERS      = func(gID string) string { return GUILDS + gID + "/members" }
	GUILD_MEMBER       = func(gID, uID string) string { return GUILDS + gID + "/members/" + uID }
	GUILD_BANS         = func(gID string) string { return GUILDS + gID + "/bans" }
	GUILD_BAN          = func(gID, uID string) string { return GUILDS + gID + "/bans/" + uID }
	GUILD_INTEGRATIONS = func(gID string) string { return GUILDS + gID + "/integrations" }
	GUILD_ROLES        = func(gID string) string { return GUILDS + gID + "/roles" }
	GUILD_ROLE         = func(gID, rID string) string { return GUILDS + gID + "/roles/" + rID }
	GUILD_INVITES      = func(gID string) string { return GUILDS + gID + "/invites" }
	GUILD_EMBED        = func(gID string) string { return GUILDS + gID + "/embed" }
	GUILD_PRUNE        = func(gID string) string { return GUILDS + gID + "/prune" }
	GUILD_ICON         = func(gID, hash string) string { return GUILDS + gID + "/icons/" + hash + ".jpg" }
	GUILD_SPLASH       = func(gID, hash string) string { return GUILDS + gID + "/splashes/" + hash + ".jpg" }

	CHANNEL             = func(cID string) string { return CHANNELS + cID }
	CHANNEL_PERMISSIONS = func(cID string) string { return CHANNELS + cID + "/permissions" }
	CHANNEL_PERMISSION  = func(cID, tID string) string { return CHANNELS + cID + "/permissions/" + tID }
	CHANNEL_INVITES     = func(cID string) string { return CHANNELS + cID + "/invites" }
	CHANNEL_TYPING      = func(cID string) string { return CHANNELS + cID + "/typing" }
	CHANNEL_MESSAGES    = func(cID string) string { return CHANNELS + cID + "/messages" }
	CHANNEL_MESSAGE     = func(cID, mID string) string { return CHANNELS + cID + "/messages/" + mID }
	CHANNEL_MESSAGE_ACK = func(cID, mID string) string { return CHANNELS + cID + "/messages/" + mID + "/ack" }

	INVITE = func(iID string) string { return API + "invite/" + iID }

	INTEGRATIONS_JOIN = func(iID string) string { return API + "integrations/" + iID + "/join" }

	EMOJI = func(eID string) string { return API + "emojis/" + eID + ".png" }
)
