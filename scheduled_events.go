package discordgo

import (
	"net/url"
	"strconv"
)

// GuildScheduledEvents returns an array of GuildScheduledEvent for a guild
// guildID        : The ID of a Guild
// withUserCounts : Whether to include the user count in the response
func (s *Session) GuildScheduledEvents(guildID string, withUserCount bool) (st []*GuildScheduledEvent, err error) {
	uri := EndpointGuildScheduledEvents(guildID)

	queryParams := url.Values{}
	if withUserCount {
		queryParams.Set("with_user_count", "true")
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, EndpointGuildScheduledEvents(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEvent returns a specific GuildScheduledEvent for a guild
// guildID        : The ID of a Guild
// eventID        : The ID of the event
// withUserCounts : Whether to include the user count in the response
func (s *Session) GuildScheduledEvent(guildID, eventID string, withUserCount bool) (st *GuildScheduledEvent, err error) {
	uri := EndpointGuildScheduledEvent(guildID, eventID)

	queryParams := url.Values{}
	if withUserCount {
		queryParams.Set("with_user_count", "true")
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, EndpointGuildScheduledEvent(guildID, eventID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildScheduledEventCreate creates a GuildScheduledEvent for a guild and returns it
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

// GuildScheduledEventUpdate updates a specific event for a guild and returns it.
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

// GuildScheduledEventDelete deletes a specific GuildScheduledEvent for a guild
// guildID   : The ID of a Guild
// eventID   : The ID of the event
func (s *Session) GuildScheduledEventDelete(guildID, eventID string) (err error) {
	_, err = s.RequestWithBucketID("DELETE", EndpointGuildScheduledEvent(guildID, eventID), nil, EndpointGuildScheduledEvent(guildID, eventID))
	return
}

// GuildScheduledEventUsers returns an array of GuildScheduledEventUser for a guild
// guildID    : The ID of a Guild
// eventID    : The ID of the event
// limit      : The maximum number of users to return (Max 100)
// withMember : Whether to include the member object in the response
// beforeID   : If provided all events users entries returned will be before the given ID.
// afterID    : If provided all events users entries returned will be after the given ID.
func (s *Session) GuildScheduledEventUsers(guildID, eventID string, limit int, withMember bool, beforeID, afterID string) (st []*GuildScheduledEventUser, err error) {
	uri := EndpointGuildScheduledEventUsers(guildID, eventID)

	queryParams := url.Values{}
	if withMember {
		queryParams.Set("with_member", "true")
	}
	if limit > 0 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}
	if beforeID != "" {
		queryParams.Set("before", beforeID)
	}
	if afterID != "" {
		queryParams.Set("after", afterID)
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, EndpointGuildScheduledEventUsers(guildID, eventID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}
