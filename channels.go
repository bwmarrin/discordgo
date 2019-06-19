package discordgo

import (
	"errors"
	"fmt"
	"time"
)

// ErrNotATextChannel gets returned when an action gets called on a channel
// that does not support sending messages to them
var ErrNotATextChannel = errors.New("not a text or dm channel")

// ErrNotAVoiceChannel gets thrown when an action gets called on a channel
// that is not a Guild Voice channel
var ErrNotAVoiceChannel = errors.New("not a voice channel")

// A Channel holds all data related to an individual Discord channel.
type Channel struct {
	// The ID of the channel.
	ID string `json:"id"`

	// The ID of the guild to which the channel belongs, if it is in a guild.
	// Else, this ID is empty (e.g. DM channels).
	GuildID string `json:"guild_id"`

	// The name of the channel.
	Name string `json:"name"`

	// The topic of the channel.
	Topic string `json:"topic"`

	// The type of the channel.
	Type ChannelType `json:"type"`

	// The ID of the last message sent in the channel. This is not
	// guaranteed to be an ID of a valid message.
	LastMessageID string `json:"last_message_id"`

	// Whether the channel is marked as NSFW.
	NSFW bool `json:"nsfw"`

	// Icon of the group DM channel.
	Icon string `json:"icon"`

	// The position of the channel, used for sorting in client.
	Position int `json:"position"`

	// The bitrate of the channel, if it is a voice channel.
	Bitrate int `json:"bitrate"`

	// The recipients of the channel. This is only populated in DM channels.
	Recipients []*User `json:"recipients"`

	// The messages in the channel. This is only present in state-cached channels,
	// and State.MaxMessageCount must be non-zero.
	Messages []*Message `json:"-"`

	// A list of permission overwrites present for the channel.
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites"`

	// The user limit of the voice channel.
	UserLimit int `json:"user_limit"`

	// The ID of the parent channel, if the channel is under a category
	ParentID string `json:"parent_id"`

	// The Session to call the API and retrieve other objects
	Session *Session `json:"session,omitempty"`
}

// Mention returns a string which mentions the channel
func (c *Channel) Mention() string {
	return fmt.Sprintf("<#%s>", c.ID)
}

// GetID returns the ID of the channel
func (c *Channel) GetID() string {
	return c.ID
}

// CreatedAt returns the channels creation time in UTC
func (c *Channel) CreatedAt() (creation time.Time, err error) {
	return SnowflakeToTime(c.ID)
}

// A ChannelEdit holds Channel Field data for a channel edit.
type ChannelEdit struct {
	Name                 string                 `json:"name,omitempty"`
	Topic                string                 `json:"topic,omitempty"`
	NSFW                 bool                   `json:"nsfw,omitempty"`
	Position             int                    `json:"position"`
	Bitrate              int                    `json:"bitrate,omitempty"`
	UserLimit            int                    `json:"user_limit,omitempty"`
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             string                 `json:"parent_id,omitempty"`
	RateLimitPerUser     int                    `json:"rate_limit_per_user,omitempty"`
}

// A PermissionOverwrite holds permission overwrite data for a Channel
type PermissionOverwrite struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Deny  int    `json:"deny"`
	Allow int    `json:"allow"`
}

// SendMessage sends a message to the channel
// content         : message content to send if provided
// embed           : embed to attach to the message if provided
// files           : files to attach to the message if provided
func (c *Channel) SendMessage(content string, embed *MessageEmbed, files []*File) (message *Message, err error) {
	if c.Type == ChannelTypeGuildVoice || c.Type == ChannelTypeGuildCategory {
		err = ErrNotATextChannel
		return
	}

	data := &MessageSend{
		Content: content,
		Embed:   embed,
		Files:   files,
	}

	return c.SendMessageComplex(data)
}

// SendMessageComplex sends a message to the channel
// data          : MessageSend object with the data to send
func (c *Channel) SendMessageComplex(data *MessageSend) (message *Message, err error) {
	if c.Type == ChannelTypeGuildVoice || c.Type == ChannelTypeGuildCategory {
		err = ErrNotATextChannel
		return
	}

	return c.Session.ChannelMessageSendComplex(c.ID, data)
}

// EditMessage edits a message, replacing it entirely with the corresponding
// fields in the given message struct
func (c *Channel) EditMessage(message *Message) (edited *Message, err error) {
	if c.Type == ChannelTypeGuildVoice || c.Type == ChannelTypeGuildCategory {
		err = ErrNotATextChannel
		return
	}

	data := &MessageEdit{
		ID:      message.ID,
		Channel: c.ID,
		Content: &message.Content,
	}
	if len(message.Embeds) > 0 {
		data.SetEmbed(message.Embeds[0])
	}

	return c.EditMessageComplex(data)
}

// EditMessageComplex edits an existing message, replacing it entirely with
// the given MessageEdit struct
func (c *Channel) EditMessageComplex(data *MessageEdit) (edited *Message, err error) {
	if c.Type == ChannelTypeGuildVoice || c.Type == ChannelTypeGuildCategory {
		err = ErrNotATextChannel
		return
	}

	data.Channel = c.ID
	return c.Session.ChannelMessageEditComplex(data)
}

// FetchMessage fetches a message with the given ID from the channel
// ID        : ID of the message to fetch
func (c *Channel) FetchMessage(ID string) (message *Message, err error) {
	if c.Type == ChannelTypeGuildVoice || c.Type == ChannelTypeGuildCategory {
		err = ErrNotATextChannel
		return
	}

	return c.Session.ChannelMessage(c.ID, ID)
}

// GetHistory fetches up to limit messages from the channel
// limit     : The number messages that can be returned. (max 100)
// beforeID  : If provided all messages returned will be before given ID.
// afterID   : If provided all messages returned will be after given ID.
// aroundID  : If provided all messages returned will be around given ID.
func (c *Channel) GetHistory(limit int, beforeID, afterID, aroundID string) (st []*Message, err error) {
	if c.Type == ChannelTypeGuildVoice || c.Type == ChannelTypeGuildCategory {
		err = ErrNotATextChannel
		return
	}

	return c.Session.ChannelMessages(c.ID, limit, beforeID, afterID, aroundID)
}
