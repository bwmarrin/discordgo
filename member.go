package discordgo

import "sort"

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

func (m *Member) String() string {
	return m.User.String()
}

func (m *Member) GetID() string {
	return m.User.ID
}

// Mention creates a member mention
func (m *Member) Mention() string {
	return m.User.Mention()
}

// AvatarURL returns a URL to the user's avatar.
//    size:    The size of the user's avatar as a power of two
//             if size is an empty string, no size parameter will
//             be added to the URL.
func (m *Member) AvatarURL(size string) string {
	return m.User.AvatarURL(size)
}

func (m *Member) GetDisplayName() string {
	if m.Nick != "" {
		return m.Nick
	} else {
		return m.User.Username
	}
}

func (m *Member) GetGuild() (g *Guild, err error) {
	return m.User.Session.State.Guild(m.GuildID)
}

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

func (m *Member) GetColour() (color int, err error) {
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

func (m *Member) GetTopRole() (role *Role, err error) {
	roles, err := m.GetRoles()
	if err != nil {
		return
	}

	role = roles[0]
	return
}

func (m *Member) Kick(reason string) (err error) {
	g, err := m.GetGuild()
	if err != nil {
		return
	}

	return g.Kick(m.User, reason)
}

func (m *Member) Ban(reason string, days int) (err error) {
	g, err := m.GetGuild()
	if err != nil {
		return
	}

	return g.Ban(m.User, reason, days)
}

func (m *Member) EditRoles(roles Roles) (err error) {
	var roleIDs []string
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}
	return m.User.Session.GuildMemberEdit(m.GuildID, m.User.ID, roleIDs)
}

func (m *Member) EditNickname(nick string) (err error) {
	return m.User.Session.GuildMemberNickname(m.GuildID, m.User.ID, nick)
}

func (m *Member) MoveTo(channel *Channel) (err error) {
	return m.User.Session.GuildMemberMove(m.GuildID, m.User.ID, channel.ID)
}

func (m *Member) AddRole(role *Role) (err error) {
	return m.User.Session.GuildMemberRoleAdd(m.GuildID, m.User.ID, role.ID)
}

func (m *Member) RemoveRole(role *Role) (err error) {
	return m.User.Session.GuildMemberRoleRemove(m.GuildID, m.User.ID, role.ID)
}
