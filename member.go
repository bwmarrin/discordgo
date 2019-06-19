package discordgo

import (
	"sort"
	"time"
)

// A Member stores user information for Guild members. A guild
// member represents a certain user's presence in a guild.
type Member struct {
	// The guild ID on which the member exists.
	GuildID string `json:"guild_id"`

	// The time at which the member joined the guild, in ISO8601.
	JoinedAt Timestamp `json:"joined_at"`

	// The nickname of the member, if they have one.
	Nick string `json:"nick"`

	// Whether the member is deafened at a guild level.
	Deaf bool `json:"deaf"`

	// Whether the member is muted at a guild level.
	Mute bool `json:"mute"`

	// The underlying user on which the member is based.
	User *User `json:"user"`

	// A list of IDs of the roles which are possessed by the member.
	Roles []string `json:"roles"`
}

// String returns a unique identifier of the form username#discriminator
func (m *Member) String() string {
	return m.User.String()
}

// GetID returns the members ID
func (m *Member) GetID() string {
	return m.User.ID
}

// CreatedAt returns the members creation time in UTC
func (m *Member) CreatedAt() (creation time.Time, err error) {
	return m.User.CreatedAt()
}

// Mention creates a member mention
func (m *Member) Mention() string {
	if m.Nick != "" {
		return "<@!" + m.User.ID + ">"
	}
	return m.User.Mention()
}

// IsMentionedIn checks if the member is mentioned in the given message
// message      : message to check for mentions
func (m *Member) IsMentionedIn(message *Message) bool {
	if m.User.IsMentionedIn(message) {
		return true
	}

	roles, err := m.GetRoles()
	if err != nil {
		return false
	}
	roles = Roles(roles)

	for _, roleID := range message.MentionRoles {
		if roles.ContainsID(roleID) {
			return true
		}
	}

	return false
}

// AvatarURL returns a URL to the user's avatar.
//    size:    The size of the user's avatar as a power of two
//             if size is an empty string, no size parameter will
//             be added to the URL.
func (m *Member) AvatarURL(size string) string {
	return m.User.AvatarURL(size)
}

// GetDisplayName returns the members nick if one has been set and else their username
func (m *Member) GetDisplayName() string {
	if m.Nick != "" {
		return m.Nick
	}
	return m.User.Username
}

// GetGuild returns the guild object where the Member belongs to
func (m *Member) GetGuild() (g *Guild, err error) {
	return m.User.Session.State.Guild(m.GuildID)
}

// GetRoles returns a slice with all roles the Member has, sorted from highest to lowest
func (m *Member) GetRoles() (roles Roles, err error) {
	g, err := m.GetGuild()
	if err != nil {
		return
	}

	for _, roleID := range m.Roles {
		r, errGR := g.GetRole(roleID)
		if errGR != nil {
			err = errGR
			return
		}
		roles = append(roles, r)
	}
	sort.Sort(roles)
	return
}

// GetColor returns the hex code of the members color as displayed in the server
func (m *Member) GetColor() (color int, err error) {
	roles, err := m.GetRoles()
	if err != nil {
		return
	}

	for _, role := range roles {
		if role.Color != 0 {
			return role.Color, nil
		}
	}

	return
}

// GetTopRole returns the members highest role
func (m *Member) GetTopRole() (role *Role, err error) {
	roles, err := m.GetRoles()
	if err != nil {
		return
	}

	role = roles[0]
	return
}

// Kick kicks the member from their guild
// reason   : reason for the kick
func (m *Member) Kick(reason string) (err error) {
	g, err := m.GetGuild()
	if err != nil {
		return
	}

	return g.Kick(m.User, reason)
}

// Ban bans the member from their guild
// reason     : reason for the ban as it will be displayed in the auditlog
// days       : days of messages to delete
func (m *Member) Ban(reason string, days int) (err error) {
	g, err := m.GetGuild()
	if err != nil {
		return
	}

	return g.Ban(m.User, reason, days)
}

// EditRoles replaces all roles of the user with the provided slice of roles
// roles     : a slice of Role objects
func (m *Member) EditRoles(roles Roles) (err error) {
	var roleIDs []string
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}
	return m.User.Session.GuildMemberEdit(m.GuildID, m.User.ID, roleIDs)
}

// EditNickname sets the nickname of the member
// nick      : the new nickname the member will have
func (m *Member) EditNickname(nick string) (err error) {
	return m.User.Session.GuildMemberNickname(m.GuildID, m.User.ID, nick)
}

// MoveTo moves the member to a voice channel
// channel   : voice channel to move the user to
func (m *Member) MoveTo(channel *Channel) (err error) {
	if channel.Type != ChannelTypeGuildVoice {
		return ErrNotAVoiceChannel
	}
	return m.User.Session.GuildMemberMove(m.GuildID, m.User.ID, channel.ID)
}

// AddRole adds a role to the member
// role     : role to add
func (m *Member) AddRole(role *Role) (err error) {
	return m.User.Session.GuildMemberRoleAdd(m.GuildID, m.User.ID, role.ID)
}

// RemoveRole removes a role from the member
// role     : role to remove
func (m *Member) RemoveRole(role *Role) (err error) {
	return m.User.Session.GuildMemberRoleRemove(m.GuildID, m.User.ID, role.ID)
}
