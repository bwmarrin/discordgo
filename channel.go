package discordgo

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
	Session               *Session
}

/*
func (c *Channel) Messages() (messages []Message) {
}

func (c *Channel) SendMessage() (content string) {
}
*/
