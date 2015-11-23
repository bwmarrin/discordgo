// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains functions for interacting with the Discord REST/JSON API
// at the lowest level.

package discordgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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
		err = fmt.Errorf("StatusCode: %d, %s", resp.StatusCode, string(response))
		return
	}

	if s.Debug {
		printJSON(response)
	}
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Sessions
// ------------------------------------------------------------------------------------------------

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

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Users
// ------------------------------------------------------------------------------------------------

// User returns the user details of the given userID
// userID    : A user ID or "@me" which is a shortcut of current user ID
func (s *Session) User(userID string) (st User, err error) {

	body, err := s.Request("GET", USER(userID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// UserAvatar returns a ?? of a users Avatar
// userID    : A user ID or "@me" which is a shortcut of current user ID
func (s *Session) UserAvatar(userID string) (st User, err error) {

	u, err := s.User(userID)
	_, err = s.Request("GET", USER_AVATAR(userID, u.Avatar), ``)
	// TODO need to figure out how to handle returning a file
	return
}

// UserSettings returns the settings for a given user
// userID    : A user ID or "@me" which is a shortcut of current user ID
// This seems to only return a result for "@me"
func (s *Session) UserSettings(userID string) (st Settings, err error) {

	body, err := s.Request("GET", USER_SETTINGS(userID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// UserChannels returns an array of Channel structures for all private
// channels for a user
// userID    : A user ID or "@me" which is a shortcut of current user ID
func (s *Session) UserChannels(userID string) (st []Channel, err error) {

	body, err := s.Request("GET", USER_CHANNELS(userID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// UserChannelCreate creates a new User (Private) Channel with another User
// userID      : A user ID or "@me" which is a shortcut of current user ID
// recipientID : A user ID for the user to which this channel is opened with.
func (s *Session) UserChannelCreate(userID, recipientID string) (st Channel, err error) {

	body, err := s.Request(
		"POST",
		USER_CHANNELS(userID),
		fmt.Sprintf(`{"recipient_id": "%s"}`, recipientID))

	err = json.Unmarshal(body, &st)
	return
}

// UserGuilds returns an array of Guild structures for all guilds for a given user
// userID    : A user ID or "@me" which is a shortcut of current user ID
func (s *Session) UserGuilds(userID string) (st []Guild, err error) {

	body, err := s.Request("GET", USER_GUILDS(userID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Guilds
// ------------------------------------------------------------------------------------------------

// Guild returns a Guild structure of a specific Guild.
// guildID   : The ID of a Guild
func (s *Session) Guild(guildID string) (st Guild, err error) {

	body, err := s.Request("GET", GUILD(guildID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// GuildCreate creates a new Guild
// name      : A name for the Guild (2-100 characters)
func (s *Session) GuildCreate(name string) (st Guild, err error) {

	body, err := s.Request("POST", GUILDS, fmt.Sprintf(`{"name":"%s"}`, name))
	err = json.Unmarshal(body, &st)
	return
}

// GuildEdit edits a new Guild
// guildID   : The ID of a Guild
// name      : A name for the Guild (2-100 characters)
func (s *Session) GuildEdit(guildID, name string) (st Guild, err error) {

	body, err := s.Request("POST", GUILD(guildID), fmt.Sprintf(`{"name":"%s"}`, name))
	err = json.Unmarshal(body, &st)
	return
}

// GuildDelete deletes or leaves a Guild.
// guildID   : The ID of a Guild
func (s *Session) GuildDelete(guildID string) (st Guild, err error) {

	body, err := s.Request("DELETE", GUILD(guildID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// GuildBans returns an array of User structures for all bans of a
// given guild.
// guildID   : The ID of a Guild.
func (s *Session) GuildBans(guildID string) (st []User, err error) {

	body, err := s.Request("GET", GUILD_BANS(guildID), ``)
	err = json.Unmarshal(body, &st)

	return
}

// GuildBanAdd bans the given user from the given guild.
// guildID   : The ID of a Guild.
// userID    : The ID of a User
func (s *Session) GuildBanAdd(guildID, userID string) (err error) {

	_, err = s.Request("PUT", GUILD_BAN(guildID, userID), ``)
	return
}

// GuildBanDelete removes the given user from the guild bans
// guildID   : The ID of a Guild.
// userID    : The ID of a User
func (s *Session) GuildBanDelete(guildID, userID string) (err error) {

	_, err = s.Request("DELETE", GUILD_BAN(guildID, userID), ``)
	return
}

// GuildMembers returns an array of Member structures for all members of a
// given guild.
// guildID   : The ID of a Guild.
func (s *Session) GuildMembers(guildID string) (st []Member, err error) {

	body, err := s.Request("GET", GUILD_MEMBERS(guildID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// GuildMemberDelete removes the given user from the given guild.
// guildID   : The ID of a Guild.
// userID    : The ID of a User
func (s *Session) GuildMemberDelete(guildID, userID string) (err error) {

	_, err = s.Request("DELETE", GUILD_MEMBER_DEL(guildID, userID), ``)
	return
}

// GuildChannels returns an array of Channel structures for all channels of a
// given guild.
// guildID   : The ID of a Guild.
func (s *Session) GuildChannels(guildID string) (st []Channel, err error) {

	body, err := s.Request("GET", GUILD_CHANNELS(guildID), ``)
	err = json.Unmarshal(body, &st)

	return
}

// GuildChannelCreate creates a new channel in the given guild
// guildID   : The ID of a Guild.
// name      : Name of the channel (2-100 chars length)
// ctype     : Tpye of the channel (voice or text)
func (s *Session) GuildChannelCreate(guildID, name, ctype string) (st Channel, err error) {

	body, err := s.Request("POST", GUILD_CHANNELS(guildID), fmt.Sprintf(`{"name":"%s", "type":"%s"}`, name, ctype))
	err = json.Unmarshal(body, &st)
	return
}

// GuildInvites returns an array of Invite structures for the given guild
// guildID   : The ID of a Guild.
func (s *Session) GuildInvites(guildID string) (st []Invite, err error) {
	body, err := s.Request("GET", GUILD_INVITES(guildID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// GuildInviteCreate creates a new invite for the given guild.
// guildID   : The ID of a Guild.
// i         : An Invite struct with the values MaxAge, MaxUses, Temporary,
//             and XkcdPass defined.
func (s *Session) GuildInviteCreate(guildID string, i Invite) (st Invite, err error) {

	payload := fmt.Sprintf(
		`{"max_age":%d, "max_uses":%d, "temporary":%t, "xkcdpass":%t}`,
		i.MaxAge, i.MaxUses, i.Temporary, i.XkcdPass)

	body, err := s.Request("POST", GUILD_INVITES(guildID), payload)
	err = json.Unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Channels
// ------------------------------------------------------------------------------------------------

// Channel returns a Channel strucutre of a specific Channel.
// channelID  : The ID of the Channel you want returend.
func (s *Session) Channel(channelID string) (st Channel, err error) {
	body, err := s.Request("GET", CHANNEL(channelID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// ChannelEdit edits the given channel
// channelID  : The ID of a Channel
// name       : The new name to assign the channel.
func (s *Session) ChannelEdit(channelID, name string) (st Channel, err error) {

	body, err := s.Request("PATCH", CHANNEL(channelID), fmt.Sprintf(`{"name":"%s"}`, name))
	err = json.Unmarshal(body, &st)
	return
}

// ChannelDelete deletes the given channel
// channelID  : The ID of a Channel
func (s *Session) ChannelDelete(channelID string) (st Channel, err error) {

	body, err := s.Request("DELETE", CHANNEL(channelID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// ChannelTyping broadcasts to all members that authenticated user is typing in
// the given channel.
// channelID  : The ID of a Channel
func (s *Session) ChannelTyping(channelID string) (err error) {

	_, err = s.Request("POST", CHANNEL_TYPING(channelID), ``)
	return
}

// ChannelMessages returns an array of Message structures for messaages within
// a given channel.
// channelID : The ID of a Channel.
// limit     : The number messages that can be returned.
// beforeID  : If provided all messages returned will be before given ID.
// afterID   : If provided all messages returned will be after given ID.
func (s *Session) ChannelMessages(channelID string, limit int, beforeID int, afterID int) (st []Message, err error) {

	var urlStr string

	if limit > 0 {
		urlStr = fmt.Sprintf("?limit=%d", limit)
	}

	if afterID > 0 {
		if urlStr != "" {
			urlStr = urlStr + fmt.Sprintf("&after=%d", afterID)
		} else {
			urlStr = fmt.Sprintf("?after=%d", afterID)
		}
	}

	if beforeID > 0 {
		if urlStr != "" {
			urlStr = urlStr + fmt.Sprintf("&before=%d", beforeID)
		} else {
			urlStr = fmt.Sprintf("?before=%d", beforeID)
		}
	}

	body, err := s.Request("GET", CHANNEL_MESSAGES(channelID)+urlStr, ``)
	err = json.Unmarshal(body, &st)
	return
}

// ChannelMessageAck acknowledges and marks the given message as read
// channeld  : The ID of a Channel
// messageID : the ID of a Message
func (s *Session) ChannelMessageAck(channelID, messageID string) (err error) {

	_, err = s.Request("POST", CHANNEL_MESSAGE_ACK(channelID, messageID), ``)
	return
}

// ChannelMessageSend sends a message to the given channel.
// channelID : The ID of a Channel.
// content   : The message to send.
func (s *Session) ChannelMessageSend(channelID string, content string) (st Message, err error) {

	response, err := s.Request("POST", CHANNEL_MESSAGES(channelID), fmt.Sprintf(`{"content":"%s"}`, content))
	err = json.Unmarshal(response, &st)
	return
}

// ChannelMessageEdit edits an existing message, replacing it entirely with
// the given content.
// channeld  : The ID of a Channel
// messageID : the ID of a Message
func (s *Session) ChannelMessageEdit(channelID, messageID, content string) (st Message, err error) {

	response, err := s.Request("PATCH", CHANNEL_MESSAGE(channelID, messageID), fmt.Sprintf(`{"content":"%s"}`, content))
	err = json.Unmarshal(response, &st)
	return
}

// ChannelMessageDelete deletes a message from the Channel.
func (s *Session) ChannelMessageDelete(channelID, messageID string) (err error) {

	_, err = s.Request("DELETE", CHANNEL_MESSAGE(channelID, messageID), ``)
	return
}

// ChannelInvites returns an array of Invite structures for the given channel
// channelID   : The ID of a Channel
func (s *Session) ChannelInvites(channelID string) (st []Invite, err error) {
	body, err := s.Request("GET", CHANNEL_INVITES(channelID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// ChannelInviteCreate creates a new invite for the given channel.
// channelID   : The ID of a Channel
// i           : An Invite struct with the values MaxAge, MaxUses, Temporary,
//               and XkcdPass defined.
func (s *Session) ChannelInviteCreate(channelID string, i Invite) (st Invite, err error) {

	payload := fmt.Sprintf(
		`{"max_age":%d, "max_uses":%d, "temporary":%t, "xkcdpass":%t}`,
		i.MaxAge, i.MaxUses, i.Temporary, i.XkcdPass)

	body, err := s.Request("POST", CHANNEL_INVITES(channelID), payload)
	err = json.Unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Invites
// ------------------------------------------------------------------------------------------------

// Invite returns an Invite structure of the given invite
// inviteID : The invite code (or maybe xkcdpass?)
func (s *Session) Invite(inviteID string) (st Invite, err error) {
	body, err := s.Request("GET", INVITE(inviteID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// InviteDelete deletes an existing invite
// inviteID   : the code (or maybe xkcdpass?) of an invite
func (s *Session) InviteDelete(inviteID string) (st Invite, err error) {

	body, err := s.Request("DELETE", INVITE(inviteID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// InviteAccept accepts an Invite to a Guild or Channel
// inviteID : The invite code (or maybe xkcdpass?)
func (s *Session) InviteAccept(inviteID string) (st Invite, err error) {
	body, err := s.Request("POST", INVITE(inviteID), ``)
	err = json.Unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Voice
// ------------------------------------------------------------------------------------------------

// VoiceRegions returns the voice server regions
func (s *Session) VoiceRegions() (st []VoiceRegion, err error) {

	body, err := s.Request("GET", VOICE_REGIONS, ``)
	err = json.Unmarshal(body, &st)
	return
}

// VoiceICE returns the voice server ICE information
func (s *Session) VoiceICE() (st VoiceICE, err error) {

	body, err := s.Request("GET", VOICE_ICE, ``)
	err = json.Unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Websockets
// ------------------------------------------------------------------------------------------------

// Gateway returns the a websocket Gateway address
func (s *Session) Gateway() (gateway string, err error) {

	response, err := s.Request("GET", GATEWAY, ``)

	var temp map[string]interface{}
	err = json.Unmarshal(response, &temp)
	gateway = temp["url"].(string)
	return
}
