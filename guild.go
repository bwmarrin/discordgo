package discordgo

type Guild struct {
	Id               int          `json:"id,string"`
	Name             string       `json:"name"`
	Icon             string       `json:"icon"`
	Region           string       `json:"region"`
	Joined_at        string       `json:"joined_at"` // make time stamp
	Afk_timeout      int          `json:"afk_timeout"`
	Afk_channel_id   int          `json:"afk_channel_id,string"`
	Embed_channel_id int          `json:"embed_channel_id"`
	Embed_enabled    bool         `json:"embed_enabled"`
	Owner_id         int          `json:"owner_id,string"`
	Large            bool         `json:"large"`     // ??
	JoinedAt         string       `json:"joined_at"` // make this a timestamp
	Roles            []Role       `json:"roles"`
	Members          []Member     `json:"members"`
	Presences        []Presence   `json:"presences"`
	Channels         []Channel    `json:"channels"`
	VoiceStates      []VoiceState `json:"voice_states"`
}

type Role struct {
	Id          int    `json:"id,string"`
	Name        string `json:"name"`
	Managed     bool   `json:"managed"`
	Color       int    `json:"color"`
	Hoist       bool   `json:"hoist"`
	Position    int    `json:"position"`
	Permissions int    `json:"permissions"`
}

type VoiceState struct {
	UserId    int    `json:"user_id,string"`
	Suppress  bool   `json:"suppress"`
	SessionId string `json:"session_id"`
	SelfMute  bool   `json:"self_mute"`
	SelfDeaf  bool   `json:"self_deaf"`
	Mute      bool   `json:"mute"`
	Deaf      bool   `json:"deaf"`
	ChannelId int    `json:"channel_id,string"`
}

type Presence struct {
	User   User   `json:"user"`
	Status string `json:"status"`
	GameId int    `json:"game_id"`
}

// TODO: Member vs User?
type Member struct {
	GuildId  int      `json:"guild_id,string"`
	JoinedAt string   `json:"joined_at"`
	Deaf     bool     `json:"deaf"`
	mute     bool     `json:"mute"`
	User     User     `json:"user"`
	Roles    []string `json:"roles"` // TODO: See below
}

//Roles   []string `json:"roles"` // TODO: Should be ints, see below
// Above "Roles" should be an array of ints
// TODO: Figure out how to make it be one.
/*
	{
		"roles": [
			"89544728336416768",
			"110429733396676608"
		],
	}
*/

// Channels returns an array of Channel structures for channels within
// this Server
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
