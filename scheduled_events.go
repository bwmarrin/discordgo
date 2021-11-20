package discordgo

// GuildScheduledEvents returns an array of GuildScheduledEvent for a guild
// guildID   : The ID of a Guild
func (s *Session) GuildScheduledEvents(guildID string) (st []*GuildScheduledEvent, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointGuildScheduledEvents(guildID), nil, EndpointGuildScheduledEvents(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEvents returns an array of GuildScheduledEvent for a guild
// guildID   : The ID of a Guild
// eventID   : The ID of the event
func (s *Session) GuildScheduledEvent(guildID, eventID string) (st *GuildScheduledEvent, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointGuildScheduledEvent(guildID, eventID), nil, EndpointGuildScheduledEvent(guildID, eventID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEvents returns an array of GuildScheduledEvent for a guild
// guildID   : The ID of a Guild
// eventID   : The ID of the event
func (s *Session) GuildScheduledEventCreate(guildID string, event *GuildScheduledEvent) (st *GuildScheduledEvent, err error) {
	body, err := s.RequestWithBucketID("POST", EndpointGuildScheduledEvents(guildID), event, EndpointGuildScheduledEvents(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEvents returns an array of GuildScheduledEvent for a guild
// guildID   : The ID of a Guild
// eventID   : The ID of the event
func (s *Session) GuildScheduledEventUpdate(guildID, eventID string, event *GuildScheduledEvent) (st *GuildScheduledEvent, err error) {
	body, err := s.RequestWithBucketID("PATCH", EndpointGuildScheduledEvent(guildID, eventID), event, EndpointGuildScheduledEvent(guildID, eventID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEvents returns an array of GuildScheduledEvent for a guild
// guildID   : The ID of a Guild
// eventID   : The ID of the event
func (s *Session) GuildScheduledEventDelete(guildID, eventID string) (err error) {
	_, err = s.RequestWithBucketID("DELETE", EndpointGuildScheduledEvent(guildID, eventID), nil, EndpointGuildScheduledEvent(guildID, eventID))
	return
}

// GuildScheduledEvents returns an array of GuildScheduledEvent for a guild
// guildID   : The ID of a Guild
// eventID   : The ID of the event
func (s *Session) GuildScheduledEventUsers(guildID, eventID string) (st []*GuildScheduledEventUser, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointGuildScheduledEventUsers(guildID, eventID), nil, EndpointGuildScheduledEventUsers(guildID, eventID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}
