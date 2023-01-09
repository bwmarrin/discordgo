package discordgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

// This file contains all the possible structs that can be
// handled by AddHandler/EventHandler.
// DO NOT ADD ANYTHING BUT EVENT HANDLER STRUCTS TO THIS FILE.
//go:generate go run tools/cmd/eventhandlers/main.go

// Connect is the data for a Connect event.
// This is a synthetic event and is not dispatched by Discord.
type Connect struct{}

// Disconnect is the data for a Disconnect event.
// This is a synthetic event and is not dispatched by Discord.
type Disconnect struct{}

// RateLimit is the data for a RateLimit event.
// This is a synthetic event and is not dispatched by Discord.
type RateLimit struct {
	*TooManyRequests
	URL string
}

// Event provides a basic initial struct for all websocket events.
type Event struct {
	Operation int             `json:"op"`
	Sequence  int64           `json:"s"`
	Type      string          `json:"t"`
	RawData   json.RawMessage `json:"d"`
	// Struct contains one of the other types in this file.
	Struct interface{} `json:"-"`
}

// A Ready stores all data for the websocket READY event.
type Ready struct {
	Version           int                  `json:"v"`
	SessionID         string               `json:"session_id"`
	User              *User                `json:"user"`
	ReadState         []*ReadState         `json:"read_state"`
	PrivateChannels   []*Channel           `json:"private_channels"`
	Guilds            []*Guild             `json:"guilds"`
	Sessions          []*DiscordSession    `json:"sessions"`
	UserGuildSettings []*UserGuildSettings `json:"user_guild_settings"`
	Relationships     []*Relationship      `json:"relationships"`
}

func (r *Ready) UnmarshalJSON(data []byte) error {
	type rawReady struct {
		Version           int                    `json:"v"`
		Users             []*User                `json:"users"`
		UserGuildSettings *UserGuildSettingsData `json:"user_guild_settings"`
		User              *User                  `json:"user"`
		Sessions          []*DiscordSession      `json:"sessions"`
		SessionID         string                 `json:"session_id"`
		Relationships     []*Relationship        `json:"relationships"`
		ReadState         *ReadStateData         `json:"read_state"`
		PrivateChannels   []*Channel             `json:"private_channels"`
		Guilds            []*ReadyGuild          `json:"guilds"`
		CountryCode       string                 `json:"country_code"`
		ConnectedAccounts []*UserConnection      `json:"connected_accounts"`
	}

	var ready rawReady

	ioutil.WriteFile("C:\\Users\\kaani\\Desktop\\ready.json", data, os.ModePerm)

	if err := json.Unmarshal(data, &ready); err != nil {
		return err
	}

	r.Version = ready.Version
	r.SessionID = ready.SessionID
	r.User = ready.User
	r.ReadState = ready.ReadState.Entries
	r.PrivateChannels = ready.PrivateChannels
	r.Sessions = ready.Sessions
	r.UserGuildSettings = ready.UserGuildSettings.Entries
	r.Relationships = ready.Relationships

	for _, readyGuild := range ready.Guilds {
		guild := Guild{
			Properties:                  readyGuild.Properties,
			ID:                          readyGuild.ID,
			Large:                       readyGuild.Large,
			Lazy:                        readyGuild.Lazy,
			MemberCount:                 readyGuild.MemberCount,
			PremiumSubscriptionCount:    readyGuild.PremiumSubscriptionCount,
			JoinedAt:                    readyGuild.JoinedAt,
			Threads:                     readyGuild.Threads,
			Stickers:                    readyGuild.Stickers,
			StageInstances:              readyGuild.StageInstances,
			ScheduledEvents:             readyGuild.ScheduledEvents,
			Roles:                       readyGuild.Roles,
			Channels:                    readyGuild.Channels,
			Emojis:                      readyGuild.Emojis,
			VerificationLevel:           readyGuild.Properties.VerificationLevel,
			VanityURLCode:               readyGuild.Properties.VanityURLCode,
			SystemChannelID:             readyGuild.Properties.SystemChannelID,
			SystemChannelFlags:          readyGuild.Properties.SystemChannelFlags,
			Splash:                      readyGuild.Properties.Splash,
			SafetyAlertsChannelID:       readyGuild.Properties.SafetyAlertsChannelID,
			RulesChannelID:              readyGuild.Properties.RulesChannelID,
			PublicUpdatesChannelID:      readyGuild.Properties.PublicUpdatesChannelID,
			PremiumTier:                 readyGuild.Properties.PremiumTier,
			PremiumProgressBarEnabled:   readyGuild.Properties.PremiumProgressBarEnabled,
			PreferredLocale:             readyGuild.Properties.PreferredLocale,
			OwnerID:                     readyGuild.Properties.OwnerID,
			NSFWLevel:                   readyGuild.Properties.NSFWLevel,
			NSFW:                        readyGuild.Properties.NSFW,
			Name:                        readyGuild.Properties.Name,
			MfaLevel:                    readyGuild.Properties.MfaLevel,
			MaxVideoChannelUsers:        readyGuild.Properties.MaxVideoChannelUsers,
			MaxStageVideoChannelUsers:   readyGuild.Properties.MaxStageVideoChannelUsers,
			MaxMembers:                  readyGuild.Properties.MaxMembers,
			Icon:                        readyGuild.Properties.Icon,
			HubType:                     readyGuild.Properties.HubType,
			HomeHeader:                  readyGuild.Properties.HomeHeader,
			Features:                    readyGuild.Properties.Features,
			ExplicitContentFilter:       readyGuild.Properties.ExplicitContentFilter,
			DiscoverySplash:             readyGuild.Properties.DiscoverySplash,
			Description:                 readyGuild.Properties.Description,
			DefaultMessageNotifications: readyGuild.Properties.DefaultMessageNotifications,
			Banner:                      readyGuild.Properties.Banner,
			ApplicationID:               readyGuild.Properties.ApplicationID,
			AfkTimeout:                  readyGuild.Properties.AfkTimeout,
			AfkChannelID:                readyGuild.Properties.AfkChannelID,
		}

		r.Guilds = append(r.Guilds, &guild)
	}

	runtime.GC()

	return nil
}

