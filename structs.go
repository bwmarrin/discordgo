package discordgo

type Server struct {
	Afk_timeout int
	Joined_at   string
	// Afk_channel_id int `json:",string"`
	Id   int `json:",string"`
	Icon string
	Name string
	//	Roles          []Role
	Region string
	//Embed_channel_id int `json:",string"`
	//	Embed_channel_id string
	//	Embed_enabled    bool
	Owner_id int `json:",string"`
}

type Role struct {
	Permissions int
	Id          int `json:",string"`
	Name        string
}

type Channel struct {
	Guild_id        int `json:",string"`
	Id              int `json:",string"`
	Name            string
	Last_message_id string
	Is_private      string

	//	Permission_overwrites string
	//	Position              int `json:",string"`
	//	Type                  string
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

type User struct {
}
