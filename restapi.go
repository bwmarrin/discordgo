/******************************************************************************
 * A Discord API for Golang.
 * See discord.go for more information.
 *
 * This file contains functions for interacting with the Discord HTTP REST API
 * at the lowest level.
 */

package discordgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Constants of all known Discord API Endpoints
// Please let me know if you know of any others.
const (
	STATUS      = "https://status.discordapp.com/api/v2/"
	SM          = STATUS + "scheduled-maintenances/"
	SM_ACTIVE   = SM + "active.json"
	SM_UPCOMING = SM + "upcoming.json"

	DISCORD  = "http://discordapp.com" // TODO consider removing
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
)

// Almost like the constants above :) Except can't be constants
var (
	USER             = func(uId string) string { return USERS + uId }
	USER_AVATAR      = func(uId, aId string) string { return USERS + uId + "/avatars/" + aId + ".jpg" }
	USER_SETTINGS    = func(uId string) string { return USERS + uId + "/settings" }
	USER_GUILDS      = func(uId string) string { return USERS + uId + "/guilds" }
	USER_CHANNELS    = func(uId string) string { return USERS + uId + "/channels" }
	USER_DEVICES     = func(uId string) string { return USERS + uId + "/devices" }
	USER_CONNECTIONS = func(uId string) string { return USERS + uId + "/connections" }

	GUILD              = func(gId string) string { return GUILDS + gId }
	GUILD_INIVTES      = func(gId string) string { return GUILDS + gId + "/invites" }
	GUILD_CHANNELS     = func(gId string) string { return GUILDS + gId + "/channels" }
	GUILD_MEMBERS      = func(gId string) string { return GUILDS + gId + "/members" }
	GUILD_MEMBER_DEL   = func(gId, uId string) string { return GUILDS + gId + "/members/" + uId }
	GUILD_BANS         = func(gId string) string { return GUILDS + gId + "/bans" }
	GUILD_BAN          = func(gId, uId string) string { return GUILDS + gId + "/bans/" + uId }
	GUILD_INTEGRATIONS = func(gId string) string { return GUILDS + gId + "/integrations" }
	GUILD_ROLES        = func(gId string) string { return GUILDS + gId + "/roles" }
	GUILD_INVITES      = func(gId string) string { return GUILDS + gId + "/invites" }
	GUILD_EMBED        = func(gId string) string { return GUILDS + gId + "/embed" }
	GUILD_PRUNE        = func(gId string) string { return GUILDS + gId + "/prune" }
	GUILD_ICON         = func(gId, hash string) string { return GUILDS + gId + "/icons/" + hash + ".jpg" }

	CHANNEL             = func(cId string) string { return CHANNELS + cId }
	CHANNEL_PERMISSIONS = func(cId string) string { return CHANNELS + cId + "/permissions" }
	CHANNEL_INVITES     = func(cId string) string { return CHANNELS + cId + "/invites" }
	CHANNEL_TYPING      = func(cId string) string { return CHANNELS + cId + "/typing" }
	CHANNEL_MESSAGES    = func(cId string) string { return CHANNELS + cId + "/messages" }
	CHANNEL_MESSAGE     = func(cId, mId string) string { return CHANNELS + cId + "/messages/" + mId }
	CHANNEL_MESSAGE_ACK = func(cId, mId string) string { return CHANNELS + cId + "/messages/" + mId + "/ack" }

	INVITE = func(iId string) string { return API + "invite/" + iId }

	INTEGRATIONS_JOIN = func(iId string) string { return API + "integrations/" + iId + "/join" }
)

// Request makes a (GET/POST/?) Requests to Discord REST API.
// All the other Discord REST Calls in this file use this function.
func (s *Session) Request(method, urlStr, body string) (response []byte, err error) {

	if s.Debug {
		fmt.Println("REQUEST  :: " + method + " " + urlStr + "\n" + body)
	}

	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return
	}

	// Not used on initial login..
	if s.Token != "" {
		req.Header.Set("authorization", s.Token)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: (20 * time.Second)}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(response)))
		return
	}

	if s.Debug {
		printJSON(response)
	}
	return
}

/***************************************************************************************************
 * Functions specific to this session.
 */

// Login asks the Discord server for an authentication token
func (s *Session) Login(email string, password string) (token string, err error) {

	response, err := s.Request("POST", LOGIN, fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password))

	var temp map[string]interface{}
	err = json.Unmarshal(response, &temp)
	token = temp["token"].(string)
	return
}

