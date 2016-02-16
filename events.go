package discordgo

// Connect is an empty struct for an event.
type Connect struct{}

// Disconnect is an empty struct for an event.
type Disconnect struct{}

// MessageCreate is a wrapper struct for an event.
type MessageCreate struct {
	*Message
}

// MessageUpdate is a wrapper struct for an event.
type MessageUpdate struct {
	*Message
}

// MessageDelete is a wrapper struct for an event.
type MessageDelete struct {
	*Message
}

// ChannelCreate is a wrapper struct for an event.
type ChannelCreate struct {
	*Channel
}

// ChannelUpdate is a wrapper struct for an event.
type ChannelUpdate struct {
	*Channel
}

// ChannelDelete is a wrapper struct for an event.
type ChannelDelete struct {
	*Channel
}

// GuildCreate is a wrapper struct for an event.
type GuildCreate struct {
	*Guild
}

// GuildUpdate is a wrapper struct for an event.
type GuildUpdate struct {
	*Guild
}

// GuildDelete is a wrapper struct for an event.
type GuildDelete struct {
	*Guild
}

// GuildBanAdd is a wrapper struct for an event.
type GuildBanAdd struct {
	*GuildBan
}

// GuildBanRemove is a wrapper struct for an event.
type GuildBanRemove struct {
	*GuildBan
}

// GuildMemberAdd is a wrapper struct for an event.
type GuildMemberAdd struct {
	*Member
}

// GuildMemberUpdate is a wrapper struct for an event.
type GuildMemberUpdate struct {
	*Member
}

// GuildMemberRemove is a wrapper struct for an event.
type GuildMemberRemove struct {
	*Member
}

// GuildRoleCreate is a wrapper struct for an event.
type GuildRoleCreate struct {
	*GuildRole
}

// GuildRoleUpdate is a wrapper struct for an event.
type GuildRoleUpdate struct {
	*GuildRole
}

// VoiceStateUpdate is a wrapper struct for an event.
type VoiceStateUpdate struct {
	*VoiceState
}

// UserUpdate is a wrapper struct for an event.
type UserUpdate struct {
	*UserUpdate
}

// UserSettingsUpdate is a map for an event.
type UserSettingsUpdate map[string]interface{}
