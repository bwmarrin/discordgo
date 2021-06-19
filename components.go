package discordgo

import (
	"encoding/json"
)

// ComponentType is type of component.
type ComponentType uint

// MessageComponent types.
const (
	ActionsRowComponent ComponentType = 1
	ButtonComponent     ComponentType = 2
)

// MessageComponent is a base interface for all message components.
type MessageComponent interface {
	json.Marshaler
	Type() ComponentType
}

type unmarshalableMessageComponent struct {
	MessageComponent
}

// UnmarshalJSON is a helper function to unmarshal MessageComponent object.
func (umc *unmarshalableMessageComponent) UnmarshalJSON(src []byte) (err error) {
	var v struct {
		Type ComponentType `json:"type"`
	}
	err = json.Unmarshal(src, &v)
	if err != nil {
		return
	}

	var data MessageComponent
	switch v.Type {
	case ActionsRowComponent:
		v := ActionsRow{}
		err = json.Unmarshal(src, &v)
		data = v
	case ButtonComponent:
		v := Button{}
		err = json.Unmarshal(src, &v)
		data = v
	}
	if err != nil {
		return
	}
	umc.MessageComponent = data
	return
}

// ActionsRow is a container for components within one row.
type ActionsRow struct {
	Components []MessageComponent `json:"components"`
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

// Type is a method to get the type of a component.
func (r ActionsRow) Type() ComponentType {
	return ActionsRowComponent
}

// ButtonStyle is style of button.
type ButtonStyle uint

// Button styles.
const (
	// PrimaryButton is a button with blurple color.
	PrimaryButton ButtonStyle = 1
	// SecondaryButton is a button with grey color.
	SecondaryButton ButtonStyle = 2
	// SuccessButton is a button with green color.
	SuccessButton ButtonStyle = 3
	// DangerButton is a button with red color.
	DangerButton ButtonStyle = 4
	// LinkButton is a special type of button which navigates to a URL. Has grey color.
	LinkButton ButtonStyle = 5
)

// ButtonEmoji represents button emoji, if it does have one.
type ButtonEmoji struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}

// Button represents button component.
type Button struct {
	Label    string      `json:"label"`
	Style    ButtonStyle `json:"style"`
	Disabled bool        `json:"disabled"`
	Emoji    ButtonEmoji `json:"emoji"`

	// NOTE: Only button with LinkButton style can have link. Also, Link is mutually exclusive with CustomID.
	Link     string `json:"url,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
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

// Type is a method to get the type of a component.
func (b Button) Type() ComponentType {
	return ButtonComponent
}
