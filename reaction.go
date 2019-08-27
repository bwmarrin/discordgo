package discordgo

// MessageReaction stores the data for a message reaction.
type MessageReaction struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	Emoji     Emoji  `json:"emoji"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`

	Session *Session `json:"-"`
}

func (r *MessageReaction) Remove() error {
	return r.Session.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
}
