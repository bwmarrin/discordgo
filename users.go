package discordgo

// A User stores all data for an individual Discord user.
type User struct {
	ID       string `json:"id"`
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

// A PrivateChannel stores all data for a specific user private channel.
type PrivateChannel struct {
	ID            string `json:"id"`
	IsPrivate     bool   `json:"is_private"`
	LastMessageID string `json:"last_message_id"`
	Recipient     User   `json:"recipient"`
} // merge with channel?

// A Settings stores data for a specific users Discord client settings.
type Settings struct {
	RenderEmbeds          bool     `json:"render_embeds"`
	InlineEmbedMedia      bool     `json:"inline_embed_media"`
	EnableTtsCommand      bool     `json:"enable_tts_command"`
	MessageDisplayCompact bool     `json:"message_display_compact"`
	Locale                string   `json:"locale"`
	ShowCurrentGame       bool     `json:"show_current_game"`
	Theme                 string   `json:"theme"`
	MutedChannels         []string `json:"muted_channels"`
}
