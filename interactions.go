package discordgo

import (
	"time"
)

const InteractionDeadline = time.Second * 3

// ApplicationCommand is representing application's slash command.
type ApplicationCommand struct {
	ID            string                      `json:"id"`
	ApplicationID string                      `json:"application_id"`
	Name          string                      `json:"name"`
	Description   string                      `json:"description"`
	Options       []*ApplicationCommandOption `json:"options"`
}

// ApplicationCommandOptionType is type of an slash-command's option.
type ApplicationCommandOptionType uint8

const (
	_ = ApplicationCommandOptionType(iota)
	// ApplicationCommandOptionSubCommand is a type of option for sub-commands.
	ApplicationCommandOptionSubCommand
	// ApplicationCommandOptionSubCommand is a type of option for sub-commands groups.
	ApplicationCommandOptionSubCommandGroup
	// ApplicationCommandOptionSubCommand is a type of option for string options.
	ApplicationCommandOptionString
	// ApplicationCommandOptionSubCommand is a type of option for integer options.
	ApplicationCommandOptionInteger
	// ApplicationCommandOptionSubCommand is a type of option for boolean options.
	ApplicationCommandOptionBoolean
	// ApplicationCommandOptionSubCommand is a type of option when user id/mention is needed.
	ApplicationCommandOptionUser
	// ApplicationCommandOptionSubCommand is a type of option when channel id/mention is needed.
	ApplicationCommandOptionChannel
	// ApplicationCommandOptionSubCommand is a type of option when role id/mention is needed.
	ApplicationCommandOptionRole
)

// ApplicationCommandOption is representing an option/subcommand/subcommands group.
type ApplicationCommandOption struct {
	Type        ApplicationCommandOptionType      `json:"type"`
	Name        string                            `json:"name"`
	Description string                            `json:"description"`
	Default     bool                              `json:"default"`
	Required    bool                              `json:"required"`
	Choices     []*ApplicationCommandOptionChoice `json:"choices"`
	Options     []*ApplicationCommandOption       `json:"options"`
}

// ApplicationCommandOption is representing slash-command's option choice.
type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// InteractionType is representing interaction type.
type InteractionType uint8

const (
	_ = InteractionType(iota)
	// InteractionPing is type of interaction for ping.
	InteractionPing
	// InteractionApplicationCommand is type of interaction for application commands.
	InteractionApplicationCommand
)

// Interaction is representing interaction with application.
type Interaction struct {
	Id        string                            `json:"id"`
	Type      InteractionType                   `json:"type"`
	Data      ApplicationCommandInteractionData `json:"data"`
	GuildID   string                            `json:"guild_id"`
	ChannelID string                            `json:"channel_id"`
	Member    *Member                           `json:"member"`
	Token     string                            `json:"token"`
	Version   int                               `json:"version"`
}

// ApplicationCommandInteractionData is representing interaction data for application command
type ApplicationCommandInteractionData struct {
	ID      string
	Name    string
	Options []*ApplicationCommandInteractionDataOption
}

// ApplicationCommandInteractionDataOption is representing an option of application's command
type ApplicationCommandInteractionDataOption struct {
	Name string `json:"name"`
	// Contains the value specified by InteractionType
	Value   interface{}                                `json:"value"`
	Options []*ApplicationCommandInteractionDataOption `json:"options"`
}

// InteractionResponseType is type of interaction response
type InteractionResponseType uint8

const (
	_ = InteractionResponseType(iota)
	// InteractionResponsePong is an interaction response type when you need to just ACK a "Ping".
	InteractionResponsePong
	// InteractionResponsePong is an interaction response type when you need to ACK a command without sending a message, eating the user's input.
	InteractionResponseAcknowledge
	// InteractionResponsePong is an interaction response type when you need to respond with a message, eating the user's input.
	InteractionResponseChannelMessage
	// InteractionResponsePong is an interaction response type when you need to respond with a message, showing the user's input.
	InteractionResponseChannelMessageWithSource
	// InteractionResponsePong is an interaction response type when you need to ACK a command without sending a message, showing the user's input.
	InteractionResponseACKWithSource
)

// InteractionResponse is representing response for interaction with application
type InteractionResponse struct {
	Type InteractionResponseType                    `json:"type"`
	Data *InteractionApplicationCommandCallbackData `json:"data"`
}

// InteractionApplicationCommandCallbackData is callback data for application command interaction
type InteractionApplicationCommandCallbackData struct {
	TTS             bool                    `json:"tts"`
	Content         string                  `json:"content"`
	Embeds          []*MessageEmbed         `json:"embeds"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions"`
}
