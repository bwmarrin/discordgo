package discordgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/textproto"
	"time"
)

// ErrJSONUnmarshal represents a json unmarshall error.
var ErrJSONUnmarshal = errors.New("json unmarshal")

// MultipartBodyWithJSON returns the contentType and body for a discord request
// data  : The object to encode for payload_json in the multipart request
// files : Files to include in the request
func MultipartBodyWithJSON(data interface{}, files []*File) (requestContentType string, requestBody []byte, err error) {
	body := &bytes.Buffer{}
	bodywriter := multipart.NewWriter(body)

	payload, err := json.Marshal(data)
	if err != nil {
		return
	}

	var p io.Writer

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="payload_json"`)
	h.Set("Content-Type", "application/json")

	p, err = bodywriter.CreatePart(h)
	if err != nil {
		return
	}

	if _, err = p.Write(payload); err != nil {
		return
	}

	for i, file := range files {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file%d"; filename="%s"`, i, quoteEscaper.Replace(file.Name)))
		contentType := file.ContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		h.Set("Content-Type", contentType)

		p, err = bodywriter.CreatePart(h)
		if err != nil {
			return
		}

		if _, err = io.Copy(p, file.Reader); err != nil {
			return
		}
	}

	err = bodywriter.Close()
	if err != nil {
		return
	}

	return bodywriter.FormDataContentType(), body.Bytes(), nil
}

// unmarshal is a universal Unmarshal method.
func unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrJSONUnmarshal, err)
	}

	return nil
}

// unmarshalableMessageComponent represents a message component that can't be unmarshalled.
type unmarshalableMessageComponent struct {
	MessageComponent
}

// UnmarshalJSON is a helper function to unmarshal MessageComponent object.
func (umc *unmarshalableMessageComponent) UnmarshalJSON(src []byte) error {
	var v struct {
		Type ComponentType `json:"type"`
	}
	err := json.Unmarshal(src, &v)
	if err != nil {
		return err
	}

	switch v.Type {
	case ActionsRowComponent:
		umc.MessageComponent = &ActionsRow{}
	case ButtonComponent:
		umc.MessageComponent = &Button{}
	case SelectMenuComponent:
		umc.MessageComponent = &SelectMenu{}
	case TextInputComponent:
		umc.MessageComponent = &TextInput{}
	default:
		return fmt.Errorf("unknown component type: %d", v.Type)
	}
	return json.Unmarshal(src, umc.MessageComponent)
}

// MessageComponentFromJSON is a helper function for unmarshaling message components
func MessageComponentFromJSON(b []byte) (MessageComponent, error) {
	var u unmarshalableMessageComponent
	err := u.UnmarshalJSON(b)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal into MessageComponent: %w", err)
	}
	return u.MessageComponent, nil
}

// MarshalJSON is a method for marshaling ActionsRow to a JSON object.
func (r ActionsRow) MarshalJSON() ([]byte, error) {
	type actionsRow ActionsRow

	return json.Marshal(struct {
		actionsRow
		Type ComponentType `json:"type"`
	}{
		actionsRow: actionsRow(r),
		Type:       r.Type(),
	})
}

