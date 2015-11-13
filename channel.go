package discordgo

type Channel struct {
	GuildId              int                   `json:"guild_id,string,omitempty"`
	Id                   int                   `json:"id,string"`
	Name                 string                `json:"name"`
	Topic                string                `json:"topic"`
	Position             int                   `json:"position"`
	Type                 string                `json:"type"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	IsPrivate            bool                  `json:"is_private"`
	LastMessageId        int                   `json:"last_message_id,string"`
	Recipient            User                  `json:"recipient"`
	Session              *Session
}

type PermissionOverwrite struct {
	Type  string `json:"type"`
	Id    int    `json:"id,string"`
	Deny  int    `json:"deny"`
	Allow int    `json:"allow"`
}

/*
func (c *Channel) Messages() (messages []Message) {
}

func (c *Channel) SendMessage() (content string) {
}
*/
