// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains custom types, currently only a timestamp wrapper.

package discordgo

import (
	"encoding/json"
	"image"
	"io"
	"net/http"
	"time"
)

// Timestamp stores a timestamp, as sent by the Discord API.
type Timestamp string

// Parse parses a timestamp string into a time.Time object.
// The only time this can fail is if Discord changes their timestamp format.
func (t Timestamp) Parse() (time.Time, error) {
	return time.Parse(time.RFC3339, string(t))
}

// RESTError stores error information about a request with a bad response code.
// Message is not always present, there are cases where api calls can fail
// without returning a json message.
type RESTError struct {
	Request      *http.Request
	Response     *http.Response
	ResponseBody []byte

	Message *APIErrorMessage // Message may be nil.
}

func newRestError(req *http.Request, resp *http.Response, body []byte) *RESTError {
	restErr := &RESTError{
		Request:      req,
		Response:     resp,
		ResponseBody: body,
	}

	// Attempt to decode the error and assume no message was provided if it fails
	var msg *APIErrorMessage
	err := json.Unmarshal(body, &msg)
	if err == nil {
		restErr.Message = msg
	}

	return restErr
}

func (r RESTError) Error() string {
	return "HTTP " + r.Response.Status + ", " + string(r.ResponseBody)
}

// Applicationer is the interface type to describe *Session functionality for
// managing applications.
type Applicationer interface {
	// Application returns an Application structure of a specific Application
	//   appID : The ID of an Application
	Application(appID string) (st *Application, err error)

	// ApplicationBotCreate creates an Application Bot Account
	//
	//   appID : The ID of an Application
	ApplicationBotCreate(appID string) (st *User, err error)

	// ApplicationCreate creates a new Application
	//    name : Name of Application / Bot
	//    uris : Redirect URIs (Not required)
	ApplicationCreate(ap *Application) (st *Application, err error)

	// ApplicationDelete deletes an existing Application
	//   appID : The ID of an Application
	ApplicationDelete(appID string) (err error)

	// ApplicationUpdate updates an existing Application
	//   var : desc
	ApplicationUpdate(appID string, ap *Application) (st *Application, err error)

	// Applications returns all applications for the authenticated user
	Applications() (st []*Application, err error)
}

// Channeler is the interface type describing *Session functionality for
// managing channels.
type Channeler interface {
	// Channel returns a Channel structure of a specific Channel.
	//   channelID  : The ID of the Channel you want returned.
	Channel(channelID string) (st *Channel, err error)

	// ChannelDelete deletes the given channel
	//   channelID  : The ID of a Channel
	ChannelDelete(channelID string) (st *Channel, err error)

	// ChannelEdit edits the given channel
	//   channelID  : The ID of a Channel
	//   name       : The new name to assign the channel.
	ChannelEdit(channelID, name string) (*Channel, error)

	// ChannelEditComplex edits an existing channel, replacing the parameters entirely with ChannelEdit struct
	//   channelID  : The ID of a Channel
	//   data          : The channel struct to send
	ChannelEditComplex(channelID string, data *ChannelEdit) (st *Channel, err error)

	// ChannelInviteCreate creates a new invite for the given channel.
	//   channelID   : The ID of a Channel
	//   i           : An Invite struct with the values MaxAge, MaxUses and Temporary defined.
	ChannelInviteCreate(channelID string, i Invite) (st *Invite, err error)

	// ChannelInvites returns an array of Invite structures for the given channel
	//   channelID   : The ID of a Channel
	ChannelInvites(channelID string) (st []*Invite, err error)

	// ChannelPermissionDelete deletes a specific permission override for the given channel.
	// NOTE: Name of this func may change.
	ChannelPermissionDelete(channelID, targetID string) (err error)

	// ChannelPermissionSet creates a Permission Override for the given channel.
	// NOTE: This func name may changed.  Using Set instead of Create because
	// you can both create a new override or update an override with this function.
	ChannelPermissionSet(channelID, targetID, targetType string, allow, deny int) (err error)

	// GuildChannelCreate creates a new channel in the given guild
	//   guildID   : The ID of a Guild.
	//   name      : Name of the channel (2-100 chars length)
	//   ctype     : Type of the channel
	GuildChannelCreate(guildID, name string, ctype ChannelType) (st *Channel, err error)

	// GuildChannelCreateComplex creates a new channel in the given guild
	//   guildID      : The ID of a Guild
	//   data         : A data struct describing the new Channel, Name and Type are mandatory, other fields depending on the type
	GuildChannelCreateComplex(guildID string, data GuildChannelCreateData) (st *Channel, err error)

	// GuildChannels returns an array of Channel structures for all channels of a
	// given guild.
	//   guildID   : The ID of a Guild.
	GuildChannels(guildID string) (st []*Channel, err error)

	// GuildChannelsReorder updates the order of channels in a guild
	//   guildID   : The ID of a Guild.
	//   channels  : Updated channels.
	GuildChannelsReorder(guildID string, channels []*Channel) (err error)
}