// ChannelCreate is the data for a ChannelCreate event.
type ChannelCreate struct {
	*Channel
}

// ChannelUpdate is the data for a ChannelUpdate event.
type ChannelUpdate struct {
	*Channel
}

// ChannelDelete is the data for a ChannelDelete event.
type ChannelDelete struct {
	*Channel
}

// ChannelPinsUpdate stores data for a ChannelPinsUpdate event.
type ChannelPinsUpdate struct {
	LastPinTimestamp string `json:"last_pin_timestamp"`
	ChannelID        string `json:"channel_id"`
	GuildID          string `json:"guild_id,omitempty"`
}

// ThreadCreate is the data for a ThreadCreate event.
type ThreadCreate struct {
	*Channel
	NewlyCreated bool `json:"newly_created"`
}

// ThreadUpdate is the data for a ThreadUpdate event.
type ThreadUpdate struct {
	*Channel
	BeforeUpdate *Channel `json:"-"`
}

// ThreadDelete is the data for a ThreadDelete event.
type ThreadDelete struct {
	*Channel
}

// ThreadListSync is the data for a ThreadListSync event.
type ThreadListSync struct {
	// The id of the guild
	GuildID string `json:"guild_id"`
	// The parent channel ids whose threads are being synced.
	// If omitted, then threads were synced for the entire guild.
	// This array may contain channel_ids that have no active threads as well, so you know to clear that data.
	ChannelIDs []string `json:"channel_ids"`
	// All active threads in the given channels that the current user can access
	Threads []*Channel `json:"threads"`
	// All thread member objects from the synced threads for the current user,
	// indicating which threads the current user has been added to
	Members []*ThreadMember `json:"members"`
}

// ThreadMemberUpdate is the data for a ThreadMemberUpdate event.
type ThreadMemberUpdate struct {
	*ThreadMember
	GuildID string `json:"guild_id"`
}

// ThreadMembersUpdate is the data for a ThreadMembersUpdate event.
type ThreadMembersUpdate struct {
	ID             string              `json:"id"`
	GuildID        string              `json:"guild_id"`
	MemberCount    int                 `json:"member_count"`
	AddedMembers   []AddedThreadMember `json:"added_members"`
	RemovedMembers []string            `json:"removed_member_ids"`
}

// GuildCreate is the data for a GuildCreate event.
type GuildCreate struct {
	*Guild
}

// GuildUpdate is the data for a GuildUpdate event.
type GuildUpdate struct {
	*Guild
}

