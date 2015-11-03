package discordgo

type Server struct {
	Id               int      `json:"id,string"`
	Name             string   `json:"name"`
	Icon             string   `json:"icon"`
	Region           string   `json:"region"`
	Joined_at        string   `json:"joined_at"`
	Afk_timeout      int      `json:"afk_timeout"`
	Afk_channel_id   int      `json:"afk_channel_id"`
	Embed_channel_id int      `json:"embed_channel_id"`
	Embed_enabled    bool     `json:"embed_enabled"`
	Owner_id         int      `json:"owner_id,string"`
	Roles            []Role   `json:"roles"`
	Session          *Session // I got to be doing it wrong here.
}

// Channels returns an array of Channel structures for channels within
// this Server
func (s *Server) Channels() (c []Channel, err error) {
	c, err = Channels(s.Session, s.Id)
	return
}