// ChannelFileSender is the interface type describing *Session functionality for
// sending files.
type ChannelFileSender interface {
	// ChannelFileSend sends a file to the given channel.
	//   channelID : The ID of a Channel.
	//   name: The name of the file.
	//   io.Reader : A reader for the file contents.
	ChannelFileSend(channelID, name string, r io.Reader) (*Message, error)

	// ChannelFileSendWithMessage sends a file to the given channel with an message.
	// DEPRECATED. Use ChannelMessageSendComplex instead.
	//   channelID : The ID of a Channel.
	//   content: Optional Message content.
	//   name: The name of the file.
	//   io.Reader : A reader for the file contents.
	ChannelFileSendWithMessage(channelID, content string, name string, r io.Reader) (*Message, error)
}

// ChannelMessager is the interface type to describe *Session functionality for
// interacting with messages to channels.
type ChannelMessager interface {
	// ChannelMessage gets a single message by ID from a given channel.
	//   channeld  : The ID of a Channel
	//   messageID : the ID of a Message
	ChannelMessage(channelID, messageID string) (st *Message, err error)

	// ChannelMessageAck acknowledges and marks the given message as read
	//   channeld  : The ID of a Channel
	//   messageID : the ID of a Message
	//   lastToken : token returned by last ack
	ChannelMessageAck(channelID, messageID, lastToken string) (st *Ack, err error)

	// ChannelMessageDelete deletes a message from the Channel.
	ChannelMessageDelete(channelID, messageID string) (err error)

	// ChannelMessageEdit edits an existing message, replacing it entirely with
	// the given content.
	//   channelID  : The ID of a Channel
	//   messageID  : The ID of a Message
	//   content    : The contents of the message
	ChannelMessageEdit(channelID, messageID, content string) (*Message, error)

	// ChannelMessageEditComplex edits an existing message, replacing it entirely with
	// the given MessageEdit struct
	ChannelMessageEditComplex(m *MessageEdit) (st *Message, err error)

	// ChannelMessageEditEmbed edits an existing message with embedded data.
	//   channelID : The ID of a Channel
	//   messageID : The ID of a Message
	//   embed     : The embed data to send
	ChannelMessageEditEmbed(channelID, messageID string, embed *MessageEmbed) (*Message, error)

	// ChannelMessagePin pins a message within a given channel.
	//   channelID: The ID of a channel.
	//   messageID: The ID of a message.
	ChannelMessagePin(channelID, messageID string) (err error)

	// ChannelMessageSend sends a message to the given channel.
	//   channelID : The ID of a Channel.
	//   content   : The message to send.
	ChannelMessageSend(channelID string, content string) (*Message, error)

	// ChannelMessageSendComplex sends a message to the given channel.
	//   channelID : The ID of a Channel.
	//   data      : The message struct to send.
	ChannelMessageSendComplex(channelID string, data *MessageSend) (st *Message, err error)

	// ChannelMessageSendEmbed sends a message to the given channel with embedded data.
	//   channelID : The ID of a Channel.
	//   embed     : The embed data to send.
	ChannelMessageSendEmbed(channelID string, embed *MessageEmbed) (*Message, error)

	// ChannelMessageSendTTS sends a message to the given channel with Text to Speech.
	//   channelID : The ID of a Channel.
	//   content   : The message to send.
	ChannelMessageSendTTS(channelID string, content string) (*Message, error)

	// ChannelMessageUnpin unpins a message within a given channel.
	//   channelID: The ID of a channel.
	//   messageID: The ID of a message.
	ChannelMessageUnpin(channelID, messageID string) (err error)

	// ChannelMessages returns an array of Message structures for messages within
	// a given channel.
	//   channelID : The ID of a Channel.
	//   limit     : The number messages that can be returned. (max 100)
	//   beforeID  : If provided all messages returned will be before given ID.
	//   afterID   : If provided all messages returned will be after given ID.
	//   aroundID  : If provided all messages returned will be around given ID.
	ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) (st []*Message, err error)

	// ChannelMessagesBulkDelete bulk deletes the messages from the channel for the provided messageIDs.
	// If only one messageID is in the slice call channelMessageDelete function.
	// If the slice is empty do nothing.
	//   channelID : The ID of the channel for the messages to delete.
	//   messages  : The IDs of the messages to be deleted. A slice of string IDs. A maximum of 100 messages.
	ChannelMessagesBulkDelete(channelID string, messages []string) (err error)

	// ChannelMessagesPinned returns an array of Message structures for pinned messages
	// within a given channel
	//   channelID : The ID of a Channel.
	ChannelMessagesPinned(channelID string) (st []*Message, err error)

	// ChannelTyping broadcasts to all members that authenticated user is typing in
	// the given channel.
	//   channelID  : The ID of a Channel
	ChannelTyping(channelID string) (err error)
}

