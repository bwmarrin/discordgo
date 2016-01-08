// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains code related to state tracking.  If enabled, state
// tracking will capture the initial READY packet and many other websocket
// events and maintain an in-memory state of of guilds, channels, users, and
// so forth.  This information can be accessed through the Session.State struct.

package discordgo

import "errors"

var nilError error = errors.New("State not instantiated, please use discordgo.New() or assign Session.State.")

// NewState creates an empty state.
func NewState() *State {
	return &State{
		Ready: Ready{
			PrivateChannels: []*Channel{},
			Guilds:          []*Guild{},
		},
	}
}

// OnReady takes a Ready event and updates all internal state.
func (s *State) OnReady(r *Ready) error {
	if s == nil {
		return nilError
	}

	s.Ready = *r
	return nil
}

// GuildAdd adds a guild to the current world state, or
// updates it if it already exists.
func (s *State) GuildAdd(guild *Guild) error {
	if s == nil {
		return nilError
	}

	for _, g := range s.Guilds {
		if g.ID == guild.ID {
			// This could be a little faster ;)
			for _, m := range guild.Members {
				s.MemberAdd(m)
			}
			for _, c := range guild.Channels {
				s.ChannelAdd(c)
			}
			return nil
		}
	}

	s.Guilds = append(s.Guilds, guild)
	return nil
}

// GuildRemove removes a guild from current world state.
func (s *State) GuildRemove(guild *Guild) error {
	if s == nil {
		return nilError
	}

	for i, g := range s.Guilds {
		if g.ID == guild.ID {
			s.Guilds = append(s.Guilds[:i], s.Guilds[i+1:]...)
			return nil
		}
	}

	return errors.New("Guild not found.")
}

// Guild gets a guild by ID.
// Useful for querying if @me is in a guild:
//     _, err := discordgo.Session.State.Guild(guildID)
//     isInGuild := err == nil
func (s *State) Guild(guildID string) (*Guild, error) {
	if s == nil {
		return nil, nilError
	}

	for _, g := range s.Guilds {
		if g.ID == guildID {
			return g, nil
		}
	}

	return nil, errors.New("Guild not found.")
}

// TODO: Consider moving Guild state update methods onto *Guild.

// MemberAdd adds a member to the current world state, or
// updates it if it already exists.
func (s *State) MemberAdd(member *Member) error {
	if s == nil {
		return nilError
	}

	guild, err := s.Guild(member.GuildID)
	if err != nil {
		return err
	}

	for i, m := range guild.Members {
		if m.User.ID == member.User.ID {
			guild.Members[i] = member
			return nil
		}
	}

	guild.Members = append(guild.Members, member)
	return nil
}

// MemberRemove removes a member from current world state.
func (s *State) MemberRemove(member *Member) error {
	if s == nil {
		return nilError
	}

	guild, err := s.Guild(member.GuildID)
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

// Member gets a member by ID from a guild.
func (s *State) Member(guildID, userID string) (*Member, error) {
	if s == nil {
		return nil, nilError
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, m := range guild.Members {
		if m.User.ID == userID {
			return m, nil
		}
	}

	return nil, errors.New("Member not found.")
}

// ChannelAdd adds a guild to the current world state, or
// updates it if it already exists.
// Channels may exist either as PrivateChannels or inside
// a guild.
func (s *State) ChannelAdd(channel *Channel) error {
	if s == nil {
		return nilError
	}

	if channel.IsPrivate {
		for i, c := range s.PrivateChannels {
			if c.ID == channel.ID {
				s.PrivateChannels[i] = channel
				return nil
			}
		}

		s.PrivateChannels = append(s.PrivateChannels, channel)
	} else {
		guild, err := s.Guild(channel.GuildID)
		if err != nil {
			return err
		}

		for i, c := range guild.Channels {
			if c.ID == channel.ID {
				guild.Channels[i] = channel
				return nil
			}
		}

		guild.Channels = append(guild.Channels, channel)
	}

	return nil
}

// ChannelRemove removes a channel from current world state.
func (s *State) ChannelRemove(channel *Channel) error {
	if s == nil {
		return nilError
	}

	if channel.IsPrivate {
		for i, c := range s.PrivateChannels {
			if c.ID == channel.ID {
				s.PrivateChannels = append(s.PrivateChannels[:i], s.PrivateChannels[i+1:]...)
				return nil
			}
		}
	} else {
		guild, err := s.Guild(channel.GuildID)
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

// GuildChannel gets a channel by ID from a guild.
func (s *State) GuildChannel(guildID, channelID string) (*Channel, error) {
	if s == nil {
		return nil, nilError
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, c := range guild.Channels {
		if c.ID == channelID {
			return c, nil
		}
	}

	return nil, errors.New("Channel not found.")
}

// PrivateChannel gets a private channel by ID.
func (s *State) PrivateChannel(channelID string) (*Channel, error) {
	if s == nil {
		return nil, nilError
	}

	for _, c := range s.PrivateChannels {
		if c.ID == channelID {
			return c, nil
		}
	}

	return nil, errors.New("Channel not found.")
}

// Channel gets a channel by ID, it will look in all guilds an private channels.
func (s *State) Channel(channelID string) (*Channel, error) {
	if s == nil {
		return nil, nilError
	}

	c, err := s.PrivateChannel(channelID)
	if err == nil {
		return c, nil
	}

	for _, g := range s.Guilds {
		c, err := s.GuildChannel(g.ID, channelID)
		if err == nil {
			return c, nil
		}
	}

	return nil, errors.New("Channel not found.")
}

// Emoji returns an emoji for a guild and emoji id.
func (s *State) Emoji(guildID, emojiID string) (*Emoji, error) {
	if s == nil {
		return nil, nilError
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, e := range guild.Emojis {
		if e.ID == emojiID {
			return e, nil
		}
	}

	return nil, errors.New("Emoji not found.")
}

// EmojiAdd adds an emoji to the current world state.
func (s *State) EmojiAdd(guildID string, emoji *Emoji) error {
	if s == nil {
		return nilError
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	for i, e := range guild.Emojis {
		if e.ID == emoji.ID {
			guild.Emojis[i] = emoji
			return nil
		}
	}

	guild.Emojis = append(guild.Emojis, emoji)
	return nil
}

// EmojisAdd adds multiple emojis to the world state.
func (s *State) EmojisAdd(guildID string, emojis []*Emoji) error {
	for _, e := range emojis {
		if err := s.EmojiAdd(guildID, e); err != nil {
			return err
		}
	}
	return nil
}
