package discordgo

// TODO: Eventually everything here gets moved to a better place.

// A Message stores all data related to a specific Discord message.
type Message struct {
	ID              string       `json:"id"`
	Author          User         `json:"author"`
	Content         string       `json:"content"`
	Attachments     []Attachment `json:"attachments"`
	Tts             bool         `json:"tts"`
	Embeds          []Embed      `json:"embeds"`
	Timestamp       string       `json:"timestamp"`
	MentionEveryone bool         `json:"mention_everyone"`
	EditedTimestamp string       `json:"edited_timestamp"`
	Mentions        []User       `json:"mentions"`
	ChannelID       string       `json:"channel_id"`
}

// An Attachment stores data for message attachments.
type Attachment struct { //TODO figure this out
}

// An Embed stores data for message embeds.
type Embed struct { // TODO figure this out
}

// A VoiceRegion stores data for a specific voice region server.
type VoiceRegion struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"sample_hostname"`
	Port     int    `json:"sample_port"`
}

// A VoiceICE stores data for voice ICE servers.
type VoiceICE struct {
	TTL     string      `json:"ttl"`
	Servers []ICEServer `json:"servers"`
}

// A ICEServer stores data for a specific voice ICE server.
type ICEServer struct {
	URL        string `json:"url"`
	Username   string `json:"username"`
	Credential string `json:"credential"`
}

// A Invite stores all data related to a specific Discord Guild or Channel invite.
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
