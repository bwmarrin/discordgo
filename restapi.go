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
	sv "strconv"
	"time"
)

// Constants of all known Discord API Endpoints
// Please let me know if you know of any others.
const (
	DISCORD  = "http://discordapp.com"
	API      = DISCORD + "/api/"
	GUILDS   = API + "guilds/"
	CHANNELS = API + "channels/"
	USERS    = API + "users/"
	GATEWAY  = API + "gateway"

	AUTH            = API + "auth/"
	LOGIN           = API + AUTH + "login"
	LOGOUT          = API + AUTH + "logout"
	VERIFY          = API + AUTH + "verify"
	VERIFY_RESEND   = API + AUTH + "verify/resend"
	FORGOT_PASSWORD = API + AUTH + "forgot"
	RESET_PASSWORD  = API + AUTH + "reset"
	REGISTER        = API + AUTH + "register"

	VOICE   = API + "/voice/"
	REGIONS = API + VOICE + "regions"
	ICE     = API + VOICE + "ice"

	TUTORIAL            = API + "tutorial/"
	TUTORIAL_INDICATORS = TUTORIAL + "indicators"

	INVITE       = API + "invite"
	TRACK        = API + "track"
	SSO          = API + "sso"
	REPORT       = API + "report"
	INTEGRATIONS = API + "integrations"
)

// Almost like the constants above :) Except can't be constants
var (
	USER             = func(userId string) string { return USERS + userId }
	USER_AVATAR      = func(userId, hash string) string { return USERS + userId + "/avatars/" + hash + ".jpg" }
	USER_SETTINGS    = func(userId string) string { return USERS + userId + "/settings" }
	USER_GUILDS      = func(userId string) string { return USERS + userId + "/guilds" }
	USER_CHANNELS    = func(userId string) string { return USERS + userId + "/channels" }
	USER_DEVICES     = func(userId string) string { return USERS + userId + "/devices" }
	USER_CONNECTIONS = func(userId string) string { return USERS + userId + "/connections" }

	GUILD              = func(guildId int) string { return GUILDS + sv.Itoa(guildId) }
	GUILD_CHANNELS     = func(guildId int) string { return GUILDS + sv.Itoa(guildId) + "/channels" }
	GUILD_MEMBERS      = func(guildId int) string { return GUILDS + sv.Itoa(guildId) + "/members" }
	GUILD_INTEGRATIONS = func(guildId int) string { return GUILDS + sv.Itoa(guildId) + "/integrations" }
	GUILD_BANS         = func(guildId int) string { return GUILDS + sv.Itoa(guildId) + "/bans" }
	GUILD_ROLES        = func(guildId int) string { return GUILDS + sv.Itoa(guildId) + "/roles" }
	GUILD_INVITES      = func(guildId int) string { return GUILDS + sv.Itoa(guildId) + "/invites" }
	GUILD_EMBED        = func(guildId int) string { return GUILDS + sv.Itoa(guildId) + "/embed" }
	GUILD_PRUNE        = func(guildId int) string { return GUILDS + sv.Itoa(guildId) + "/prune" }
	GUILD_ICON         = func(guildId int, hash string) string { return GUILDS + sv.Itoa(guildId) + "/icons/" + hash + ".jpg" }

	CHANNEL             = func(channelId int) string { return CHANNELS + sv.Itoa(channelId) }
	CHANNEL_MESSAGES    = func(channelId int) string { return CHANNELS + sv.Itoa(channelId) + "/messages" }
	CHANNEL_PERMISSIONS = func(channelId int) string { return CHANNELS + sv.Itoa(channelId) + "/permissions" }
	CHANNEL_INVITES     = func(channelId int) string { return CHANNELS + sv.Itoa(channelId) + "/invites" }
	CHANNEL_TYPING      = func(channelId int) string { return CHANNELS + sv.Itoa(channelId) + "/typing" }

	INTEGRATIONS_JOIN = func(intId int) string { return API + "integrations/" + sv.Itoa(intId) + "/join" }
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

// Gateway returns the a websocket Gateway address
func (s *Session) Gateway() (gateway string, err error) {

	response, err := s.Request("GET", GATEWAY, ``)

	var temp map[string]interface{}
	err = json.Unmarshal(response, &temp)
	gateway = temp["url"].(string)
	return
}

// VoiceRegions returns the voice server regions
func (s *Session) VoiceRegions() (st []VoiceRegion, err error) {

	body, err := s.Request("GET", REGIONS, ``)
	err = json.Unmarshal(body, &st)
	return
}

// VoiceIce returns the voice server ICE information
func (s *Session) VoiceIce() (st VoiceIce, err error) {

	body, err := s.Request("GET", ICE, ``)
	err = json.Unmarshal(body, &st)
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
// guildId   : The ID of the Guild you want returend.
func (s *Session) Guild(guildId int) (st []Guild, err error) {

	body, err := s.Request("GET", GUILD(guildId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// GuildMembers returns an array of Member structures for all members of a
// given guild.
// guildId   : The ID of a Guild.
func (s *Session) GuildMembers(guildId int) (st []Member, err error) {

	body, err := s.Request("GET", GUILD_MEMBERS(guildId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// GuildChannels returns an array of Channel structures for all channels of a
// given guild.
// guildId   : The ID of a Guild.
func (s *Session) GuildChannels(guildId int) (st []Channel, err error) {

	body, err := s.Request("GET", GUILD_CHANNELS(guildId), ``)
	err = json.Unmarshal(body, &st)

	return
}

/***************************************************************************************************
 * Functions related to a specific channel
 */

// Channel returns a Channel strucutre of a specific Channel.
// channelId  : The ID of the Channel you want returend.
func (s *Session) Channel(channelId int) (st Channel, err error) {
	body, err := s.Request("GET", CHANNEL(channelId), ``)
	err = json.Unmarshal(body, &st)
	return
}

// ChannelMessages returns an array of Message structures for messaages within
// a given channel.
// channelId : The ID of a Channel.
// limit     : The number messages that can be returned.
// beforeId  : If provided all messages returned will be before given ID.
// afterId   : If provided all messages returned will be after given ID.
func (s *Session) ChannelMessages(channelId int, limit int, beforeId int, afterId int) (st []Message, err error) {

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

// ChannelMessageSend sends a message to the given channel.
// channelId : The ID of a Channel.
// content   : The message to send.
func (s *Session) ChannelMessageSend(channelId int, content string) (st Message, err error) {

	response, err := s.Request("POST", CHANNEL_MESSAGES(channelId), fmt.Sprintf(`{"content":"%s"}`, content))
	err = json.Unmarshal(response, &st)
	return
}
