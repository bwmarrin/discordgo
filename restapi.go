// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains functions for interacting with the Discord REST/JSON API
// at the lowest level.

package discordgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // For JPEG decoding
	_ "image/png"  // For PNG decoding
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ErrJSONUnmarshal is returned for JSON Unmarshall errors.
var ErrJSONUnmarshal = errors.New("json unmarshal")

// Request makes a (GET/POST/...) Requests to Discord REST API with JSON data.
// All the other Discord REST Calls in this file use this function.
func (s *Session) Request(method, urlStr string, data interface{}) (response []byte, err error) {

	if s.Debug {
		fmt.Println("API REQUEST  PAYLOAD :: [" + fmt.Sprintf("%+v", data) + "]")
	}

	var body []byte
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return
		}
	}

	return s.request(method, urlStr, "application/json", body)
}

// request makes a (GET/POST/...) Requests to Discord REST API.
func (s *Session) request(method, urlStr, contentType string, b []byte) (response []byte, err error) {

	if s.Debug {
		fmt.Printf("API REQUEST %8s :: %s\n", method, urlStr)
	}

	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(b))
	if err != nil {
		return
	}

	// Not used on initial login..
	// TODO: Verify if a login, otherwise complain about no-token
	if s.Token != "" {
		req.Header.Set("authorization", s.Token)
	}

	req.Header.Set("Content-Type", contentType)
	// TODO: Make a configurable static variable.
	req.Header.Set("User-Agent", fmt.Sprintf("DiscordBot (https://github.com/bwmarrin/discordgo, v%s)", VERSION))

	if s.Debug {
		for k, v := range req.Header {
			fmt.Printf("API REQUEST   HEADER :: [%s] = %+v\n", k, v)
		}
	}

	client := &http.Client{Timeout: (20 * time.Second)}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Println("error closing resp body")
		}
	}()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if s.Debug {

		fmt.Printf("API RESPONSE  STATUS :: %s\n", resp.Status)
		for k, v := range resp.Header {
			fmt.Printf("API RESPONSE  HEADER :: [%s] = %+v\n", k, v)
		}
		fmt.Printf("API RESPONSE    BODY :: [%s]\n", response)
	}

	// See http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html
	switch resp.StatusCode {

	case 200: // OK
	case 204: // No Content

		// TODO check for 401 response, invalidate token if we get one.

	case 429: // TOO MANY REQUESTS - Rate limiting
		rl := RateLimit{}
		err = json.Unmarshal(response, &rl)
		if err != nil {
			err = fmt.Errorf("Request unmarshal rate limit error : %+v", err)
			return
		}
		time.Sleep(rl.RetryAfter)
		response, err = s.request(method, urlStr, contentType, b)

	default: // Error condition
		err = fmt.Errorf("HTTP %s, %s", resp.Status, response)
	}

	return
}

func unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return ErrJSONUnmarshal
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Sessions
// ------------------------------------------------------------------------------------------------

// Login asks the Discord server for an authentication token.
func (s *Session) Login(email, password string) (err error) {

	data := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{email, password}

	response, err := s.Request("POST", LOGIN, data)
	if err != nil {
		return
	}

	temp := struct {
		Token string `json:"token"`
	}{}

	err = unmarshal(response, &temp)
	if err != nil {
		return
	}

	s.Token = temp.Token
	return
}

// Register sends a Register request to Discord, and returns the authentication token
// Note that this account is temporary and should be verified for future use.
// Another option is to save the authentication token external, but this isn't recommended.
func (s *Session) Register(username string) (token string, err error) {

	data := struct {
		Username string `json:"username"`
	}{username}

	response, err := s.Request("POST", REGISTER, data)
	if err != nil {
		return
	}

	temp := struct {
		Token string `json:"token"`
	}{}

	err = unmarshal(response, &temp)
	if err != nil {
		return
	}

	token = temp.Token
	return
}

