package discordgo

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	// "fmt"
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
	Type        ApplicationCommandOptionType      `json:"type"`
	Name        string                            `json:"name"`
	Description string                            `json:"description,omitempty"`
	// Default     bool                              `json:"default"`
	Required    bool                              `json:"required"`
	Choices     []*ApplicationCommandOptionChoice `json:"choices"`
	Options     []*ApplicationCommandOption       `json:"options"`
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
	ID      string `json:"id"`
	Name    string `json:"name"`
	Options []*ApplicationCommandInteractionDataOption `json:"options"`
}

// ApplicationCommandInteractionDataOption is representing an option of application's command.
type ApplicationCommandInteractionDataOption struct {
	Name string `json:"name"`
	// Contains the value specified by InteractionType
	Value   interface{}                                `json:"value,omitempty"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
}

func (o ApplicationCommandInteractionDataOption) IntValue() (i int64) {
	if v, ok := o.Value.(float64); ok {
		i = int64(v)
	}
	return
}

func (o ApplicationCommandInteractionDataOption) UintValue() (i uint64) {
	if v, ok := o.Value.(float64); ok {
		i = uint64(v)
	}
	return
}

func (o ApplicationCommandInteractionDataOption) FloatValue() (n float64) {
	if v, ok := o.Value.(float64); ok {
		n = float64(v)
	}
	return
}

func (o ApplicationCommandInteractionDataOption) StringValue() (s string) {
	if v, ok := o.Value.(string); ok {
		s = v
	}
	return
}
func (o ApplicationCommandInteractionDataOption) BoolValue() (b bool) {
	if v, ok := o.Value.(bool); ok {
		b = v
	}
	return
}

func (o ApplicationCommandInteractionDataOption) ChannelValue(s *Session) (ch *Channel) {
	chanID := o.StringValue()
	if chanID == "" { return }

	if s == nil {
		return &Channel {ID: chanID}
	}

	ch, err := s.State.Channel(chanID)
	if err != nil {
		ch, err = s.Channel(chanID)
	}

	return
}

func (o ApplicationCommandInteractionDataOption) RoleValue(s *Session, gID string) (r *Role) {
	roleID := o.StringValue()
	if roleID == "" { return }

	if s == nil || gID == "" {
		return &Role { ID: roleID }
	}

	var err error
	r, err = s.State.Role(roleID, gID)
	if err != nil {
		roles, err := s.GuildRoles(gID)
		if err != nil {
			return
		}
		for _, r = range roles {
			if r.ID == roleID {
				return
			}
		}
		r = nil
	}

	return
}

func (o ApplicationCommandInteractionDataOption) UserValue(s *Session) (u *User) {
	userID := o.StringValue()
	if userID == "" { return }

	if s == nil {
		return &User { ID: userID }
	}

	u, _ = s.User(userID)

	return
}

// func (o ApplicationCommandInteractionDataOption) String() string {
// 	return fmt.Sprintf("%v", o.Value)
// }

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

	Flags uint64 `json:"flags,omitempty"` // NOTE: Undocumented feature, be careful with it.
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
