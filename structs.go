package discordgo

type User struct {
	Id            int    `json:"id,string"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	Avatar        string `json:"Avatar"`
	Verified      bool   `json:"verified"`
	Discriminator string `json:"discriminator"`
}

type Server struct {
	Id               int    `json:"id,string"`
	Name             string `json:"name"`
	Icon             string `json:"icon"`
	Region           string `json:"region"`
	Joined_at        string `json:"joined_at"`
	Afk_timeout      int    `json:"afk_timeout"`
	Afk_channel_id   int    `json:"afk_channel_id"`
	Embed_channel_id int    `json:"embed_channel_id"`
	Embed_enabled    bool   `json:"embed_enabled"`
	Owner_id         int    `json:"owner_id,string"`
	Roles            []Role `json:"roles"`
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

type Channel struct {
	Server_id             int    `json:"guild_id,string,omitempty"`
	Id                    int    `json:"id,string"`
	Name                  string `json:"name"`
	Topic                 string `json:"topic"`
	Position              int    `json:"position"`
	Last_message_id       int    `json:"last_message_id,string"`
	Type                  string `json:"type"`
	Is_private            bool   `json:"is_private"`
	Permission_overwrites string `json:"-"` // ignored for now
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