// Guilder is the interface type to describe *Session functionality for managing
// guilds and self membership.
type Guilder interface {
	// Guild returns a Guild structure of a specific Guild.
	//   guildID   : The ID of a Guild
	Guild(guildID string) (st *Guild, err error)

	// GuildCreate creates a new Guild
	//   name      : A name for the Guild (2-100 characters)
	GuildCreate(name string) (st *Guild, err error)

	// GuildDelete deletes a Guild.
	//   guildID   : The ID of a Guild
	GuildDelete(guildID string) (st *Guild, err error)

	// GuildEdit edits a new Guild
	//   guildID   : The ID of a Guild
	//   g 		 : A GuildParams struct with the values Name, Region and VerificationLevel defined.
	GuildEdit(guildID string, g GuildParams) (st *Guild, err error)

	// GuildIcon returns an image.Image of a guild icon.
	//   guildID   : The ID of a Guild.
	GuildIcon(guildID string) (img image.Image, err error)

	// GuildLeave leaves a Guild.
	//   guildID   : The ID of a Guild
	GuildLeave(guildID string) (err error)

	// GuildSplash returns an image.Image of a guild splash image.
	//   guildID   : The ID of a Guild.
	GuildSplash(guildID string) (img image.Image, err error)
}

// GuildAuditLoger is the interface type to describe *Session functionality for
// accessing the guild's audit log.
type GuildAuditLoger interface {
	// GuildAuditLog returns the audit log for a Guild.
	//   guildID     : The ID of a Guild.
	//   userID      : If provided the log will be filtered for the given ID.
	//   beforeID    : If provided all log entries returned will be before the given ID.
	//   actionType  : If provided the log will be filtered for the given Action Type.
	//   limit       : The number messages that can be returned. (default 50, min 1, max 100)
	GuildAuditLog(guildID, userID, beforeID string, actionType, limit int) (st *GuildAuditLog, err error)
}

// GuildBanner is the interface type to describe *Session functionality for
// managing the guild's bans.
type GuildBanner interface {
	GuildBanCreate(guildID, userID string, days int) (err error)
	GuildBanCreateWithReason(guildID, userID, reason string, days int) (err error)
	GuildBanDelete(guildID, userID string) (err error)
	GuildBans(guildID string) (st []*GuildBan, err error)
}

// GuildEmbeder is the intereface type to describe *Session functionality for
// managing the guild's embeds.
type GuildEmbeder interface {
	// GuildEmbed returns the embed for a Guild.
	//   guildID   : The ID of a Guild.
	GuildEmbed(guildID string) (st *GuildEmbed, err error)

	// GuildEmbed returns the embed for a Guild.
	//   guildID   : The ID of a Guild.
	GuildEmbedEdit(guildID string, enabled bool, channelID string) (err error)
}

