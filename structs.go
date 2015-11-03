package discordgo

type User struct {
	Id            int    `json:"id,string"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	Avatar        string `json:"Avatar"`
	Verified      bool   `json:"verified"`
	Discriminator string `json:"discriminator"`
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
	Attachments      []Attachment
	Tts              bool
	Embeds           []Embed
	Timestamp        string
	Mention_everyone bool
	Id               int `json:",string"`
	Edited_timestamp string
	Author           *Author
	Content          string
	Channel_id       int `json:",string"`
	Mentions         []Mention
}

type Mention struct {
}

type Attachment struct {
}

type Embed struct {
}

type Author struct {
	Username      string
	Discriminator int `json:",string"`
	Id            int `json:",string"`
	Avatar        string
}
