package discordgo

import (
	"time"
)

const InteractionDeadline = time.Second * 3

// ApplicationCommand is representing application's slash command.
type ApplicationCommand struct {
	ID            string                      `json:"id,omitempty"`
	ApplicationID string                      `json:"application_id,omitempty"`
	Name          string                      `json:"name,omitempty"`
	Description   string                      `json:"description,omitempty"`
	Options       []*ApplicationCommandOption `json:"options,omitempty"`
}

// ApplicationCommandOptionType is type of an slash-command's option.
type ApplicationCommandOptionType uint8

// Application command option types.
const (
	_ = ApplicationCommandOptionType(iota)

	ApplicationCommandOptionSubCommand
	ApplicationCommandOptionSubCommandGroup
	ApplicationCommandOptionString
	ApplicationCommandOptionInteger
	ApplicationCommandOptionBoolean
	ApplicationCommandOptionUser
	ApplicationCommandOptionChannel
	ApplicationCommandOptionRole
)

// ApplicationCommandOption is representing an option/subcommand/subcommands group.
type ApplicationCommandOption struct {
	Type        ApplicationCommandOptionType      `json:"type,omitempty"`
	Name        string                            `json:"name,omitempty"`
	Description string                            `json:"description,omitempty"`
	Default     bool                              `json:"default,omitempty"`
	Required    bool                              `json:"required,omitempty"`
	Choices     []*ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options     []*ApplicationCommandOption       `json:"options,omitempty"`
}

// ApplicationCommandOption is representing slash-command's option choice.
type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
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
	ID        string                            `json:"id"`
	Type      InteractionType                   `json:"type"`
	Data      ApplicationCommandInteractionData `json:"data"`
	GuildID   string                            `json:"guild_id"`
	ChannelID string                            `json:"channel_id"`
	Member    *Member                           `json:"member"`
	Token     string                            `json:"token"`
	Version   int                               `json:"version"`
}

// ApplicationCommandInteractionData is representing interaction data for application command.
type ApplicationCommandInteractionData struct {
	ID      string
	Name    string
	Options []*ApplicationCommandInteractionDataOption
}

// ApplicationCommandInteractionDataOption is representing an option of application's command.
type ApplicationCommandInteractionDataOption struct {
	Name string `json:"name,omitempty"`
	// Contains the value specified by InteractionType
	Value   interface{}                                `json:"value,omitempty"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
}

// InteractionResponseType is type of interaction response.
type InteractionResponseType uint8

// Interaction response types.
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

// InteractionResponse is representing response for interaction with application.
type InteractionResponse struct {
	Type InteractionResponseType                    `json:"type,omitempty"`
	Data *InteractionApplicationCommandResponseData `json:"data,omitempty"`
}

// InteractionApplicationCommandCallbackData is callback data for application command interaction.
type InteractionApplicationCommandResponseData struct {
	TTS             bool                    `json:"tts,omitempty"`
	Content         string                  `json:"content,omitempty"`
	Embeds          []*MessageEmbed         `json:"embeds,omitempty"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`

	Flags uint64 `json:"flags,omitempty"` // NOTE: Undocumented feature, be careful with it.
}