// GuildEmojier is the interface type to describe *Session functionality for
// managing the guild's emoji.
type GuildEmojier interface {
	// GuildEmojiCreate creates a new emoji
	//   guildID : The ID of a Guild.
	//   name    : The Name of the Emoji.
	//   image   : The base64 encoded emoji image, has to be smaller than 256KB.
	//   roles   : The roles for which this emoji will be whitelisted, can be nil.
	GuildEmojiCreate(guildID, name, image string, roles []string) (emoji *Emoji, err error)

	// GuildEmojiDelete deletes an Emoji.
	//   guildID : The ID of a Guild.
	//   emojiID : The ID of an Emoji.
	GuildEmojiDelete(guildID, emojiID string) (err error)

	// GuildEmojiEdit modifies an emoji
	//   guildID : The ID of a Guild.
	//   emojiID : The ID of an Emoji.
	//   name    : The Name of the Emoji.
	//   roles   : The roles for which this emoji will be whitelisted, can be nil.
	GuildEmojiEdit(guildID, emojiID, name string, roles []string) (emoji *Emoji, err error)
}

// GuildIntegrationer is the interface type to describe *Session functionality
// for managing the guild's integrations.
type GuildIntegrationer interface {
	// GuildIntegrationCreate creates a Guild Integration.
	//   guildID          : The ID of a Guild.
	//   integrationType  : The Integration type.
	//   integrationID    : The ID of an integration.
	GuildIntegrationCreate(guildID, integrationType, integrationID string) (err error)

	// GuildIntegrationDelete removes the given integration from the Guild.
	//   guildID          : The ID of a Guild.
	//   integrationID    : The ID of an integration.
	GuildIntegrationDelete(guildID, integrationID string) (err error)

	// GuildIntegrationEdit edits a Guild Integration.
	//   guildID              : The ID of a Guild.
	//   integrationType      : The Integration type.
	//   integrationID        : The ID of an integration.
	//   expireBehavior	      : The behavior when an integration subscription lapses (see the integration object documentation).
	//   expireGracePeriod    : Period (in seconds) where the integration will ignore lapsed subscriptions.
	//   enableEmoticons	    : Whether emoticons should be synced for this integration (twitch only currently).
	GuildIntegrationEdit(guildID, integrationID string, expireBehavior, expireGracePeriod int, enableEmoticons bool) (err error)

	// GuildIntegrationSync syncs an integration.
	//   guildID          : The ID of a Guild.
	//   integrationID    : The ID of an integration.
	GuildIntegrationSync(guildID, integrationID string) (err error)

	// GuildIntegrations returns an array of Integrations for a guild.
	//   guildID   : The ID of a Guild.
	GuildIntegrations(guildID string) (st []*Integration, err error)
}

