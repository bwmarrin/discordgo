package discordgo

type Messageable interface {
	// Sends a message to the channel
	SendMessage(content string, embed *MessageEmbed, files []*File) (message *Message, err error)
	SendMessageComplex(data *MessageSend) (message *Message, err error)
	EditMessage(message *Message) (edited *Message, err error)
	EditMessageComplex(data *MessageEdit) (edited *Message, err error)

	// gets a single message by ID from the channel
	// ID : the ID of a Message
	FetchMessage(ID string) (message *Message, err error)

	// returns an array of Message structures for messages within
	// a given channel.
	// channelID : The ID of a Channel.
	// limit     : The number messages that can be returned. (max 100)
	// beforeID  : If provided all messages returned will be before given ID.
	// afterID   : If provided all messages returned will be after given ID.
	// aroundID  : If provided all messages returned will be around given ID.
	GetHistory(limit int, beforeID, afterID, aroundID string) (st []*Message, err error)
}

type IDGettable interface {
	GetID() string
}

type Mentionable interface {
	IDGettable
	Mention() string
}
