package discordgo

import (
	"errors"
	"strings"
)

// A Guild holds all data related to a specific Discord Guild.  Guilds are also
// sometimes referred to as Servers in the Discord client.
type Guild struct {
	// The ID of the guild.
	ID string `json:"id"`

	// The name of the guild. (2–100 characters)
	Name string `json:"name"`

	// The hash of the guild's icon. Use Session.GuildIcon
	// to retrieve the icon itself.
	Icon string `json:"icon"`

	// The voice region of the guild.
	Region string `json:"region"`

	// The ID of the AFK voice channel.
	AfkChannelID string `json:"afk_channel_id"`

	// The ID of the embed channel ID, used for embed widgets.
	EmbedChannelID string `json:"embed_channel_id"`

	// The user ID of the owner of the guild.
	OwnerID string `json:"owner_id"`

	// The time at which the current user joined the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	JoinedAt Timestamp `json:"joined_at"`

	// The hash of the guild's splash.
	Splash string `json:"splash"`

	// The timeout, in seconds, before a user is considered AFK in voice.
	AfkTimeout int `json:"afk_timeout"`

	// The number of members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	MemberCount int `json:"member_count"`

	// The verification level required for the guild.
	VerificationLevel VerificationLevel `json:"verification_level"`

	// Whether the guild has embedding enabled.
	EmbedEnabled bool `json:"embed_enabled"`

	// Whether the guild is considered large. This is
	// determined by a member threshold in the identify packet,
	// and is currently hard-coded at 250 members in the library.
	Large bool `json:"large"`

	// The default message notification setting for the guild.
	// 0 == all messages, 1 == mentions only.
	DefaultMessageNotifications int `json:"default_message_notifications"`

	// A list of roles in the guild.
	Roles []*Role `json:"roles"`

	// A list of the custom emojis present in the guild.
	Emojis []*Emoji `json:"emojis"`

	// A list of the members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Members []*Member `json:"members"`

	// A list of partial presence objects for members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Presences []*Presence `json:"presences"`

	// A list of channels in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Channels []*Channel `json:"channels"`

	// A list of voice states for the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	VoiceStates []*VoiceState `json:"voice_states"`

	// Whether this guild is currently unavailable (most likely due to outage).
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Unavailable bool `json:"unavailable"`

	// The explicit content filter level
	ExplicitContentFilter ExplicitContentFilterLevel `json:"explicit_content_filter"`

	// The list of enabled guild features
	Features []string `json:"features"`

	// Required MFA level for the guild
	MfaLevel MfaLevel `json:"mfa_level"`

	// Whether or not the Server Widget is enabled
	WidgetEnabled bool `json:"widget_enabled"`

	// The Channel ID for the Server Widget
	WidgetChannelID string `json:"widget_channel_id"`

	// The Channel ID to which system messages are sent (eg join and leave messages)
	SystemChannelID string `json:"system_channel_id"`

	// The Session to call the API and retrieve other objects
	Session *Session `json:"session,omitempty"`
}

// A UserGuild holds a brief version of a Guild
type UserGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Owner       bool   `json:"owner"`
	Permissions int    `json:"permissions"`
}

// A GuildParams stores all the data needed to update discord guild settings
type GuildParams struct {
	Name                        string             `json:"name,omitempty"`
	Region                      string             `json:"region,omitempty"`
	VerificationLevel           *VerificationLevel `json:"verification_level,omitempty"`
	DefaultMessageNotifications int                `json:"default_message_notifications,omitempty"` // TODO: Separate type?
	AfkChannelID                string             `json:"afk_channel_id,omitempty"`
	AfkTimeout                  int                `json:"afk_timeout,omitempty"`
	Icon                        string             `json:"icon,omitempty"`
	OwnerID                     string             `json:"owner_id,omitempty"`
	Splash                      string             `json:"splash,omitempty"`
}

// gets the role with the given ID as it is sored in Guild.Roles
func (g *Guild) GetRole(roleID string) (role *Role, err error) {
	for _, role = range g.Roles {
		if role.ID == roleID {
			return role, nil
		}
	}

	err = errors.New("no role found")
	return
}

// gets the role with the given name as it is stored in Guild.Roles
// It is semi-case-sensitive; if a name matches full, the first role with that name gets returned
// if a name matches but with different capitalization, the last role with that name gets returned
func (g *Guild) GetRoleNamed(name string) (role *Role, err error) {
	var savedRole *Role
	lowerCaseName := strings.ToLower(name)

	for _, role = range g.Roles {
		if role.Name == name {
			return
		} else if role.Name == lowerCaseName {
			savedRole = role
		}
	}

	if savedRole != nil {
		role = savedRole
		return
	}

	err = errors.New("no role found")
	return
}

// gets the channel with the given ID as it is sored in Guild.Channels
func (g *Guild) GetChannel(channelID string) (role *Channel, err error) {
	for _, channel := range g.Channels {
		if channel.ID == channelID {
			return
		}
	}

	err = errors.New("no channel found")
	return
}

// gets the channel with the given name as it is stored in Guild.Channels
// It is semi-case-sensitive; if a name matches full, the first channel with that name gets returned
// if a name matches but with different capitalization, the last channel with that name gets returned
func (g *Guild) GetChannelNamed(name string) (channel *Channel, err error) {
	var savedChannel *Channel
	lowerCaseName := strings.ToLower(name)

	for _, channel = range g.Channels {
		if channel.Name == name {
			return
		} else if strings.ToLower(channel.Name) == lowerCaseName {
			savedChannel = channel
		}
	}

	if savedChannel != nil {
		channel = savedChannel
		return
	}

	err = errors.New("no channel found")
	return
}