// GuildMemberer is the interface type to describe *Session functionality for
// managing the guild's integrations.
type GuildMemberer interface {
	// GuildMember returns a member of a guild.
	//   guildID   : The ID of a Guild.
	//   userID    : The ID of a User
	GuildMember(guildID, userID string) (st *Member, err error)

	// GuildMemberAdd force joins a user to the guild.
	//   accessToken   : Valid access_token for the user.
	//   guildID       : The ID of a Guild.
	//   userID        : The ID of a User.
	//   nick          : Value to set users nickname to
	//   roles         : A list of role ID's to set on the member.
	//   mute          : If the user is muted.
	//   deaf          : If the user is deafened.
	GuildMemberAdd(accessToken, guildID, userID, nick string, roles []string, mute, deaf bool) (err error)

	// GuildMemberDelete removes the given user from the given guild.
	//   guildID   : The ID of a Guild.
	//   userID    : The ID of a User
	GuildMemberDelete(guildID, userID string) (err error)

	// GuildMemberDeleteWithReason removes the given user from the given guild.
	//   guildID   : The ID of a Guild.
	//   userID    : The ID of a User
	//   reason    : The reason for the kick
	GuildMemberDeleteWithReason(guildID, userID, reason string) (err error)

	// GuildMemberEdit edits the roles of a member.
	//   guildID  : The ID of a Guild.
	//   userID   : The ID of a User.
	//   roles    : A list of role ID's to set on the member.
	GuildMemberEdit(guildID, userID string, roles []string) (err error)

	// GuildMemberMove moves a guild member from one voice channel to another/none
	//   guildID   : The ID of a Guild.
	//   userID    : The ID of a User.
	//   channelID : The ID of a channel to move user to, or null?
	// NOTE : I am not entirely set on the name of this function and it may change
	// prior to the final 1.0.0 release of Discordgo
	GuildMemberMove(guildID, userID, channelID string) (err error)

	// GuildMemberNickname updates the nickname of a guild member
	//   guildID   : The ID of a guild
	//   userID    : The ID of a user
	//   userID    : The ID of a user or "@me" which is a shortcut of the current user ID
	GuildMemberNickname(guildID, userID, nickname string) (err error)

	// GuildMemberRoleAdd adds the specified role to a given member
	//   guildID   : The ID of a Guild.
	//   userID    : The ID of a User.
	//   roleID 	  : The ID of a Role to be assigned to the user.
	GuildMemberRoleAdd(guildID, userID, roleID string) (err error)

	// GuildMemberRoleRemove removes the specified role to a given member
	//   guildID   : The ID of a Guild.
	//   userID    : The ID of a User.
	//   roleID 	  : The ID of a Role to be removed from the user.
	GuildMemberRoleRemove(guildID, userID, roleID string) (err error)

	// GuildMembers returns a list of members for a guild.
	//    guildID  : The ID of a Guild.
	//    after    : The id of the member to return members after
	//    limit    : max number of members to return (max 1000)
	GuildMembers(guildID string, after string, limit int) (st []*Member, err error)

	// GuildPrune Begin as prune operation. Requires the 'KICK_MEMBERS' permission.
	// Returns an object with one 'pruned' key indicating the number of members that were removed in the prune operation.
	//   guildID	: The ID of a Guild.
	//   days		: The number of days to count prune for (1 or more).
	GuildPrune(guildID string, days uint32) (count uint32, err error)

	// GuildPruneCount Returns the number of members that would be removed in a prune operation.
	// Requires 'KICK_MEMBER' permission.
	//   guildID	: The ID of a Guild.
	//   days		: The number of days to count prune for (1 or more).
	GuildPruneCount(guildID string, days uint32) (count uint32, err error)

	// RequestGuildMembers requests guild members from the gateway
	// The gateway responds with GuildMembersChunk events
	//   guildID  : The ID of the guild to request members of
	//   query    : String that username starts with, leave empty to return all members
	//   limit    : Max number of items to return, or 0 to request all members matched
	RequestGuildMembers(guildID, query string, limit int) (err error)
}

// GuildRoler is the interface type to describe *Session functionality for
// managing the guild's roles.
type GuildRoler interface {
	// GuildRoleCreate returns a new Guild Role.
	//   guildID: The ID of a Guild.
	GuildRoleCreate(guildID string) (st *Role, err error)

	// GuildRoleDelete deletes an existing role.
	//   guildID   : The ID of a Guild.
	//   roleID    : The ID of a Role.
	GuildRoleDelete(guildID, roleID string) (err error)

	// GuildRoleEdit updates an existing Guild Role with new values
	//   guildID   : The ID of a Guild.
	//   roleID    : The ID of a Role.
	//   name      : The name of the Role.
	//   color     : The color of the role (decimal, not hex).
	//   hoist     : Whether to display the role's users separately.
	//   perm      : The permissions for the role.
	//   mention   : Whether this role is mentionable
	GuildRoleEdit(guildID, roleID, name string, color int, hoist bool, perm int, mention bool) (st *Role, err error)

	// GuildRoleReorder reoders guild roles
	//   guildID   : The ID of a Guild.
	//   roles     : A list of ordered roles.
	GuildRoleReorder(guildID string, roles []*Role) (st []*Role, err error)

	// GuildRoles returns all roles for a given guild.
	//   guildID   : The ID of a Guild.
	GuildRoles(guildID string) (st []*Role, err error)
}

// Inviter is the interface type to describe *Session functionality for managing
// the guild's invites.
type Inviter interface {
	// GuildInvites returns an array of Invite structures for the given guild
	//   guildID   : The ID of a Guild.
	GuildInvites(guildID string) (st []*Invite, err error)

	// Invite returns an Invite structure of the given invite
	//   inviteID : The invite code
	Invite(inviteID string) (st *Invite, err error)

	// InviteAccept accepts an Invite to a Guild or Channel
	//   inviteID : The invite code
	InviteAccept(inviteID string) (st *Invite, err error)

	// InviteDelete deletes an existing invite
	//   inviteID   : the code of an invite
	InviteDelete(inviteID string) (st *Invite, err error)

	// InviteWithCounts returns an Invite structure of the given invite including approximate member counts
	//   inviteID : The invite code
	InviteWithCounts(inviteID string) (st *Invite, err error)
}

