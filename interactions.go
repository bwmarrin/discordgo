package discordgo

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// InteractionDeadline is a deadline for responding to an interaction, if you haven't responded in the time, you won't be able to respond later.
const InteractionDeadline = time.Second * 3

// ApplicationCommand is representing application's slash command.
type ApplicationCommand struct {
	ID            string                      `json:"id"`
	ApplicationID string                      `json:"application_id,omitempty"`
	Name          string                      `json:"name"`
	Description   string                      `json:"description,omitempty"`
	Options       []*ApplicationCommandOption `json:"options"`
}

// ApplicationCommandOptionType is type of an slash-command's option.
type ApplicationCommandOptionType uint8

// Application command option types.
const (
	ApplicationCommandOptionSubCommand = ApplicationCommandOptionType(iota + 1)
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
	Type        ApplicationCommandOptionType `json:"type"`
	Name        string                       `json:"name"`
	Description string                       `json:"description,omitempty"`
	// NOTE: This feature was on the API, but at some point developers decided to remove it.
	// So I commented it, until it will be officially on the docs.
	// Default     bool                              `json:"default"`
	Required bool                              `json:"required"`
	Choices  []*ApplicationCommandOptionChoice `json:"choices"`
	Options  []*ApplicationCommandOption       `json:"options"`
}

// ApplicationCommandOptionChoice is representing slash-command's option choice.
type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// InteractionType is representing interaction type.
type InteractionType uint8

// Interaction types
const (
	InteractionPing = InteractionType(iota + 1)
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
	ID      string                                     `json:"id"`
	Name    string                                     `json:"name"`
	Options []*ApplicationCommandInteractionDataOption `json:"options"`
}

// ApplicationCommandInteractionDataOption is representing an option of application's command.
type ApplicationCommandInteractionDataOption struct {
	Name string `json:"name"`
	// NOTE: Contains the value specified by InteractionType.
	Value   interface{}                                `json:"value,omitempty"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
}

// IntValue is utility function for casting option value to integer
func (o ApplicationCommandInteractionDataOption) IntValue() int64 {
	if v, ok := o.Value.(float64); ok {
		return int64(v)
	}

	return 0
}

// UintValue is utility function for casting option value to unsigned integer
func (o ApplicationCommandInteractionDataOption) UintValue() uint64 {
	if v, ok := o.Value.(float64); ok {
		return uint64(v)
	}

	return 0
}

// FloatValue is utility function for casting option value to float
func (o ApplicationCommandInteractionDataOption) FloatValue() float64 {
	if v, ok := o.Value.(float64); ok {
		return float64(v)
	}

	return 0.0
}

// StringValue is utility function for casting option value to string
func (o ApplicationCommandInteractionDataOption) StringValue() string {
	if v, ok := o.Value.(string); ok {
		return v
	}

	return ""
}

// BoolValue is utility function for casting option value to bool
func (o ApplicationCommandInteractionDataOption) BoolValue() bool {
	if v, ok := o.Value.(bool); ok {
		return v
	}

	return false
}

// ChannelValue is utility function for casting option value to channel object.
// s : Session object, if not nil, function additionaly fetches all channel's data
func (o ApplicationCommandInteractionDataOption) ChannelValue(s *Session) *Channel {
	chanID := o.StringValue()
	if chanID == "" {
		return nil
	}

	if s == nil {
		return &Channel{ID: chanID}
	}

	ch, err := s.State.Channel(chanID)
	if err != nil {
		ch, err = s.Channel(chanID)
		if err != nil {
			return &Channel{ID: chanID}
		}
	}

	return ch
}

// RoleValue is utility function for casting option value to role object.
// s : Session object, if not nil, function additionaly fetches all role's data
func (o ApplicationCommandInteractionDataOption) RoleValue(s *Session, gID string) *Role {
	roleID := o.StringValue()
	if roleID == "" {
		return nil
	}

	if s == nil || gID == "" {
		return &Role{ID: roleID}
	}

	r, err := s.State.Role(roleID, gID)
	if err != nil {
		roles, err := s.GuildRoles(gID)
		if err == nil {
			for _, r = range roles {
				if r.ID == roleID {
					return r
				}
			}
		}
		return &Role{ID: roleID}
	}

	return r
}

// UserValue is utility function for casting option value to user object.
// s : Session object, if not nil, function additionaly fetches all user's data
func (o ApplicationCommandInteractionDataOption) UserValue(s *Session) *User {
	userID := o.StringValue()
	if userID == "" {
		return nil
	}

	if s == nil {
		return &User{ID: userID}
	}

	u, err := s.User(userID)
	if err != nil {
		return &User{ID: userID}
	}

	return u
}

// InteractionResponseType is type of interaction response.
type InteractionResponseType uint8

// Interaction response types.
const (
	// InteractionResponsePong is for ACK ping event.
	InteractionResponsePong = InteractionResponseType(iota + 1)
	// InteractionResponseAcknowledge is for ACK a command without sending a message, eating the user's input.
	InteractionResponseAcknowledge
	// InteractionResponseChannelMessage is for responding with a message, eating the user's input.
	InteractionResponseChannelMessage
	// InteractionResponseChannelMessageWithSource is for responding with a message, showing the user's input.
	InteractionResponseChannelMessageWithSource
	// InteractionResponseACKWithSource is for ACK a command without sending a message, showing the user's input.
	InteractionResponseACKWithSource
)

// InteractionResponse is representing response for interaction with application.
type InteractionResponse struct {
	Type InteractionResponseType                    `json:"type,omitempty"`
	Data *InteractionApplicationCommandResponseData `json:"data,omitempty"`
}

// InteractionApplicationCommandResponseData is callback data for application command interaction.
type InteractionApplicationCommandResponseData struct {
	TTS             bool                    `json:"tts,omitempty"`
	Content         string                  `json:"content,omitempty"`
	Embeds          []*MessageEmbed         `json:"embeds,omitempty"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`

	// NOTE: Undocumented feature, be careful with it.
	Flags uint64 `json:"flags,omitempty"` 
}

// VerifyInteraction implements message verification of the discord interactions api
// signing algorithm, as documented here:
// https://discord.com/developers/docs/interactions/slash-commands#security-and-authorization
func VerifyInteraction(r *http.Request, key ed25519.PublicKey) bool {
	var msg bytes.Buffer

	signature := r.Header.Get("X-Signature-Ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	if len(sig) != ed25519.SignatureSize {
		return false
	}

	timestamp := r.Header.Get("X-Signature-Timestamp")
	if timestamp == "" {
		return false
	}

	msg.WriteString(timestamp)

	defer r.Body.Close()
	var body bytes.Buffer

	// at the end of the function, copy the original body back into the request
	defer func() {
		r.Body = ioutil.NopCloser(&body)
	}()

	// copy body into buffers
	_, err = io.Copy(&msg, io.TeeReader(r.Body, &body))
	if err != nil {
		return false
	}

	return ed25519.Verify(key, msg.Bytes(), sig)
}
