package discordgo

// A Channel holds all data related to an individual Discord channel.
type Channel struct {
	ID                   string                `json:"id"`
	GuildID              string                `json:"guild_id"`
	Name                 string                `json:"name"`
	Topic                string                `json:"topic"`
	Position             int                   `json:"position"`
	Type                 string                `json:"type"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	IsPrivate            bool                  `json:"is_private"`
	LastMessageID        string                `json:"last_message_id"`
	Recipient            User                  `json:"recipient"`
}

// A PermissionOverwrite holds permission overwrite data for a Channel
type PermissionOverwrite struct {
	ID    string `json:"id"`
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
