package discordgo

type User struct {
	Id            int    `json:"id,string"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	Avatar        string `json:"Avatar"`
	Verified      bool   `json:"verified"`
	Discriminator string `json:"discriminator"`
}

type Member struct {
	JoinedAt string `json:"joined_at"`
	Deaf     bool   `json:"deaf"`
	mute     bool   `json:"mute"`
	User     User   `json:"user"`
	Roles    []Role `json:"roles"`
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
