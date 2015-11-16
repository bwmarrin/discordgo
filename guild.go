package discordgo

type Guild struct {
	Id             string       `json:"id"`
	Name           string       `json:"name"`
	Icon           string       `json:"icon"`
	Region         string       `json:"region"`
	AfkTimeout     int          `json:"afk_timeout"`
	AfkChannelId   string       `json:"afk_channel_id"`
	EmbedChannelId string       `json:"embed_channel_id"`
	EmbedEnabled   bool         `json:"embed_enabled"`
	OwnerId        string       `json:"owner_id"`
	Large          bool         `json:"large"`     // ??
	JoinedAt       string       `json:"joined_at"` // make this a timestamp
	Roles          []Role       `json:"roles"`
	Members        []Member     `json:"members"`
	Presences      []Presence   `json:"presences"`
	Channels       []Channel    `json:"channels"`
	VoiceStates    []VoiceState `json:"voice_states"`
}

type Role struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Managed     bool   `json:"managed"`
	Color       int    `json:"color"`
	Hoist       bool   `json:"hoist"`
	Position    int    `json:"position"`
	Permissions int    `json:"permissions"`
}

type VoiceState struct {
	UserId    string `json:"user_id"`
	Suppress  bool   `json:"suppress"`
	SessionId string `json:"session_id"`
	SelfMute  bool   `json:"self_mute"`
	SelfDeaf  bool   `json:"self_deaf"`
	Mute      bool   `json:"mute"`
	Deaf      bool   `json:"deaf"`
	ChannelId string `json:"channel_id"`
}

type Presence struct {
	User   User   `json:"user"`
	Status string `json:"status"`
	GameId int    `json:"game_id"`
}

// TODO: Member vs User?
type Member struct {
	GuildId  string   `json:"guild_id"`
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
