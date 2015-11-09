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

type Attachment struct {
}

type Embed struct {
}
