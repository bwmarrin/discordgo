package discordgo

import "errors"

// NewState creates an empty state.
func NewState() *State {
	return &State{
		Ready: Ready{
			PrivateChannels: []Channel{},
			Guilds:          []Guild{},
		},
	}
}

// OnReady takes a Ready event and updates all internal state.
func (s *State) OnReady(r *Ready) {
	s.Ready = *r
}

// AddGuild adds a guild to the current world state, or
// updates it if it already exists.
func (s *State) AddGuild(guild *Guild) {
	for _, g := range s.Guilds {
		if g.ID == guild.ID {
			// This could be a little faster ;)
			for _, m := range guild.Members {
				s.AddMember(&m)
			}
			for _, c := range guild.Channels {
				s.AddChannel(&c)
			}
			return
		}
	}
	s.Guilds = append(s.Guilds, *guild)
}

// RemoveGuild removes a guild from current world state.
func (s *State) RemoveGuild(guild *Guild) error {
	for i, g := range s.Guilds {
		if g.ID == guild.ID {
			s.Guilds = append(s.Guilds[:i], s.Guilds[i+1:]...)
			return nil
		}
	}
	return errors.New("Guild not found.")
}

// GetGuildByID gets a guild by ID.
// Useful for querying if @me is in a guild:
//     _, err := discordgo.Session.State.GetGuildById(guildID)
//     isInGuild := err == nil
func (s *State) GetGuildByID(guildID string) (*Guild, error) {
	for _, g := range s.Guilds {
		if g.ID == guildID {
			return &g, nil
		}
	}
	return nil, errors.New("Guild not found.")
}

// TODO: Consider moving Guild state update methods onto *Guild.

// AddMember adds a member to the current world state, or
// updates it if it already exists.
func (s *State) AddMember(member *Member) error {
	guild, err := s.GetGuildByID(member.GuildID)
	if err != nil {
		return err
	}

	for i, m := range guild.Members {
		if m.User.ID == member.User.ID {
			guild.Members[i] = *member
			return nil
		}
	}

	guild.Members = append(guild.Members, *member)
	return nil
}

// RemoveMember removes a member from current world state.
func (s *State) RemoveMember(member *Member) error {
	guild, err := s.GetGuildByID(member.GuildID)
	if err != nil {
		return err
	}

	for i, m := range guild.Members {
		if m.User.ID == member.User.ID {
			guild.Members = append(guild.Members[:i], guild.Members[i+1:]...)
			return nil
		}
	}
	return errors.New("Member not found.")
}

// GetMemberByID gets a member by ID from a guild.
func (s *State) GetMemberByID(guildID string, userID string) (*Member, error) {
	guild, err := s.GetGuildByID(guildID)
	if err != nil {
		return nil, err
	}

	for _, m := range guild.Members {
		if m.User.ID == userID {
			return &m, nil
		}
	}
	return nil, errors.New("Member not found.")
}

// AddChannel adds a guild to the current world state, or
// updates it if it already exists.
// Channels may exist either as PrivateChannels or inside
// a guild.
func (s *State) AddChannel(channel *Channel) error {
	if channel.IsPrivate {
		for i, c := range s.PrivateChannels {
			if c.ID == channel.ID {
				s.PrivateChannels[i] = *channel
				return nil
			}
		}

		s.PrivateChannels = append(s.PrivateChannels, *channel)
	} else {
		guild, err := s.GetGuildByID(channel.GuildID)
		if err != nil {
			return err
		}

		for i, c := range guild.Channels {
			if c.ID == channel.ID {
				guild.Channels[i] = *channel
				return nil
			}
		}

		guild.Channels = append(guild.Channels, *channel)
	}
	return nil
}

// RemoveChannel removes a channel from current world state.
func (s *State) RemoveChannel(channel *Channel) error {
	if channel.IsPrivate {
		for i, c := range s.PrivateChannels {
			if c.ID == channel.ID {
				s.PrivateChannels = append(s.PrivateChannels[:i], s.PrivateChannels[i+1:]...)
				return nil
			}
		}
	} else {
		guild, err := s.GetGuildByID(channel.GuildID)
		if err != nil {
			return err
		}

		for i, c := range guild.Channels {
			if c.ID == channel.ID {
				guild.Channels = append(guild.Channels[:i], guild.Channels[i+1:]...)
				return nil
			}
		}
	}

	return errors.New("Channel not found.")
}

// GetGuildChannelById gets a channel by ID from a guild.
func (s *State) GetGuildChannelByID(guildID string, channelID string) (*Channel, error) {
	guild, err := s.GetGuildByID(guildID)
	if err != nil {
		return nil, err
	}

	for _, c := range guild.Channels {
		if c.ID == channelID {
			return &c, nil
		}
	}
	return nil, errors.New("Channel not found.")
}

// GetPrivateChannelByID gets a private channel by ID.
func (s *State) GetPrivateChannelByID(channelID string) (*Channel, error) {
	for _, c := range s.PrivateChannels {
		if c.ID == channelID {
			return &c, nil
		}
	}
	return nil, errors.New("Channel not found.")
}