// Logout sends a logout request to Discord.
// This does not seem to actually invalidate the token.  So you can still
// make API calls even after a Logout.  So, it seems almost pointless to
// even use.
func (s *Session) Logout() (err error) {

	//  _, err = s.Request("POST", LOGOUT, fmt.Sprintf(`{"token": "%s"}`, s.Token))

	if s.Token == "" {
		return
	}

	data := struct {
		Token string `json:"token"`
	}{s.Token}

	_, err = s.Request("POST", LOGOUT, data)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Users
// ------------------------------------------------------------------------------------------------

// User returns the user details of the given userID
// userID    : A user ID or "@me" which is a shortcut of current user ID
func (s *Session) User(userID string) (st *User, err error) {

	body, err := s.Request("GET", USER(userID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// UserAvatar returns an image.Image of a users Avatar.
// userID    : A user ID or "@me" which is a shortcut of current user ID
func (s *Session) UserAvatar(userID string) (img image.Image, err error) {
	u, err := s.User(userID)
	if err != nil {
		return
	}

	body, err := s.Request("GET", USER_AVATAR(userID, u.Avatar), nil)
	if err != nil {
		return
	}

	img, _, err = image.Decode(bytes.NewReader(body))
	return
}

// UserUpdate updates a users settings.
func (s *Session) UserUpdate(email, password, username, avatar, newPassword string) (st *User, err error) {

	// NOTE: Avatar must be either the hash/id of existing Avatar or
	// data:image/png;base64,BASE64_STRING_OF_NEW_AVATAR_PNG
	// to set a new avatar.
	// If left blank, avatar will be set to null/blank

	data := struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		Username    string `json:"username"`
		Avatar      string `json:"avatar,omitempty"`
		NewPassword string `json:"new_password,omitempty"`
	}{email, password, username, avatar, newPassword}

	body, err := s.Request("PATCH", USER("@me"), data)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// UserSettings returns the settings for a given user
func (s *Session) UserSettings() (st *Settings, err error) {

	body, err := s.Request("GET", USER_SETTINGS("@me"), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// UserChannels returns an array of Channel structures for all private
// channels.
func (s *Session) UserChannels() (st []*Channel, err error) {

	body, err := s.Request("GET", USER_CHANNELS("@me"), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// UserChannelCreate creates a new User (Private) Channel with another User
// recipientID : A user ID for the user to which this channel is opened with.
func (s *Session) UserChannelCreate(recipientID string) (st *Channel, err error) {

	data := struct {
		RecipientID string `json:"recipient_id"`
	}{recipientID}

	body, err := s.Request("POST", USER_CHANNELS("@me"), data)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// UserGuilds returns an array of Guild structures for all guilds.
func (s *Session) UserGuilds() (st []*Guild, err error) {

	body, err := s.Request("GET", USER_GUILDS("@me"), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Guilds
// ------------------------------------------------------------------------------------------------

// Guild returns a Guild structure of a specific Guild.
// guildID   : The ID of a Guild
func (s *Session) Guild(guildID string) (st *Guild, err error) {
	if s.StateEnabled {
		// Attempt to grab the guild from State first.
		st, err = s.State.Guild(guildID)
		if err == nil {
			return
		}
	}

	body, err := s.Request("GET", GUILD(guildID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildCreate creates a new Guild
// name      : A name for the Guild (2-100 characters)
func (s *Session) GuildCreate(name string) (st *Guild, err error) {

	data := struct {
		Name string `json:"name"`
	}{name}

	body, err := s.Request("POST", GUILDS, data)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildEdit edits a new Guild
// guildID   : The ID of a Guild
// name      : A name for the Guild (2-100 characters)
func (s *Session) GuildEdit(guildID, name string) (st *Guild, err error) {

	data := struct {
		Name string `json:"name"`
	}{name}

	body, err := s.Request("POST", GUILD(guildID), data)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildDelete deletes a Guild.
// guildID   : The ID of a Guild
func (s *Session) GuildDelete(guildID string) (st *Guild, err error) {

	body, err := s.Request("DELETE", GUILD(guildID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildLeave leaves a Guild.
// guildID   : The ID of a Guild
func (s *Session) GuildLeave(guildID string) (err error) {

	_, err = s.Request("DELETE", USER_GUILD("@me", guildID), nil)
	return
}

// GuildBans returns an array of User structures for all bans of a
// given guild.
// guildID   : The ID of a Guild.
func (s *Session) GuildBans(guildID string) (st []*User, err error) {

	body, err := s.Request("GET", GUILD_BANS(guildID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// GuildBanCreate bans the given user from the given guild.
// guildID   : The ID of a Guild.
// userID    : The ID of a User
// days      : The number of days of previous comments to delete.
func (s *Session) GuildBanCreate(guildID, userID string, days int) (err error) {

	uri := GUILD_BAN(guildID, userID)

	if days > 0 {
		uri = fmt.Sprintf("%s?delete-message-days=%d", uri, days)
	}

	_, err = s.Request("PUT", uri, nil)
	return
}

// GuildBanDelete removes the given user from the guild bans
// guildID   : The ID of a Guild.
// userID    : The ID of a User
func (s *Session) GuildBanDelete(guildID, userID string) (err error) {

	_, err = s.Request("DELETE", GUILD_BAN(guildID, userID), nil)
	return
}

// GuildMembers returns a list of members for a guild.
//  guildID  : The ID of a Guild.
//  offset   : A number of members to skip
//  limit    : max number of members to return
func (s *Session) GuildMembers(guildID string, offset, limit int) (st []*Member, err error) {

	uri := GUILD_MEMBERS(guildID)

	v := url.Values{}

	if offset > 0 {
		v.Set("offset", strconv.Itoa(offset))
	}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}

	if len(v) > 0 {
		uri = fmt.Sprintf("%s?%s", uri, v.Encode())
	}

	body, err := s.Request("GET", uri, nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildMember returns a member of a guild.
//  guildID   : The ID of a Guild.
//  userID    : The ID of a User
func (s *Session) GuildMember(guildID, userID string) (st *Member, err error) {

	body, err := s.Request("GET", GUILD_MEMBER(guildID, userID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildMemberDelete removes the given user from the given guild.
// guildID   : The ID of a Guild.
// userID    : The ID of a User
func (s *Session) GuildMemberDelete(guildID, userID string) (err error) {

	_, err = s.Request("DELETE", GUILD_MEMBER(guildID, userID), nil)
	return
}

// GuildMemberEdit edits the roles of a member.
// guildID  : The ID of a Guild.
// userID   : The ID of a User.
// roles    : A list of role ID's to set on the member.
func (s *Session) GuildMemberEdit(guildID, userID string, roles []string) (err error) {

	data := struct {
		Roles []string `json:"roles"`
	}{roles}

	_, err = s.Request("PATCH", GUILD_MEMBER(guildID, userID), data)
	if err != nil {
		return
	}

	return
}

// GuildMemberMove moves a guild member from one voice channel to another/none
//  guildID   : The ID of a Guild.
//  userID    : The ID of a User.
//  channelID : The ID of a channel to move user to, or null?
// NOTE : I am not entirely set on the name of this function and it may change
// prior to the final 1.0.0 release of Discordgo
func (s *Session) GuildMemberMove(guildID, userID, channelID string) (err error) {

	data := struct {
		ChannelID string `json:"channel_id"`
	}{channelID}

	_, err = s.Request("PATCH", GUILD_MEMBER(guildID, userID), data)
	if err != nil {
		return
	}

	return
}

// GuildChannels returns an array of Channel structures for all channels of a
// given guild.
// guildID   : The ID of a Guild.
func (s *Session) GuildChannels(guildID string) (st []*Channel, err error) {

	body, err := s.Request("GET", GUILD_CHANNELS(guildID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// GuildChannelCreate creates a new channel in the given guild
// guildID   : The ID of a Guild.
// name      : Name of the channel (2-100 chars length)
// ctype     : Tpye of the channel (voice or text)
func (s *Session) GuildChannelCreate(guildID, name, ctype string) (st *Channel, err error) {

	data := struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}{name, ctype}

	body, err := s.Request("POST", GUILD_CHANNELS(guildID), data)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildInvites returns an array of Invite structures for the given guild
// guildID   : The ID of a Guild.
func (s *Session) GuildInvites(guildID string) (st []*Invite, err error) {
	body, err := s.Request("GET", GUILD_INVITES(guildID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildRoles returns all roles for a given guild.
// guildID   : The ID of a Guild.
func (s *Session) GuildRoles(guildID string) (st []*Role, err error) {

	body, err := s.Request("GET", GUILD_ROLES(guildID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return // TODO return pointer
}

// GuildRoleCreate returns a new Guild Role.
// guildID: The ID of a Guild.
func (s *Session) GuildRoleCreate(guildID string) (st *Role, err error) {

	body, err := s.Request("POST", GUILD_ROLES(guildID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// GuildRoleEdit updates an existing Guild Role with new values
// guildID   : The ID of a Guild.
// roleID    : The ID of a Role.
// name      : The name of the Role.
// color     : The color of the role (decimal, not hex).
// hoist     : Whether to display the role's users separately.
// perm      : The permissions for the role.
func (s *Session) GuildRoleEdit(guildID, roleID, name string, color int, hoist bool, perm int) (st *Role, err error) {

	data := struct {
		Name        string `json:"name"`        // The color the role should have (as a decimal, not hex)
		Color       int    `json:"color"`       // Whether to display the role's users separately
		Hoist       bool   `json:"hoist"`       // The role's name (overwrites existing)
		Permissions int    `json:"permissions"` // The overall permissions number of the role (overwrites existing)
	}{name, color, hoist, perm}

	body, err := s.Request("PATCH", GUILD_ROLE(guildID, roleID), data)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// GuildRoleReorder reoders guild roles
// guildID   : The ID of a Guild.
// roles     : A list of ordered roles.
func (s *Session) GuildRoleReorder(guildID string, roles []*Role) (st []*Role, err error) {

	body, err := s.Request("PATCH", GUILD_ROLES(guildID), roles)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// GuildRoleDelete deletes an existing role.
// guildID   : The ID of a Guild.
// roleID    : The ID of a Role.
func (s *Session) GuildRoleDelete(guildID, roleID string) (err error) {

	_, err = s.Request("DELETE", GUILD_ROLE(guildID, roleID), nil)

	return
}

// GuildIcon returns an image.Image of a guild icon.
// guildID   : The ID of a Guild.
func (s *Session) GuildIcon(guildID string) (img image.Image, err error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return
	}

	if g.Icon == "" {
		err = errors.New("Guild does not have an icon set.")
		return
	}

	body, err := s.Request("GET", GUILD_ICON(guildID, g.Icon), nil)
	if err != nil {
		return
	}

	img, _, err = image.Decode(bytes.NewReader(body))
	return
}

// GuildSplash returns an image.Image of a guild splash image.
// guildID   : The ID of a Guild.
func (s *Session) GuildSplash(guildID string) (img image.Image, err error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return
	}

	if g.Splash == "" {
		err = errors.New("Guild does not have a splash set.")
		return
	}

	body, err := s.Request("GET", GUILD_SPLASH(guildID, g.Splash), nil)
	if err != nil {
		return
	}

	img, _, err = image.Decode(bytes.NewReader(body))
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Channels
// ------------------------------------------------------------------------------------------------

// Channel returns a Channel strucutre of a specific Channel.
// channelID  : The ID of the Channel you want returned.
func (s *Session) Channel(channelID string) (st *Channel, err error) {
	body, err := s.Request("GET", CHANNEL(channelID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ChannelEdit edits the given channel
// channelID  : The ID of a Channel
// name       : The new name to assign the channel.
func (s *Session) ChannelEdit(channelID, name string) (st *Channel, err error) {

	data := struct {
		Name string `json:"name"`
	}{name}

	body, err := s.Request("PATCH", CHANNEL(channelID), data)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ChannelDelete deletes the given channel
// channelID  : The ID of a Channel
func (s *Session) ChannelDelete(channelID string) (st *Channel, err error) {

	body, err := s.Request("DELETE", CHANNEL(channelID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ChannelTyping broadcasts to all members that authenticated user is typing in
// the given channel.
// channelID  : The ID of a Channel
func (s *Session) ChannelTyping(channelID string) (err error) {

	_, err = s.Request("POST", CHANNEL_TYPING(channelID), nil)
	return
}

// ChannelMessages returns an array of Message structures for messages within
// a given channel.
// channelID : The ID of a Channel.
// limit     : The number messages that can be returned.
// beforeID  : If provided all messages returned will be before given ID.
// afterID   : If provided all messages returned will be after given ID.
func (s *Session) ChannelMessages(channelID string, limit int, beforeID, afterID string) (st []*Message, err error) {

	uri := CHANNEL_MESSAGES(channelID)

	v := url.Values{}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if afterID != "" {
		v.Set("after", afterID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}
	if len(v) > 0 {
		uri = fmt.Sprintf("%s?%s", uri, v.Encode())
	}

	body, err := s.Request("GET", uri, nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ChannelMessageAck acknowledges and marks the given message as read
// channeld  : The ID of a Channel
// messageID : the ID of a Message
func (s *Session) ChannelMessageAck(channelID, messageID string) (err error) {

	_, err = s.Request("POST", CHANNEL_MESSAGE_ACK(channelID, messageID), nil)
	return
}

// channelMessageSend sends a message to the given channel.
// channelID : The ID of a Channel.
// content   : The message to send.
// tts       : Whether to send the message with TTS.
func (s *Session) channelMessageSend(channelID, content string, tts bool) (st *Message, err error) {

	// TODO: nonce string ?
	data := struct {
		Content string `json:"content"`
		TTS     bool   `json:"tts"`
	}{content, tts}

	// Send the message to the given channel
	response, err := s.Request("POST", CHANNEL_MESSAGES(channelID), data)
	if err != nil {
		return
	}

	err = unmarshal(response, &st)
	return
}

// ChannelMessageSend sends a message to the given channel.
// channelID : The ID of a Channel.
// content   : The message to send.
func (s *Session) ChannelMessageSend(channelID string, content string) (st *Message, err error) {

	return s.channelMessageSend(channelID, content, false)
}

// ChannelMessageSendTTS sends a message to the given channel with Text to Speech.
// channelID : The ID of a Channel.
// content   : The message to send.
func (s *Session) ChannelMessageSendTTS(channelID string, content string) (st *Message, err error) {

	return s.channelMessageSend(channelID, content, true)
}

// ChannelMessageEdit edits an existing message, replacing it entirely with
// the given content.
// channeld  : The ID of a Channel
// messageID : the ID of a Message
func (s *Session) ChannelMessageEdit(channelID, messageID, content string) (st *Message, err error) {

	data := struct {
		Content string `json:"content"`
	}{content}

	response, err := s.Request("PATCH", CHANNEL_MESSAGE(channelID, messageID), data)
	if err != nil {
		return
	}

	err = unmarshal(response, &st)
	return
}

// ChannelMessageDelete deletes a message from the Channel.
func (s *Session) ChannelMessageDelete(channelID, messageID string) (err error) {

	_, err = s.Request("DELETE", CHANNEL_MESSAGE(channelID, messageID), nil)
	return
}

// ChannelFileSend sends a file to the given channel.
// channelID : The ID of a Channel.
// io.Reader : A reader for the file contents.
func (s *Session) ChannelFileSend(channelID, name string, r io.Reader) (st *Message, err error) {

	body := &bytes.Buffer{}
	bodywriter := multipart.NewWriter(body)

	writer, err := bodywriter.CreateFormFile("file", name)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(writer, r)
	if err != nil {
		return
	}

	err = bodywriter.Close()
	if err != nil {
		return
	}

	response, err := s.request("POST", CHANNEL_MESSAGES(channelID), bodywriter.FormDataContentType(), body.Bytes())
	if err != nil {
		return
	}

	err = unmarshal(response, &st)
	return
}

// ChannelInvites returns an array of Invite structures for the given channel
// channelID   : The ID of a Channel
func (s *Session) ChannelInvites(channelID string) (st []*Invite, err error) {

	body, err := s.Request("GET", CHANNEL_INVITES(channelID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ChannelInviteCreate creates a new invite for the given channel.
// channelID   : The ID of a Channel
// i           : An Invite struct with the values MaxAge, MaxUses, Temporary,
//               and XkcdPass defined.
func (s *Session) ChannelInviteCreate(channelID string, i Invite) (st *Invite, err error) {

	data := struct {
		MaxAge    int  `json:"max_age"`
		MaxUses   int  `json:"max_uses"`
		Temporary bool `json:"temporary"`
		XKCDPass  bool `json:"xkcdpass"`
	}{i.MaxAge, i.MaxUses, i.Temporary, i.XkcdPass}

	body, err := s.Request("POST", CHANNEL_INVITES(channelID), data)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ChannelPermissionSet creates a Permission Override for the given channel.
// NOTE: This func name may changed.  Using Set instead of Create because
// you can both create a new override or update an override with this function.
func (s *Session) ChannelPermissionSet(channelID, targetID, targetType string, allow, deny int) (err error) {

	data := struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Allow int    `json:"allow"`
		Deny  int    `json:"deny"`
	}{targetID, targetType, allow, deny}

	_, err = s.Request("PUT", CHANNEL_PERMISSION(channelID, targetID), data)
	return
}

// ChannelPermissionDelete deletes a specific permission override for the given channel.
// NOTE: Name of this func may change.
func (s *Session) ChannelPermissionDelete(channelID, targetID string) (err error) {

	_, err = s.Request("DELETE", CHANNEL_PERMISSION(channelID, targetID), nil)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Invites
// ------------------------------------------------------------------------------------------------

// Invite returns an Invite structure of the given invite
// inviteID : The invite code (or maybe xkcdpass?)
func (s *Session) Invite(inviteID string) (st *Invite, err error) {

	body, err := s.Request("GET", INVITE(inviteID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// InviteDelete deletes an existing invite
// inviteID   : the code (or maybe xkcdpass?) of an invite
func (s *Session) InviteDelete(inviteID string) (st *Invite, err error) {

	body, err := s.Request("DELETE", INVITE(inviteID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// InviteAccept accepts an Invite to a Guild or Channel
// inviteID : The invite code (or maybe xkcdpass?)
func (s *Session) InviteAccept(inviteID string) (st *Invite, err error) {

	body, err := s.Request("POST", INVITE(inviteID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Voice
// ------------------------------------------------------------------------------------------------

// VoiceRegions returns the voice server regions
func (s *Session) VoiceRegions() (st []*VoiceRegion, err error) {

	body, err := s.Request("GET", VOICE_REGIONS, nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// VoiceICE returns the voice server ICE information
func (s *Session) VoiceICE() (st *VoiceICE, err error) {

	body, err := s.Request("GET", VOICE_ICE, nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Websockets
// ------------------------------------------------------------------------------------------------

// Gateway returns the a websocket Gateway address
func (s *Session) Gateway() (gateway string, err error) {

	response, err := s.Request("GET", GATEWAY, nil)
	if err != nil {
		return
	}

	temp := struct {
		URL string `json:"url"`
	}{}

	err = unmarshal(response, &temp)
	if err != nil {
		return
	}

	gateway = temp.URL
	return
}
