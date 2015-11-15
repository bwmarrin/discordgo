package discordgo

type Channel struct {
	Id                   string                `json:"id"`
	GuildId              string                `json:"guild_idomitempty"`
	Name                 string                `json:"name"`
	Topic                string                `json:"topic"`
	Position             int                   `json:"position"`
	Type                 string                `json:"type"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	IsPrivate            bool                  `json:"is_private"`
	LastMessageId        string                `json:"last_message_id"`
	Recipient            User                  `json:"recipient"`
	Session              *Session
}

type PermissionOverwrite struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Deny  int    `json:"deny"`
	Allow int    `json:"allow"`
}

/*
func (c *Channel) Messages() (messages []Message) {
}

func (c *Channel) SendMessage() (content string) {
}
*/