// MessageReactioner is the interface type to describe *Session functionality
// for message reactions.
type MessageReactioner interface {
	// MessageReactionAdd creates an emoji reaction to a message.
	//   channelID : The channel ID.
	//   messageID : The message ID.
	//   emojiID   : Either the unicode emoji for the reaction, or a guild emoji identifier.
	MessageReactionAdd(channelID, messageID, emojiID string) error

	// MessageReactionRemove deletes an emoji reaction to a message.
	//   channelID : The channel ID.
	//   messageID : The message ID.
	//   emojiID   : Either the unicode emoji for the reaction, or a guild emoji identifier.
	//   userID	 : @me or ID of the user to delete the reaction for.
	MessageReactionRemove(channelID, messageID, emojiID, userID string) error

	// MessageReactions gets all the users reactions for a specific emoji.
	//   channelID : The channel ID.
	//   messageID : The message ID.
	//   emojiID   : Either the unicode emoji for the reaction, or a guild emoji identifier.
	//   limit    : max number of users to return (max 100)
	MessageReactions(channelID, messageID, emojiID string, limit int) (st []*User, err error)

	// MessageReactionsRemoveAll deletes all reactions from a message
	//   channelID : The channel ID
	//   messageID : The message ID.
	MessageReactionsRemoveAll(channelID, messageID string) error
}

// Relationshiper is the interface type to describe *Session functionality for
// managing relationships.
type Relationshiper interface {
	// RelationshipDelete removes the relationship with a user.
	//   userID: ID of the user.
	RelationshipDelete(userID string) (err error)

	// RelationshipFriendRequestAccept accepts a friend request from a user.
	//   userID: ID of the user.
	RelationshipFriendRequestAccept(userID string) (err error)

	// RelationshipFriendRequestSend sends a friend request to a user.
	//   userID: ID of the user.
	RelationshipFriendRequestSend(userID string) (err error)

	// RelationshipUserBlock blocks a user.
	//   userID: ID of the user.
	RelationshipUserBlock(userID string) (err error)

	// RelationshipsGet returns an array of all the relationships of the user.
	RelationshipsGet() (r []*Relationship, err error)

	// RelationshipsMutualGet returns an array of all the users both @me and the given user is friends with.
	//   userID: ID of the user.
	RelationshipsMutualGet(userID string) (mf []*User, err error)
}

// Requester is the interface type to describe *Session functionality around
// low-level requests.
type Requester interface {
	// Request is the same as RequestWithBucketID but the bucket id is the same
	// as the urlStr
	Request(method, urlStr string, data interface{}) (response []byte, err error)

	// RequestWithBucketID makes a (GET/POST/...) Requests to Discord REST API
	// with JSON data.
	RequestWithBucketID(method, urlStr string, data interface{}, bucketID string) (response []byte, err error)

	// RequestWithLockedBucket makes a request using a bucket that's already
	// been locked
	RequestWithLockedBucket(method, urlStr, contentType string, b []byte, bucket *Bucket, sequence int) (response []byte, err error)
}

// StatusUpdater is the interface type to describe *Session functionality for
// managing user status updates.
type StatusUpdater interface {
	// UpdateListeningStatus is used to set the user to "Listening to..."
	// If game!="" then set to what user is listening to
	// Else, set user to active and no game.
	UpdateListeningStatus(game string) (err error)

	// UpdateStatus is used to update the user's status.
	// If idle>0 then set status to idle.
	// If game!="" then set game.
	// if otherwise, set status to active, and no game.
	UpdateStatus(idle int, game string) (err error)

	// UpdateStatusComplex allows for sending the raw status update data
	// untouched by discordgo.
	UpdateStatusComplex(usd UpdateStatusData) (err error)

	// UpdateStreamingStatus is used to update the user's streaming status.
	// If idle>0 then set status to idle.
	// If game!="" then set game.
	// If game!="" and url!="" then set the status type to streaming with the URL set.
	// if otherwise, set status to active, and no game.
	UpdateStreamingStatus(idle int, game string, url string) (err error)

	// UserUpdateStatus update the user status
	//   status   : The new status (Actual valid status are 'online','idle','dnd','invisible')
	UserUpdateStatus(status Status) (st *Settings, err error)
}

