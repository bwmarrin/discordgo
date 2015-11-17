package discordgo

// A Guild holds all data related to a specific Discord Guild.  Guilds are also
// sometimes referred to as Servers in the Discord client.
type Guild struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Icon           string       `json:"icon"`
	Region         string       `json:"region"`
	AfkTimeout     int          `json:"afk_timeout"`
	AfkChannelID   string       `json:"afk_channel_id"`
	EmbedChannelID string       `json:"embed_channel_id"`
	EmbedEnabled   bool         `json:"embed_enabled"`
	OwnerID        string       `json:"owner_id"`
	Large          bool         `json:"large"`     // ??
	JoinedAt       string       `json:"joined_at"` // make this a timestamp
	Roles          []Role       `json:"roles"`
	Members        []Member     `json:"members"`
	Presences      []Presence   `json:"presences"`
	Channels       []Channel    `json:"channels"`
	VoiceStates    []VoiceState `json:"voice_states"`
}

// A Role stores information about Discord guild member roles.
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Managed     bool   `json:"managed"`
	Color       int    `json:"color"`
	Hoist       bool   `json:"hoist"`
	Position    int    `json:"position"`
	Permissions int    `json:"permissions"`
}

// A VoiceState stores the voice states of Guilds
type VoiceState struct {
	UserID    string `json:"user_id"`
	Suppress  bool   `json:"suppress"`
	SessionID string `json:"session_id"`
	SelfMute  bool   `json:"self_mute"`
	SelfDeaf  bool   `json:"self_deaf"`
	Mute      bool   `json:"mute"`
	Deaf      bool   `json:"deaf"`
	ChannelID string `json:"channel_id"`
}

// A Presence stores the online, offline, or idle and game status of Guild members.
type Presence struct {
	User   User   `json:"user"`
	Status string `json:"status"`
	GameID int    `json:"game_id"`
}

// A Member stores user information for Guild members.
type Member struct {
	GuildID  string   `json:"guild_id"`
	JoinedAt string   `json:"joined_at"`
	Deaf     bool     `json:"deaf"`
	Mute     bool     `json:"mute"`
	User     User     `json:"user"`
	Roles    []string `json:"roles"`
}

/*
TODO: How to name these? If we make a variable to store
channels from READY packet, etc.  We can't have a Channel
func?  And which is better.  Channels func gets live up
to date data on each call.. so, there is some benefit there.

Maybe it should have both, but make the Channels check and
pull new data based on a cache time?

func (s *Server) Channels() (c []Channel, err error) {
	c, err = Channels(s.Session, s.Id)
	return
}
*/
/*

func (s *Server) Members() (m []Users, err error) {
}
*/
