package main

import (
	Discord "github.com/bwmarrin/discordgo"
)

// Registers all event handlers
func RegisterHandlers() {

	Session = Discord.Session{
		OnEvent:                   OnEvent,
		OnReady:                   OnReady,
		OnTypingStart:             OnTypingStart,
		OnMessageCreate:           OnMessageCreate,
		OnMessageUpdate:           OnMessageUpdate,
		OnMessageDelete:           OnMessageDelete,
		OnMessageAck:              OnMessageAck,
		OnVoiceStateUpdate:        OnVoiceStateUpdate,
		OnPresenceUpdate:          OnPresenceUpdate,
		OnChannelCreate:           OnChannelCreate,
		OnChannelUpdate:           OnChannelUpdate,
		OnGuildCreate:             OnGuildCreate,
		OnGuildUpdate:             OnGuildUpdate,
		OnGuildDelete:             OnGuildDelete,
		OnGuildRoleCreate:         OnGuildRoleCreate,
		OnGuildRoleUpdate:         OnGuildRoleUpdate,
		OnGuildRoleDelete:         OnGuildRoleDelete,
		OnGuildMemberAdd:          OnGuildMemberAdd,
		OnGuildMemberUpdate:       OnGuildMemberUpdate,
		OnGuildMemberRemove:       OnGuildMemberRemove,
		OnGuildIntegrationsUpdate: OnGuildIntegrationsUpdate,
	}

}

// OnEvent is called for unknown events or unhandled events.  It provides
// a generic interface to handle them.
func OnEvent(s *Discord.Session, e Discord.Event) {
	// Add code here to handle this event.
}

// OnReady is called when Discordgo receives a READY event
// This event must be handled and must contain the Heartbeat call.
func OnReady(s *Discord.Session, st Discord.Ready) {

	// start the Heartbeat
	go s.Heartbeat(st.HeartbeatInterval)

	// Add code here to handle this event.
}

func OnTypingStart(s *Discord.Session, st Discord.TypingStart) {
	// Add code here to handle this event.
}

func OnPresenceUpdate(s *Discord.Session, st Discord.PresenceUpdate) {
	// Add code here to handle this event.
}

func OnMessageCreate(s *Discord.Session, m Discord.Message) {
	// Add code here to handle this event.
}

func OnMessageUpdate(s *Discord.Session, m Discord.Message) {
	// Add code here to handle this event.
}

func OnMessageAck(s *Discord.Session, st Discord.MessageAck) {
	// Add code here to handle this event.
}

func OnMessageDelete(s *Discord.Session, st Discord.MessageDelete) {
	// Add code here to handle this event.
}

func OnVoiceStateUpdate(s *Discord.Session, st Discord.VoiceState) {
	// Add code here to handle this event.
}

func OnChannelCreate(s *Discord.Session, st Discord.Channel) {
	// Add code here to handle this event.
}

func OnChannelUpdate(s *Discord.Session, st Discord.Channel) {
	// Add code here to handle this event.
}

func OnGuildCreate(s *Discord.Session, st Discord.Guild) {
	// Add code here to handle this event.
}

func OnGuildUpdate(s *Discord.Session, st Discord.Guild) {
	// Add code here to handle this event.
}
func OnGuildDelete(s *Discord.Session, st Discord.Guild) {
	// Add code here to handle this event.
}

func OnGuildRoleCreate(s *Discord.Session, st Discord.GuildRole) {
	// Add code here to handle this event.
}
func OnGuildRoleUpdate(s *Discord.Session, st Discord.GuildRole) {
	// Add code here to handle this event.
}
func OnGuildRoleDelete(s *Discord.Session, st Discord.GuildRoleDelete) {
	// Add code here to handle this event.
}
func OnGuildMemberAdd(s *Discord.Session, st Discord.Member) {
	// Add code here to handle this event.
}

func OnGuildMemberUpdate(s *Discord.Session, st Discord.Member) {
	// Add code here to handle this event.
}

func OnGuildMemberRemove(s *Discord.Session, st Discord.Member) {
	// Add code here to handle this event.
}

func OnGuildIntegrationsUpdate(s *Discord.Session, st Discord.GuildIntegrationsUpdate) {
	// Add code here to handle this event.
}