// GuildDelete is the data for a GuildDelete event.
type GuildDelete struct {
	*Guild
	BeforeDelete *Guild `json:"-"`
}

// GuildBanAdd is the data for a GuildBanAdd event.
type GuildBanAdd struct {
	User    *User  `json:"user"`
	GuildID string `json:"guild_id"`
}

// GuildBanRemove is the data for a GuildBanRemove event.
type GuildBanRemove struct {
	User    *User  `json:"user"`
	GuildID string `json:"guild_id"`
}

// GuildMemberAdd is the data for a GuildMemberAdd event.
type GuildMemberAdd struct {
	*Member
}

// GuildMemberUpdate is the data for a GuildMemberUpdate event.
type GuildMemberUpdate struct {
	*Member
	BeforeUpdate *Member `json:"-"`
}

// GuildMemberRemove is the data for a GuildMemberRemove event.
type GuildMemberRemove struct {
	*Member
}

// GuildRoleCreate is the data for a GuildRoleCreate event.
type GuildRoleCreate struct {
	*GuildRole
}

// GuildRoleUpdate is the data for a GuildRoleUpdate event.
type GuildRoleUpdate struct {
	*GuildRole
}

// A GuildRoleDelete is the data for a GuildRoleDelete event.
type GuildRoleDelete struct {
	RoleID  string `json:"role_id"`
	GuildID string `json:"guild_id"`
}

// A GuildEmojisUpdate is the data for a guild emoji update event.
type GuildEmojisUpdate struct {
	GuildID string   `json:"guild_id"`
	Emojis  []*Emoji `json:"emojis"`
}

// A GuildMembersChunk is the data for a GuildMembersChunk event.
type GuildMembersChunk struct {
	GuildID    string      `json:"guild_id"`
	Members    []*Member   `json:"members"`
	ChunkIndex int         `json:"chunk_index"`
	ChunkCount int         `json:"chunk_count"`
	NotFound   []string    `json:"not_found,omitempty"`
	Presences  []*Presence `json:"presences,omitempty"`
	Nonce      string      `json:"nonce,omitempty"`
}

type GuildMemberListUpdate struct {
	OnlineCount int              `json:"online_count"`
	MemberCount int              `json:"member_count"`
	ID          string           `json:"id"`
	GuildID     string           `json:"guild_id"`
	Ops         []*Operator      `json:"ops"`
	Groups      []*SyncItemGroup `json:"groups"`
}

// GuildIntegrationsUpdate is the data for a GuildIntegrationsUpdate event.
type GuildIntegrationsUpdate struct {
	GuildID string `json:"guild_id"`
}

// StageInstanceEventCreate is the data for a StageInstanceEventCreate event.
type StageInstanceEventCreate struct {
	*StageInstance
}

// StageInstanceEventUpdate is the data for a StageInstanceEventUpdate event.
type StageInstanceEventUpdate struct {
	*StageInstance
}

// StageInstanceEventDelete is the data for a StageInstanceEventDelete event.
type StageInstanceEventDelete struct {
	*StageInstance
}

// GuildScheduledEventCreate is the data for a GuildScheduledEventCreate event.
type GuildScheduledEventCreate struct {
	*GuildScheduledEvent
}

// GuildScheduledEventUpdate is the data for a GuildScheduledEventUpdate event.
type GuildScheduledEventUpdate struct {
	*GuildScheduledEvent
}

// GuildScheduledEventDelete is the data for a GuildScheduledEventDelete event.
type GuildScheduledEventDelete struct {
	*GuildScheduledEvent
}

// GuildScheduledEventUserAdd is the data for a GuildScheduledEventUserAdd event.
type GuildScheduledEventUserAdd struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

// GuildScheduledEventUserRemove is the data for a GuildScheduledEventUserRemove event.
type GuildScheduledEventUserRemove struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

// MessageCreate is the data for a MessageCreate event.
type MessageCreate struct {
	*Message
}