// UnmarshalJSON is a helper function to unmarshal Actions Row.
func (r *ActionsRow) UnmarshalJSON(data []byte) error {
	var v struct {
		RawComponents []unmarshalableMessageComponent `json:"components"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	r.Components = make([]MessageComponent, len(v.RawComponents))
	for i, v := range v.RawComponents {
		r.Components[i] = v.MessageComponent
	}

	return err
}

// UnmarshalJSON is a helper function to unmarshal MessageCreate object.
func (m *MessageCreate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

// UnmarshalJSON is a helper function to unmarshal MessageUpdate object.
func (m *MessageUpdate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

// UnmarshalJSON is a helper function to unmarshal MessageDelete object.
func (m *MessageDelete) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

type rawInteraction struct {
	interaction
	Data json.RawMessage `json:"data"`
}

// UnmarshalJSON is a helper function to unmarshal Interaction object.
func (i *InteractionCreate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &i.Interaction)
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
	case InteractionApplicationCommand, InteractionApplicationCommandAutocomplete:
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
	case InteractionModalSubmit:
		v := ModalSubmitInteractionData{}
		err = json.Unmarshal(tmp.Data, &v)
		if err != nil {
			return err
		}
		i.Data = v
	}
	return nil
}

// UnmarshalJSON is a helper function to correctly unmarshal Components.
func (d *ModalSubmitInteractionData) UnmarshalJSON(data []byte) error {
	type modalSubmitInteractionData ModalSubmitInteractionData
	var v struct {
		modalSubmitInteractionData
		RawComponents []unmarshalableMessageComponent `json:"components"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*d = ModalSubmitInteractionData(v.modalSubmitInteractionData)
	d.Components = make([]MessageComponent, len(v.RawComponents))
	for i, v := range v.RawComponents {
		d.Components[i] = v.MessageComponent
	}
	return err
}

// UnmarshalJSON is a helper function to unmarshal the Message.
func (m *Message) UnmarshalJSON(data []byte) error {
	type message Message
	var v struct {
		message
		RawComponents []unmarshalableMessageComponent `json:"components"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*m = Message(v.message)
	m.Components = make([]MessageComponent, len(v.RawComponents))
	for i, v := range v.RawComponents {
		m.Components[i] = v.MessageComponent
	}
	return err
}

// UnmarshalJSON unmarshals JSON into TimeStamps struct
func (t *TimeStamps) UnmarshalJSON(b []byte) error {
	temp := struct {
		End   float64 `json:"end,omitempty"`
		Start float64 `json:"start,omitempty"`
	}{}
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	t.EndTimestamp = int64(temp.End)
	t.StartTimestamp = int64(temp.Start)
	return nil
}

// UnmarshalJSON helps support translation of a milliseconds-based float
// into a time.Duration on TooManyRequests.
func (t *TooManyRequests) UnmarshalJSON(b []byte) error {
	u := struct {
		Bucket     string  `json:"bucket"`
		Message    string  `json:"message"`
		RetryAfter float64 `json:"retry_after"`
	}{}
	err := json.Unmarshal(b, &u)
	if err != nil {
		return err
	}

	t.Bucket = u.Bucket
	t.Message = u.Message
	whole, frac := math.Modf(u.RetryAfter)
	t.RetryAfter = time.Duration(whole)*time.Second + time.Duration(frac*1000)*time.Millisecond
	return nil
}

// UnmarshalJSON is a custom unmarshaljson to make CreatedAt a time.Time instead of an int
func (activity *Activity) UnmarshalJSON(b []byte) error {
	temp := struct {
		Name          string       `json:"name"`
		Type          ActivityType `json:"type"`
		URL           string       `json:"url,omitempty"`
		CreatedAt     int64        `json:"created_at"`
		ApplicationID string       `json:"application_id,omitempty"`
		State         string       `json:"state,omitempty"`
		Details       string       `json:"details,omitempty"`
		Timestamps    TimeStamps   `json:"timestamps,omitempty"`
		Emoji         Emoji        `json:"emoji,omitempty"`
		Party         Party        `json:"party,omitempty"`
		Assets        Assets       `json:"assets,omitempty"`
		Secrets       Secrets      `json:"secrets,omitempty"`
		Instance      bool         `json:"instance,omitempty"`
		Flags         int          `json:"flags,omitempty"`
	}{}
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	activity.CreatedAt = time.Unix(0, temp.CreatedAt*1000000)
	activity.ApplicationID = temp.ApplicationID
	activity.Assets = temp.Assets
	activity.Details = temp.Details
	activity.Emoji = temp.Emoji
	activity.Flags = temp.Flags
	activity.Instance = temp.Instance
	activity.Name = temp.Name
	activity.Party = temp.Party
	activity.Secrets = temp.Secrets
	activity.State = temp.State
	activity.Timestamps = temp.Timestamps
	activity.Type = temp.Type
	activity.URL = temp.URL
	return nil
}

// MarshalJSON is a method for marshaling Button to a JSON object.
func (b Button) MarshalJSON() ([]byte, error) {
	type button Button

	if b.Style == 0 {
		b.Style = PrimaryButton
	}

	return json.Marshal(struct {
		button
		Type ComponentType `json:"type"`
	}{
		button: button(b),
		Type:   b.Type(),
	})
}

// MarshalJSON is a method for marshaling SelectMenu to a JSON object.
func (m SelectMenu) MarshalJSON() ([]byte, error) {
	type selectMenu SelectMenu

	return json.Marshal(struct {
		selectMenu
		Type ComponentType `json:"type"`
	}{
		selectMenu: selectMenu(m),
		Type:       m.Type(),
	})
}

// MarshalJSON is a method for marshaling TextInput to a JSON object.
func (m TextInput) MarshalJSON() ([]byte, error) {
	type inputText TextInput

	return json.Marshal(struct {
		inputText
		Type ComponentType `json:"type"`
	}{
		inputText: inputText(m),
		Type:      m.Type(),
	})
}

// MarshalJSON is a helper function to marshal GuildScheduledEventParams
func (p GuildScheduledEventParams) MarshalJSON() ([]byte, error) {
	type guildScheduledEventParams GuildScheduledEventParams

	if p.EntityType == GuildScheduledEventEntityTypeExternal && p.ChannelID == "" {
		return json.Marshal(struct {
			guildScheduledEventParams
			ChannelID json.RawMessage `json:"channel_id"`
		}{
			guildScheduledEventParams: guildScheduledEventParams(p),
			ChannelID:                 json.RawMessage("null"),
		})
	}

	return json.Marshal(guildScheduledEventParams(p))
}