// Userer is the interface type to describe *Session functionlaity for doing
// things as the authenticated user or for interacting with other users.
type Userer interface {
	// User returns the user details of the given userID
	//   userID    : A user ID or "@me" which is a shortcut of current user ID
	User(userID string) (st *User, err error)

	// UserAvatar is deprecated. Please use UserAvatarDecode
	//   userID    : A user ID or "@me" which is a shortcut of current user ID
	UserAvatar(userID string) (img image.Image, err error)

	// UserAvatarDecode returns an image.Image of a user's Avatar
	//   user : The user which avatar should be retrieved
	UserAvatarDecode(u *User) (img image.Image, err error)

	// UserChannelCreate creates a new User (Private) Channel with another User
	//   recipientID : A user ID for the user to which this channel is opened with.
	UserChannelCreate(recipientID string) (st *Channel, err error)

	// UserChannelPermissions returns the permission of a user in a channel.
	//   userID    : The ID of the user to calculate permissions for.
	//   channelID : The ID of the channel to calculate permission for.
	//
	// NOTE: This function is now deprecated and will be removed in the future.
	// Please see the same function inside state.go
	UserChannelPermissions(userID, channelID string) (apermissions int, err error)

	// UserChannels returns an array of Channel structures for all private
	// channels.
	UserChannels() (st []*Channel, err error)

	// UserConnections returns the user's connections
	UserConnections() (conn []*UserConnection, err error)

	// UserGuildSettingsEdit Edits the users notification settings for a guild
	//   guildID   : The ID of the guild to edit the settings on
	//   settings  : The settings to update
	UserGuildSettingsEdit(guildID string, settings *UserGuildSettingsEdit) (st *UserGuildSettings, err error)

	// UserGuilds returns an array of UserGuild structures for all guilds.
	//   limit     : The number guilds that can be returned. (max 100)
	//   beforeID  : If provided all guilds returned will be before given ID.
	//   afterID   : If provided all guilds returned will be after given ID.
	UserGuilds(limit int, beforeID, afterID string) (st []*UserGuild, err error)

	// UserNoteSet sets the note for a specific user.
	UserNoteSet(userID string, message string) (err error)

	// UserSettings returns the settings for a given user
	UserSettings() (st *Settings, err error)

	// UserUpdate updates a users settings.
	UserUpdate(email, password, username, avatar, newPassword string) (st *User, err error)
}

// Voicer is the interface type to describe *Session functionality for working
// with voice channels.
type Voicer interface {
	// ChannelVoiceJoin joins the session user to a voice channel.
	//    gID     : Guild ID of the channel to join.
	//    cID     : Channel ID of the channel to join.
	//    mute    : If true, you will be set to muted upon joining.
	//    deaf    : If true, you will be set to deafened upon joining.
	ChannelVoiceJoin(gID, cID string, mute, deaf bool) (voice *VoiceConnection, err error)

	// ChannelVoiceJoinManual initiates a voice session to a voice channel, but does not complete it.
	//
	// This should only be used when the VoiceServerUpdate will be intercepted and used elsewhere.
	//
	//    gID     : Guild ID of the channel to join.
	//    cID     : Channel ID of the channel to join.
	//    mute    : If true, you will be set to muted upon joining.
	//    deaf    : If true, you will be set to deafened upon joining.
	ChannelVoiceJoinManual(gID, cID string, mute, deaf bool) (err error)

	// VoiceICE returns the voice server ICE information
	VoiceICE() (st *VoiceICE, err error)

	// VoiceRegions returns the voice server regions
	VoiceRegions() (st []*VoiceRegion, err error)
}

// i'm going to make my own Discord Go client with blackjack and webhookers...
//
// and you know what, forget the blackjack

