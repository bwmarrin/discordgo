package discordgo

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
