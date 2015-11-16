package discordgo

// TODO: Eventually everything here gets moved to a better place.

type Message struct {
	Id              int          `json:"id,string"`
	Author          User         `json:"author"`
	Content         string       `json:"content"`
	Attachments     []Attachment `json:"attachments"`
	Tts             bool         `json:"tts"`
	Embeds          []Embed      `json:"embeds"`
	Timestamp       string       `json:"timestamp"`
	MentionEveryone bool         `json:"mention_everyone"`
	EditedTimestamp string       `json:"edited_timestamp"`
	Mentions        []User       `json:"mentions"`
	ChannelId       int          `json:"channel_id,string"`
}

type Attachment struct { //TODO figure this out
}

type Embed struct { // TODO figure this out
}

type VoiceRegion struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	SampleHostname string `json:"sample_hostname"`
	SamplePort     int    `json:"sample_port"`
}

type VoiceIce struct {
	Ttl     int         `json:"ttl,string"`
	Servers []IceServer `json:"servers"`
}

type IceServer struct {
	Url        string `json:"url"`
	Username   string `json:"username"`
	Credential string `json:"credential"`
}

type Invite struct {
	MaxAge    int     `json:"max_age"`
	Code      string  `json:"code"`
	Guild     Guild   `json:"guild"`
	Revoked   bool    `json:"revoked"`
	CreatedAt string  `json:"created_at"` // TODO make timestamp
	Temporary bool    `json:"temporary"`
	Uses      int     `json:"uses"`
	MaxUses   int     `json:"max_uses"`
	Inviter   User    `json:"inviter"`
	XkcdPass  bool    `json:"xkcdpass"`
	Channel   Channel `json:"channel"`
}
