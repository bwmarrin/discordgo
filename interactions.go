package discordgo

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// InteractionDeadline is the time allowed to respond to an interaction.
const InteractionDeadline = time.Second * 3

// ApplicationCommand represents an application's slash command.
type ApplicationCommand struct {
	ID            string                      `json:"id,omitempty"`
	ApplicationID string                      `json:"application_id,omitempty"`
	Name          string                      `json:"name"`
	Description   string                      `json:"description,omitempty"`
	Version       string                      `json:"version,omitempty"`
	Options       []*ApplicationCommandOption `json:"options"`
}

// ApplicationCommandOptionType indicates the type of a slash command's option.
type ApplicationCommandOptionType uint8

// Application command option types.
const (
	ApplicationCommandOptionSubCommand      ApplicationCommandOptionType = 1
	ApplicationCommandOptionSubCommandGroup ApplicationCommandOptionType = 2
	ApplicationCommandOptionString          ApplicationCommandOptionType = 3
	ApplicationCommandOptionInteger         ApplicationCommandOptionType = 4
	ApplicationCommandOptionBoolean         ApplicationCommandOptionType = 5
	ApplicationCommandOptionUser            ApplicationCommandOptionType = 6
	ApplicationCommandOptionChannel         ApplicationCommandOptionType = 7
	ApplicationCommandOptionRole            ApplicationCommandOptionType = 8
	ApplicationCommandOptionMentionable     ApplicationCommandOptionType = 9
)

// ApplicationCommandOption represents an option/subcommand/subcommands group.
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

// ApplicationCommandOptionChoice represents a slash command option choice.
type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// InteractionType indicates the type of an interaction event.
type InteractionType uint8

// Interaction types
const (
	InteractionPing               InteractionType = 1
	InteractionApplicationCommand InteractionType = 2
	InteractionMessageComponent   InteractionType = 3
)

// Interaction represents data of an interaction.
type Interaction struct {
	ID        string          `json:"id"`
	Type      InteractionType `json:"type"`
	Data      InteractionData `json:"-"`
	GuildID   string          `json:"guild_id"`
	ChannelID string          `json:"channel_id"`

	// The message on which interaction was used.
	// NOTE: this field is only filled when a button click triggered the interaction. Otherwise it will be nil.
	Message *Message `json:"message"`

	// The member who invoked this interaction.
	// NOTE: this field is only filled when the slash command was invoked in a guild;
	// if it was invoked in a DM, the `User` field will be filled instead.
	// Make sure to check for `nil` before using this field.
	Member *Member `json:"member"`
	// The user who invoked this interaction.
	// NOTE: this field is only filled when the slash command was invoked in a DM;
	// if it was invoked in a guild, the `Member` field will be filled instead.
	// Make sure to check for `nil` before using this field.
	User *User `json:"user"`

	Token   string `json:"token"`
	Version int    `json:"version"`
}

type interaction Interaction

type rawInteraction struct {
	interaction
	Data json.RawMessage `json:"data"`
}

// UnmarshalJSON is a method for unmarshalling JSON object to Interaction.
func (i *Interaction) UnmarshalJSON(raw []byte) error {
	var tmp rawInteraction
	err := json.Unmarshal(raw, &tmp)
	if err != nil {
		return err
	}

	*i = Interaction(tmp.interaction)

	switch tmp.Type {
	case InteractionApplicationCommand:
		v := ApplicationCommandInteractionData{}
		err = json.Unmarshal(tmp.Data, &v)
		if err != nil {
			return err
		}
		i.Data = v
	case InteractionMessageComponent:
		v := MessageComponentInteractionData{}
		err = json.Unmarshal(tmp.Data, &v)
		if err != nil {
			return err
		}
		i.Data = v
	}
	return nil
}

// MessageComponentData is helper function to assert the inner InteractionData to MessageComponentInteractionData.
// Make sure to check that the Type of the interaction is InteractionMessageComponent before calling.
func (i Interaction) MessageComponentData() (data MessageComponentInteractionData) {
	return i.Data.(MessageComponentInteractionData)
}

// ApplicationCommandData is helper function to assert the inner InteractionData to ApplicationCommandInteractionData.
// Make sure to check that the Type of the interaction is InteractionApplicationCommand before calling.
func (i Interaction) ApplicationCommandData() (data ApplicationCommandInteractionData) {
	return i.Data.(ApplicationCommandInteractionData)
}