// Webhooker is the interface type to describe *Session functionality for
// interacting with webhook configuration.
type Webhooker interface {
	// ChannelWebhooks returns all webhooks for a given channel.
	//   channelID: The ID of a channel.
	ChannelWebhooks(channelID string) (st []*Webhook, err error)

	// GuildWebhooks returns all webhooks for a given guild.
	//   guildID: The ID of a Guild.
	GuildWebhooks(guildID string) (st []*Webhook, err error)

	// Webhook returns a webhook for a given ID
	//   webhookID: The ID of a webhook.
	Webhook(webhookID string) (st *Webhook, err error)

	// WebhookCreate returns a new Webhook.
	//   channelID: The ID of a Channel.
	//   name     : The name of the webhook.
	//   avatar   : The avatar of the webhook.
	WebhookCreate(channelID, name, avatar string) (st *Webhook, err error)

	// WebhookDelete deletes a webhook for a given ID
	//   webhookID: The ID of a webhook.
	WebhookDelete(webhookID string) (err error)

	// WebhookDeleteWithToken deletes a webhook for a given ID with an auth token.
	//   webhookID: The ID of a webhook.
	//   token    : The auth token for the webhook.
	WebhookDeleteWithToken(webhookID, token string) (st *Webhook, err error)

	// WebhookEdit updates an existing Webhook.
	//   webhookID: The ID of a webhook.
	//   name     : The name of the webhook.
	//   avatar   : The avatar of the webhook.
	WebhookEdit(webhookID, name, avatar, channelID string) (st *Role, err error)

	// WebhookEditWithToken updates an existing Webhook with an auth token.
	//   webhookID: The ID of a webhook.
	//   token    : The auth token for the webhook.
	//   name     : The name of the webhook.
	//   avatar   : The avatar of the webhook.
	WebhookEditWithToken(webhookID, token, name, avatar string) (st *Role, err error)

	// WebhookExecute executes a webhook.
	//   webhookID: The ID of a webhook.
	//   token    : The auth token for the webhook
	WebhookExecute(webhookID, token string, wait bool, data *WebhookParams) (err error)

	// WebhookWithToken returns a webhook for a given ID
	//   webhookID: The ID of a webhook.
	//   token    : The auth token for the webhook.
	WebhookWithToken(webhookID, token string) (st *Webhook, err error)
}

// Sessioner is the interface to describe a *discordgo.Session for dependency
// injection in handlers. This composes smaller interfaces in to one giant one.
type Sessioner interface {
	// Open creates a websocket connection to Discord.
	// See: https://discordapp.com/developers/docs/topics/gateway#connecting
	Open() error

	// Close closes a websocket and stops all listening/heartbeat goroutines.
	Close() (err error)

	// HeartbeatLatency returns the latency between heartbeat acknowledgement
	// and heartbeat send.
	HeartbeatLatency() time.Duration

	// State returns the inner state that's updated when StateEnabled is true.
	State() *State

	// Applicationer is for configuring guild apps
	Applicationer

	// Channeler is for working with channels
	Channeler

	// ChannelFileSender is for sending files to the channel
	ChannelFileSender

	// ChannelMessager is for working with messages
	ChannelMessager

	// Guilder is for working with guilds
	Guilder

	// GuildAuditLoger is for working with the guild audit log
	GuildAuditLoger

	// GuildBanner is for working with user bans from the guild
	GuildBanner

	// GuildEmbeder is for working with guild embeds
	GuildEmbeder

	// GuildEmojier is for working with guild emoji
	GuildEmojier

	// GuildIntegrationer is for working with guild integrations
	GuildIntegrationer

	// GuildMemberer is for working with guild members
	GuildMemberer

	// GuildRoler is for working with guild roles
	GuildRoler

	// Inviter is for working with invites
	Inviter

	// MessageReactioner is for working with message reactions
	MessageReactioner

	// Relationshiper is for working with user relationships
	Relationshiper

	// Requester is for internal request functions
	Requester

	// StatusUpdater is for updating different user statuses
	StatusUpdater

	// Userer is for working with user data
	Userer

	// Voicer is for working with voice channels
	Voicer

	// Webhooker is for webhook management
	Webhooker

	AddHandler(handler interface{}) func()
	AddHandlerOnce(handler interface{}) func()

	Gateway() (gateway string, err error)
	GatewayBot() (st *GatewayBotResponse, err error)

	Login(email, password string) (err error)
	Logout() (err error)

	Register(username string) (token string, err error)
}

// compile-time check to ensure *Session satisfies Sessioner
var _ Sessioner = (*Session)(nil)
