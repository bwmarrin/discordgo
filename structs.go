// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains all structures for the discordgo package.  These
// may be moved about later into separate files but I find it easier to have
// them all located together.

package discordgo

import (
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// A Session represents a connection to the Discord API.
type Session struct {
	sync.RWMutex

	// General configurable settings.

	// Authentication token for this session
	Token string

	// Debug for printing JSON request/responses
	Debug bool

	// Should the session reconnect the websocket on errors.
	ShouldReconnectOnError bool

	// Should the session request compressed websocket data.
	Compress bool

	// Should state tracking be enabled.
	// State tracking is the best way for getting the the users
	// active guilds and the members of the guilds.
	StateEnabled bool

	// Exposed but should not be modified by User.

	// Whether the Data Websocket is ready
	DataReady bool

	// Whether the Voice Websocket is ready
	VoiceReady bool

	// Whether the UDP Connection is ready
	UDPReady bool

	// Stores all details related to voice connections
	Voice *Voice

	// Managed state object, updated internally with events when
	// StateEnabled is true.
	State *State

	handlersMu sync.RWMutex
	// This is a mapping of event struct to a reflected value
	// for event handlers.
	// We store the reflected value instead of the function
	// reference as it is more performant, instead of re-reflecting
	// the function each event.
	handlers map[interface{}][]reflect.Value

	// The websocket connection.
	wsConn *websocket.Conn

	// When nil, the session is not listening.
	listening chan interface{}
}

// A VoiceRegion stores data for a specific voice region server.
type VoiceRegion struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"sample_hostname"`
	Port     int    `json:"sample_port"`
}

// A VoiceICE stores data for voice ICE servers.
type VoiceICE struct {
	TTL     string       `json:"ttl"`
	Servers []*ICEServer `json:"servers"`
}

// A ICEServer stores data for a specific voice ICE server.
type ICEServer struct {
	URL        string `json:"url"`
	Username   string `json:"username"`
	Credential string `json:"credential"`
}

// A Invite stores all data related to a specific Discord Guild or Channel invite.
type Invite struct {
	Guild     *Guild   `json:"guild"`
	Channel   *Channel `json:"channel"`
	Inviter   *User    `json:"inviter"`
	Code      string   `json:"code"`
	CreatedAt string   `json:"created_at"` // TODO make timestamp
	MaxAge    int      `json:"max_age"`
	Uses      int      `json:"uses"`
	MaxUses   int      `json:"max_uses"`
	XkcdPass  bool     `json:"xkcdpass"`
	Revoked   bool     `json:"revoked"`
	Temporary bool     `json:"temporary"`
}

// A Channel holds all data related to an individual Discord channel.
type Channel struct {
	ID                   string                 `json:"id"`
	GuildID              string                 `json:"guild_id"`
	Name                 string                 `json:"name"`
	Topic                string                 `json:"topic"`
	Position             int                    `json:"position"`
	Bitrate              int                    `json:"bitrate"`
	Type                 string                 `json:"type"`
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites"`
	IsPrivate            bool                   `json:"is_private"`
	LastMessageID        string                 `json:"last_message_id"`
	Recipient            *User                  `json:"recipient"`
	Messages             []*Message             `json:"-"`
}

// A PermissionOverwrite holds permission overwrite data for a Channel
type PermissionOverwrite struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Deny  int    `json:"deny"`
	Allow int    `json:"allow"`
}

// Emoji struct holds data related to Emoji's
type Emoji struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Roles         []string `json:"roles"`
	Managed       bool     `json:"managed"`
	RequireColons bool     `json:"require_colons"`
}

// A Guild holds all data related to a specific Discord Guild.  Guilds are also
// sometimes referred to as Servers in the Discord client.
type Guild struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Icon           string        `json:"icon"`
	Region         string        `json:"region"`
	AfkChannelID   string        `json:"afk_channel_id"`
	EmbedChannelID string        `json:"embed_channel_id"`
	OwnerID        string        `json:"owner_id"`
	JoinedAt       string        `json:"joined_at"` // make this a timestamp
	Splash         string        `json:"splash"`
	AfkTimeout     int           `json:"afk_timeout"`
	EmbedEnabled   bool          `json:"embed_enabled"`
	Large          bool          `json:"large"` // ??
	Roles          []*Role       `json:"roles"`
	Emojis         []*Emoji      `json:"emojis"`
	Members        []*Member     `json:"members"`
	Presences      []*Presence   `json:"presences"`
	Channels       []*Channel    `json:"channels"`
	VoiceStates    []*VoiceState `json:"voice_states"`
}

// A Role stores information about Discord guild member roles.
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Managed     bool   `json:"managed"`
	Hoist       bool   `json:"hoist"`
	Color       int    `json:"color"`
	Position    int    `json:"position"`
	Permissions int    `json:"permissions"`
}