// InteractionData is a common interface for all types of interaction data.
type InteractionData interface {
	Type() InteractionType
}

// ApplicationCommandInteractionData contains the data of application command interaction.
type ApplicationCommandInteractionData struct {
	ID       string                                     `json:"id"`
	Name     string                                     `json:"name"`
	Resolved *ApplicationCommandInteractionDataResolved `json:"resolved"`
	Options  []*ApplicationCommandInteractionDataOption `json:"options"`
}

// ApplicationCommandInteractionDataResolved contains resolved data for command arguments.
// Partial Member objects are missing user, deaf and mute fields.
// Partial Channel objects only have id, name, type and permissions fields.
type ApplicationCommandInteractionDataResolved struct {
	Users    map[string]*User    `json:"users"`
	Members  map[string]*Member  `json:"members"`
	Roles    map[string]*Role    `json:"roles"`
	Channels map[string]*Channel `json:"channels"`
}

// Type returns the type of interaction data.
func (ApplicationCommandInteractionData) Type() InteractionType {
	return InteractionApplicationCommand
}

// MessageComponentInteractionData contains the data of message component interaction.
type MessageComponentInteractionData struct {
	CustomID      string        `json:"custom_id"`
	ComponentType ComponentType `json:"component_type"`
}

// Type returns the type of interaction data.
func (MessageComponentInteractionData) Type() InteractionType {
	return InteractionMessageComponent
}

// ApplicationCommandInteractionDataOption represents an option of a slash command.
type ApplicationCommandInteractionDataOption struct {
	Name string `json:"name"`
	// NOTE: Contains the value specified by InteractionType.
	Value   interface{}                                `json:"value,omitempty"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
}

// IntValue is a utility function for casting option value to integer
func (o ApplicationCommandInteractionDataOption) IntValue() int64 {
	if v, ok := o.Value.(float64); ok {
		return int64(v)
	}

	return 0
}

// UintValue is a utility function for casting option value to unsigned integer
func (o ApplicationCommandInteractionDataOption) UintValue() uint64 {
	if v, ok := o.Value.(float64); ok {
		return uint64(v)
	}

	return 0
}

// FloatValue is a utility function for casting option value to float
func (o ApplicationCommandInteractionDataOption) FloatValue() float64 {
	if v, ok := o.Value.(float64); ok {
		return v
	}

	return 0.0
}

// StringValue is a utility function for casting option value to string
func (o ApplicationCommandInteractionDataOption) StringValue() string {
	if v, ok := o.Value.(string); ok {
		return v
	}

	return ""
}

// BoolValue is a utility function for casting option value to bool
func (o ApplicationCommandInteractionDataOption) BoolValue() bool {
	if v, ok := o.Value.(bool); ok {
		return v
	}

	return false
}

// ChannelValue is a utility function for casting option value to channel object.
// s : Session object, if not nil, function additionally fetches all channel's data
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

// RoleValue is a utility function for casting option value to role object.
// s : Session object, if not nil, function additionally fetches all role's data
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

// UserValue is a utility function for casting option value to user object.
// s : Session object, if not nil, function additionally fetches all user's data
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
	InteractionResponsePong InteractionResponseType = 1
	// InteractionResponseChannelMessageWithSource is for responding with a message, showing the user's input.
	InteractionResponseChannelMessageWithSource InteractionResponseType = 4
	// InteractionResponseDeferredChannelMessageWithSource acknowledges that the event was received, and that a follow-up will come later.
	InteractionResponseDeferredChannelMessageWithSource InteractionResponseType = 5
	// InteractionResponseDeferredMessageUpdate acknowledges that the message component interaction event was received, and message will be updated later.
	InteractionResponseDeferredMessageUpdate InteractionResponseType = 6
	// InteractionResponseUpdateMessage is for updating the message to which message component was attached.
	InteractionResponseUpdateMessage InteractionResponseType = 7
)

// InteractionResponse represents a response for an interaction event.
type InteractionResponse struct {
	Type  InteractionResponseType  `json:"type,omitempty"`
  Files []*File                 `json:"-"`
	Data  *InteractionResponseData `json:"data,omitempty"`
}

// InteractionResponseData is response data for an interaction.
type InteractionResponseData struct {
	TTS             bool                    `json:"tts"`
	Content         string                  `json:"content"`
	Components      []MessageComponent      `json:"components"`
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
