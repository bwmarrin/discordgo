package discordgo

import (
	"encoding/json"
)

// ComponentType is type of component.
type ComponentType uint

// Component types.
const (
	ActionsRowComponent ComponentType = iota + 1
	ButtonComponent
)

// Component is a base interface for all components
type Component interface {
	json.Marshaler
	Type() ComponentType
}

// ActionsRow is a container for components within one row.
type ActionsRow struct {
	Components []Component `json:"components"`
}

func (r ActionsRow) MarshalJSON() ([]byte, error) {
	type actionRow ActionsRow

	return json.Marshal(struct {
		actionRow
		Type ComponentType `json:"type"`
	}{
		actionRow: actionRow(r),
		Type:      r.Type(),
	})
}

func (r ActionsRow) Type() ComponentType {
	return ActionsRowComponent
}

// ButtonStyle is style of button.
type ButtonStyle uint

// Button styles.
const (
	// PrimaryButton is a button with blurple color.
	PrimaryButton ButtonStyle = iota + 1
	// SecondaryButton is a button with grey color.
	SecondaryButton
	// SuccessButton is a button with green color.
	SuccessButton
	// DangerButton is a button with red color.
	DangerButton
	// LinkButton is a special type of button which navigates to a URL. Has grey color.
	LinkButton
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

func (b Button) Type() ComponentType {
	return ButtonComponent
}
