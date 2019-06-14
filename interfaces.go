package discordgo

type Messageable interface {
	SendMessage(content string, embed MessageEmbed, file File) (err error)
	EditMessage(content string, embed MessageEmbed, file File) (err error)
	FetchMessage(id string) (message Message, err error)
}

type IDGettable interface {
	GetID() string
}

type Mentionable interface {
	IDGettable
	Mention() string
}