// A VoiceState stores the voice states of Guilds
type VoiceState struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	ChannelID string `json:"channel_id"`
	Suppress  bool   `json:"suppress"`
	SelfMute  bool   `json:"self_mute"`
	SelfDeaf  bool   `json:"self_deaf"`
	Mute      bool   `json:"mute"`
	Deaf      bool   `json:"deaf"`
}

// A Presence stores the online, offline, or idle and game status of Guild members.
type Presence struct {
	User   *User  `json:"user"`
	Status string `json:"status"`
	Game   *Game  `json:"game"`
}

// A Game struct holds the name of the "playing .." game for a user
type Game struct {
	Name string `json:"name"`
}

// A Member stores user information for Guild members.
type Member struct {
	GuildID  string   `json:"guild_id"`
	JoinedAt string   `json:"joined_at"`
	Deaf     bool     `json:"deaf"`
	Mute     bool     `json:"mute"`
	User     *User    `json:"user"`
	Roles    []string `json:"roles"`
}

// A User stores all data for an individual Discord user.
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Avatar   string `json:"Avatar"`
	Verified bool   `json:"verified"`
	//Discriminator int    `json:"discriminator,string"` // TODO: See below
}

// TODO: Research issue.
// Discriminator sometimes comes as a string
// and sometimes it comes as a int.  Weird.
// to avoid errors I've just commented it out
// but it doesn't seem to just kill the whole
// program.  Heartbeat is taken on READY even
// with error and the system continues to read
// it just doesn't seem able to handle this one
// field correctly.  Need to research this more.

// A Settings stores data for a specific users Discord client settings.
type Settings struct {
	RenderEmbeds          bool     `json:"render_embeds"`
	InlineEmbedMedia      bool     `json:"inline_embed_media"`
	EnableTtsCommand      bool     `json:"enable_tts_command"`
	MessageDisplayCompact bool     `json:"message_display_compact"`
	ShowCurrentGame       bool     `json:"show_current_game"`
	Locale                string   `json:"locale"`
	Theme                 string   `json:"theme"`
	MutedChannels         []string `json:"muted_channels"`
}

// An Event provides a basic initial struct for all websocket event.
type Event struct {
	Type      string          `json:"t"`
	State     int             `json:"s"`
	Operation int             `json:"op"`
	Direction int             `json:"dir"`
	RawData   json.RawMessage `json:"d"`
}

// A Ready stores all data for the websocket READY event.
type Ready struct {
	Version           int           `json:"v"`
	SessionID         string        `json:"session_id"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	User              *User         `json:"user"`
	ReadState         []*ReadState  `json:"read_state"`
	PrivateChannels   []*Channel    `json:"private_channels"`
	Guilds            []*Guild      `json:"guilds"`
}

// A RateLimit struct holds information related to a specific rate limit.
type RateLimit struct {
	Bucket     string        `json:"bucket"`
	Message    string        `json:"message"`
	RetryAfter time.Duration `json:"retry_after"`
}

// A ReadState stores data on the read state of channels.
type ReadState struct {
	MentionCount  int
	LastMessageID string `json:"last_message_id"`
	ID            string `json:"id"`
}

// A TypingStart stores data for the typing start websocket event.
type TypingStart struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	Timestamp int    `json:"timestamp"`
}

// A PresenceUpdate stores data for the pressence update websocket event.
type PresenceUpdate struct {
	User    *User    `json:"user"`
	Status  string   `json:"status"`
	Roles   []string `json:"roles"`
	GuildID string   `json:"guild_id"`
	Game    *Game    `json:"game"`
}

// A MessageAck stores data for the message ack websocket event.
type MessageAck struct {
	MessageID string `json:"message_id"`
	ChannelID string `json:"channel_id"`
}

// A GuildIntegrationsUpdate stores data for the guild integrations update
// websocket event.
type GuildIntegrationsUpdate struct {
	GuildID string `json:"guild_id"`
}

// A GuildRole stores data for guild role websocket events.
type GuildRole struct {
	Role    *Role  `json:"role"`
	GuildID string `json:"guild_id"`
}

// A GuildRoleDelete stores data for the guild role delete websocket event.
type GuildRoleDelete struct {
	RoleID  string `json:"role_id"`
	GuildID string `json:"guild_id"`
}

// A GuildBan stores data for a guild ban.
type GuildBan struct {
	User    *User  `json:"user"`
	GuildID string `json:"guild_id"`
}

// A GuildEmojisUpdate stores data for a guild emoji update event.
type GuildEmojisUpdate struct {
	GuildID string   `json:"guild_id"`
	Emojis  []*Emoji `json:"emojis"`
}

// A State contains the current known state.
// As discord sends this in a READY blob, it seems reasonable to simply
// use that struct as the data store.
type State struct {
	sync.RWMutex
	Ready
	MaxMessageCount int
}
