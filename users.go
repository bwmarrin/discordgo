package discordgo

type User struct {
	Id            int    `json:"id,string"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	Avatar        string `json:"Avatar"`
	Verified      bool   `json:"verified"`
	Discriminator string `json:"discriminator"`
}

type PrivateChannel struct {
	Id            int  `json:"id,string"`
	IsPrivate     bool `json:"is_private"`
	LastMessageId int  `json:"last_message_id,string"`
	Recipient     User `json:"recipient"`
}

// PM function to PM a user.
