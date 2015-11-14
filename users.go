package discordgo

type User struct {
	Id       int    `json:"id,string"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Avatar   string `json:"Avatar"`
	Verified bool   `json:"verified"`
	//Discriminator int    `json:"discriminator,string"` // TODO: See below
}

// Discriminator sometimes comes as a string
// and sometimes it comes as a int.  Weird.
// to avoid errors I've just commented it out
// but it doesn't seem to just kill the whole
// program.  Heartbeat is taken on READY even
// with error and the system continues to read
// it just doesn't seem able to handle this one
// field correctly.  Need to research this more.

type PrivateChannel struct {
	Id            int  `json:"id,string"`
	IsPrivate     bool `json:"is_private"`
	LastMessageId int  `json:"last_message_id,string"`
	Recipient     User `json:"recipient"`
} // merge with channel?

type Settings struct {
	RenderEmbeds          bool   `json:"render_embeds"`
	InlineEmbedMedia      bool   `json:"inline_embed_media"`
	EnableTtsCommand      bool   `json:"enable_tts_command"`
	MessageDisplayCompact bool   `json:"message_display_compact"`
	Locale                string `json:"locale"`
	ShowCurrentGame       bool   `json:"show_current_game"`
	Theme                 string `json:"theme"`
	//MutedChannels         []string `json:"muted_channels"` // TODO, see below
	MutedChannels []int `json:"muted_channels,string"` // TODO, see below
	//  MutedChannels []MutedChannel `json:"muted_channels"`
}

type MutedChannel struct {
	mc int `json:",string"`
}

// MutedChannels should be an array of ints...
// need to find a way to make that happen

// PM function to PM a user.