// UnmarshalJSON is a helper function to unmarshal MessageCreate object.
func (m *MessageCreate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

// MessageUpdate is the data for a MessageUpdate event.
type MessageUpdate struct {
	*Message
	// BeforeUpdate will be nil if the Message was not previously cached in the state cache.
	BeforeUpdate *Message `json:"-"`
}

// UnmarshalJSON is a helper function to unmarshal MessageUpdate object.
func (m *MessageUpdate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

// MessageDelete is the data for a MessageDelete event.
type MessageDelete struct {
	*Message
	BeforeDelete *Message `json:"-"`
}

// UnmarshalJSON is a helper function to unmarshal MessageDelete object.
func (m *MessageDelete) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

// MessageReactionAdd is the data for a MessageReactionAdd event.
type MessageReactionAdd struct {
	*MessageReaction
	Member *Member `json:"member,omitempty"`
}

// MessageReactionRemove is the data for a MessageReactionRemove event.
type MessageReactionRemove struct {
	*MessageReaction
}

// MessageReactionRemoveAll is the data for a MessageReactionRemoveAll event.
type MessageReactionRemoveAll struct {
	*MessageReaction
}

// PresencesReplace is the data for a PresencesReplace event.
type PresencesReplace []*Presence

// PresenceUpdate is the data for a PresenceUpdate event.
type PresenceUpdate struct {
	Presence
	GuildID string `json:"guild_id"`
}

// Resumed is the data for a Resumed event.
type Resumed struct {
	Trace []string `json:"_trace"`
}

// RelationshipAdd is the data for a RelationshipAdd event.
type RelationshipAdd struct {
	*Relationship
}

// RelationshipRemove is the data for a RelationshipRemove event.
type RelationshipRemove struct {
	*Relationship
}

// TypingStart is the data for a TypingStart event.
type TypingStart struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
	Timestamp int    `json:"timestamp"`
}

// UserUpdate is the data for a UserUpdate event.
type UserUpdate struct {
	*User
}

// VoiceServerUpdate is the data for a VoiceServerUpdate event.
type VoiceServerUpdate struct {
	Token    string `json:"token"`
	GuildID  string `json:"guild_id"`
	Endpoint string `json:"endpoint"`
}

// VoiceStateUpdate is the data for a VoiceStateUpdate event.
type VoiceStateUpdate struct {
	*VoiceState
	// BeforeUpdate will be nil if the VoiceState was not previously cached in the state cache.
	BeforeUpdate *VoiceState `json:"-"`
}

// MessageDeleteBulk is the data for a MessageDeleteBulk event
type MessageDeleteBulk struct {
	Messages  []string `json:"ids"`
	ChannelID string   `json:"channel_id"`
	GuildID   string   `json:"guild_id"`
}

// WebhooksUpdate is the data for a WebhooksUpdate event
type WebhooksUpdate struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}

// InteractionCreate is the data for a InteractionCreate event
type InteractionCreate struct {
	*Interaction
}

// UnmarshalJSON is a helper function to unmarshal Interaction object.
func (i *InteractionCreate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &i.Interaction)
}

// InviteCreate is the data for a InviteCreate event
type InviteCreate struct {
	*Invite
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
}

// InviteDelete is the data for a InviteDelete event
type InviteDelete struct {
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	Code      string `json:"code"`
}

// ApplicationCommandPermissionsUpdate is the data for an ApplicationCommandPermissionsUpdate event
type ApplicationCommandPermissionsUpdate struct {
	*GuildApplicationCommandPermissions
}

// AutoModerationRuleCreate is the data for an AutoModerationRuleCreate event.
type AutoModerationRuleCreate struct {
	*AutoModerationRule
}

// AutoModerationRuleUpdate is the data for an AutoModerationRuleUpdate event.
type AutoModerationRuleUpdate struct {
	*AutoModerationRule
}

// AutoModerationRuleDelete is the data for an AutoModerationRuleDelete event.
type AutoModerationRuleDelete struct {
	*AutoModerationRule
}

// AutoModerationActionExecution is the data for an AutoModerationActionExecution event.
type AutoModerationActionExecution struct {
	GuildID              string                        `json:"guild_id"`
	Action               AutoModerationAction          `json:"action"`
	RuleID               string                        `json:"rule_id"`
	RuleTriggerType      AutoModerationRuleTriggerType `json:"rule_trigger_type"`
	UserID               string                        `json:"user_id"`
	ChannelID            string                        `json:"channel_id"`
	MessageID            string                        `json:"message_id"`
	AlertSystemMessageID string                        `json:"alert_system_message_id"`
	Content              string                        `json:"content"`
	MatchedKeyword       string                        `json:"matched_keyword"`
	MatchedContent       string                        `json:"matched_content"`
}