// Logout sends a logout request to Discord.
// This does not seem to actually invalidate the token.  So you can still
// make API calls even after a Logout.  So, it seems almost pointless to
// even use.
func (s *Session) Logout() (err error) {

	_, err = s.Request("POST", LOGOUT, fmt.Sprintf(`{"token": "%s"}`, s.Token))
	return
}

/***************************************************************************************************
 * Functions related to a specific user
 */

// User returns the user details of the given userId
// userId    : A user Id or "@me" which is a shortcut of current user ID
func (s *Session) User(userId string) (st User, err error) {

	body, err := s.Request("GET", USER(userId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// UserAvatar returns a ?? of a users Avatar
// userId    : A user Id or "@me" which is a shortcut of current user ID
func (s *Session) UserAvatar(userId string) (st User, err error) {

	u, err := s.User(userId)
	_, err = s.Request("GET", USER_AVATAR(userId, u.Avatar), ``)
	// TODO need to figure out how to handle returning a file
	return
}

// UserSettings returns the settings for a given user
// userId    : A user Id or "@me" which is a shortcut of current user ID
// This seems to only return a result for "@me"
func (s *Session) UserSettings(userId string) (st Settings, err error) {

	body, err := s.Request("GET", USER_SETTINGS(userId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// UserChannels returns an array of Channel structures for all private
// channels for a user
// userId    : A user Id or "@me" which is a shortcut of current user ID
func (s *Session) UserChannels(userId string) (st []Channel, err error) {

	body, err := s.Request("GET", USER_CHANNELS(userId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// UserChannelCreate creates a new User (Private) Channel with another User
// userId      : A user Id or "@me" which is a shortcut of current user ID
// recipientId : A user Id for the user to which this channel is opened with.
func (s *Session) UserChannelCreate(userId, recipientId string) (st []Channel, err error) {

	body, err := s.Request(
		"POST",
		USER_CHANNELS(userId),
		fmt.Sprintf(`{"recipient_id": "%s"}`, recipientId))

	err = json.Unmarshal(body, &st)
	return
}

// UserGuilds returns an array of Guild structures for all guilds for a given user
// userId    : A user Id or "@me" which is a shortcut of current user ID
func (s *Session) UserGuilds(userId string) (st []Guild, err error) {

	body, err := s.Request("GET", USER_GUILDS(userId), ``)
	err = json.Unmarshal(body, &st)
	return
}

/***************************************************************************************************
 * Functions related to a specific guild
 */

// Guild returns a Guild structure of a specific Guild.
// guildId   : The ID of a Guild
func (s *Session) Guild(guildId string) (st []Guild, err error) {

	body, err := s.Request("GET", GUILD(guildId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// GuildCreate creates a new Guild
// name      : A name for the Guild (2-100 characters)
func (s *Session) GuildCreate(name string) (st []Guild, err error) {

	body, err := s.Request("POST", GUILDS, fmt.Sprintf(`{"name":"%s"}`, name))
	err = json.Unmarshal(body, &st)
	return
}

// GuildEdit edits a new Guild
// guildId   : The ID of a Guild
// name      : A name for the Guild (2-100 characters)
func (s *Session) GuildEdit(guildId, name string) (st []Guild, err error) {

	body, err := s.Request("POST", GUILD(guildId), fmt.Sprintf(`{"name":"%s"}`, name))
	err = json.Unmarshal(body, &st)
	return
}

// GuildDelete deletes or leaves a Guild.
// guildId   : The ID of a Guild
func (s *Session) GuildDelete(guildId string) (st []Guild, err error) {

	body, err := s.Request("DELETE", GUILD(guildId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// GuildBans returns an array of User structures for all bans of a
// given guild.
// guildId   : The ID of a Guild.
func (s *Session) GuildBans(guildId string) (st []User, err error) {

	body, err := s.Request("GET", GUILD_BANS(guildId), ``)
	err = json.Unmarshal(body, &st)

	return
}

// GuildBanAdd bans the given user from the given guild.
// guildId   : The ID of a Guild.
// userId    : The ID of a User
func (s *Session) GuildBanAdd(guildId, userId string) (err error) {

	_, err = s.Request("PUT", GUILD_BAN(guildId, userId), ``)
	return
}

// GuildBanDelete removes the given user from the guild bans
// guildId   : The ID of a Guild.
// userId    : The ID of a User
func (s *Session) GuildBanDelete(guildId, userId string) (err error) {

	_, err = s.Request("DELETE", GUILD_BAN(guildId, userId), ``)
	return
}

// GuildMembers returns an array of Member structures for all members of a
// given guild.
// guildId   : The ID of a Guild.
func (s *Session) GuildMembers(guildId string) (st []Member, err error) {

	body, err := s.Request("GET", GUILD_MEMBERS(guildId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// GuildMemberDelete removes the given user from the given guild.
// guildId   : The ID of a Guild.
// userId    : The ID of a User
func (s *Session) GuildMemberDelete(guildId, userId string) (err error) {

	_, err = s.Request("DELETE", GUILD_MEMBER_DEL(guildId, userId), ``)
	return
}

// GuildChannels returns an array of Channel structures for all channels of a
// given guild.
// guildId   : The ID of a Guild.
func (s *Session) GuildChannels(guildId string) (st []Channel, err error) {

	body, err := s.Request("GET", GUILD_CHANNELS(guildId), ``)
	err = json.Unmarshal(body, &st)

	return
}

// GuildChannelCreate creates a new channel in the given guild
// guildId   : The ID of a Guild.
// name      : Name of the channel (2-100 chars length)
// ctype     : Tpye of the channel (voice or text)
func (s *Session) GuildChannelCreate(guildId, name, ctype string) (st []Channel, err error) {

	body, err := s.Request("POST", GUILD_CHANNELS(guildId), fmt.Sprintf(`{"name":"%s", "type":"%s"}`, name, ctype))
	err = json.Unmarshal(body, &st)
	return
}

// GuildInvites returns an array of Invite structures for the given guild
// guildId   : The ID of a Guild.
func (s *Session) GuildInvites(guildId string) (st []Invite, err error) {
	body, err := s.Request("GET", GUILD_INVITES(guildId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// GuildInviteCreate creates a new invite for the given guild.
// guildId   : The ID of a Guild.
// i         : An Invite struct with the values MaxAge, MaxUses, Temporary,
//             and XkcdPass defined.
func (s *Session) GuildInviteCreate(guildId string, i Invite) (st Invite, err error) {

	payload := fmt.Sprintf(
		`{"max_age":%d, "max_uses":%d, "temporary":%t, "xkcdpass":%t}`,
		i.MaxAge, i.MaxUses, i.Temporary, i.XkcdPass)

	body, err := s.Request("POST", GUILD_INVITES(guildId), payload)
	err = json.Unmarshal(body, &st)
	return
}

/***************************************************************************************************
 * Functions related to a specific channel
 */

// Channel returns a Channel strucutre of a specific Channel.
// channelId  : The ID of the Channel you want returend.
func (s *Session) Channel(channelId string) (st Channel, err error) {
	body, err := s.Request("GET", CHANNEL(channelId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// ChannelEdit edits the given channel
// channelId  : The ID of a Channel
// name       : The new name to assign the channel.
func (s *Session) ChannelEdit(channelId, name string) (st []Channel, err error) {

	body, err := s.Request("PATCH", CHANNEL(channelId), fmt.Sprintf(`{"name":"%s"}`, name))
	err = json.Unmarshal(body, &st)
	return
}

// ChannelDelete deletes the given channel
// channelId  : The ID of a Channel
func (s *Session) ChannelDelete(channelId string) (st []Channel, err error) {

	body, err := s.Request("DELETE", CHANNEL(channelId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// ChannelTyping broadcasts to all members that authenticated user is typing in
// the given channel.
// channelId  : The ID of a Channel
func (s *Session) ChannelTyping(channelId string) (err error) {

	_, err = s.Request("POST", CHANNEL_TYPING(channelId), ``)
	return
}

// ChannelMessages returns an array of Message structures for messaages within
// a given channel.
// channelId : The ID of a Channel.
// limit     : The number messages that can be returned.
// beforeId  : If provided all messages returned will be before given ID.
// afterId   : If provided all messages returned will be after given ID.
func (s *Session) ChannelMessages(channelId string, limit int, beforeId int, afterId int) (st []Message, err error) {

	var urlStr string = ""

	if limit > 0 {
		urlStr = fmt.Sprintf("?limit=%d", limit)
	}

	if afterId > 0 {
		if urlStr != "" {
			urlStr = urlStr + fmt.Sprintf("&after=%d", afterId)
		} else {
			urlStr = fmt.Sprintf("?after=%d", afterId)
		}
	}

	if beforeId > 0 {
		if urlStr != "" {
			urlStr = urlStr + fmt.Sprintf("&before=%d", beforeId)
		} else {
			urlStr = fmt.Sprintf("?before=%d", beforeId)
		}
	}

	body, err := s.Request("GET", CHANNEL_MESSAGES(channelId)+urlStr, ``)
	err = json.Unmarshal(body, &st)
	return
}

// ChannelMessageAck acknowledges and marks the given message as read
// channeld  : The ID of a Channel
// messageId : the ID of a Message
func (s *Session) ChannelMessageAck(channelId, messageId string) (err error) {

	_, err = s.Request("POST", CHANNEL_MESSAGE_ACK(channelId, messageId), ``)
	return
}

// ChannelMessageSend sends a message to the given channel.
// channelId : The ID of a Channel.
// content   : The message to send.
func (s *Session) ChannelMessageSend(channelId string, content string) (st Message, err error) {

	response, err := s.Request("POST", CHANNEL_MESSAGES(channelId), fmt.Sprintf(`{"content":"%s"}`, content))
	err = json.Unmarshal(response, &st)
	return
}

// ChannelMessageEdit edits an existing message, replacing it entirely with
// the given content.
// channeld  : The ID of a Channel
// messageId : the ID of a Message
func (s *Session) ChannelMessageEdit(channelId, messageId, content string) (st Message, err error) {

	response, err := s.Request("PATCH", CHANNEL_MESSAGE(channelId, messageId), fmt.Sprintf(`{"content":"%s"}`, content))
	err = json.Unmarshal(response, &st)
	return
}

// ChannelMessageDelete deletes a message from the Channel.
func (s *Session) ChannelMessageDelete(channelId, messageId string) (err error) {

	_, err = s.Request("DELETE", CHANNEL_MESSAGE(channelId, messageId), ``)
	return
}

// ChannelInvites returns an array of Invite structures for the given channel
// channelId   : The ID of a Channel
func (s *Session) ChannelInvites(channelId string) (st []Invite, err error) {
	body, err := s.Request("GET", CHANNEL_INVITES(channelId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// ChannelInviteCreate creates a new invite for the given channel.
// channelId   : The ID of a Channel
// i           : An Invite struct with the values MaxAge, MaxUses, Temporary,
//               and XkcdPass defined.
func (s *Session) ChannelInviteCreate(channelId string, i Invite) (st Invite, err error) {

	payload := fmt.Sprintf(
		`{"max_age":%d, "max_uses":%d, "temporary":%t, "xkcdpass":%t}`,
		i.MaxAge, i.MaxUses, i.Temporary, i.XkcdPass)

	body, err := s.Request("POST", CHANNEL_INVITES(channelId), payload)
	err = json.Unmarshal(body, &st)
	return
}

/***************************************************************************************************
 * Functions related to an invite
 */

// Inivte returns an Invite structure of the given invite
// inviteId : The invite code (or maybe xkcdpass?)
func (s *Session) Invite(inviteId string) (st Invite, err error) {
	body, err := s.Request("GET", INVITE(inviteId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// InviteDelete deletes an existing invite
// inviteId   : the code (or maybe xkcdpass?) of an invite
func (s *Session) InviteDelete(inviteId string) (st Invite, err error) {

	body, err := s.Request("DELETE", INVITE(inviteId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// InivteAccept accepts an Invite to a Guild or Channel
// inviteId : The invite code (or maybe xkcdpass?)
func (s *Session) InviteAccept(inviteId string) (st Invite, err error) {
	body, err := s.Request("POST", INVITE(inviteId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// https://discordapi.readthedocs.org/en/latest/reference/guilds/invites.html#get-and-accept-invite
func (s *Session) InviteValidate(validateId string) (i Invite, err error) {
	return
}

/***************************************************************************************************
 * Functions related to Voice/Audio
 */

// VoiceRegions returns the voice server regions
func (s *Session) VoiceRegions() (st []VoiceRegion, err error) {

	body, err := s.Request("GET", VOICE_REGIONS, ``)
	err = json.Unmarshal(body, &st)
	return
}

// VoiceIce returns the voice server ICE information
func (s *Session) VoiceIce() (st VoiceIce, err error) {

	body, err := s.Request("GET", VOICE_ICE, ``)
	err = json.Unmarshal(body, &st)
	return
}

/***************************************************************************************************
 * Functions related to Websockets
 */

// Gateway returns the a websocket Gateway address
func (s *Session) Gateway() (gateway string, err error) {

	response, err := s.Request("GET", GATEWAY, ``)

	var temp map[string]interface{}
	err = json.Unmarshal(response, &temp)
	gateway = temp["url"].(string)
	return
}
