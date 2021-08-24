package discordgo

import (
	"encoding/json"
)

// ComponentType is type of component.
type ComponentType uint

// MessageComponent types.
const (
	ComponentActionRow  ComponentType = 1
	ComponentButton     ComponentType = 2
	ComponentSelectMenu ComponentType = 3
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
func (umc *unmarshalableMessageComponent) UnmarshalJSON(src []byte) error {
	var v struct {
		Type ComponentType `json:"type"`
	}
	err := json.Unmarshal(src, &v)
	if err != nil {
		return err
	}

	var data MessageComponent
	switch v.Type {
	case ComponentActionRow:
		v := ActionsRow{}
		err = json.Unmarshal(src, &v)
		data = v
	case ComponentButton:
		v := Button{}
		err = json.Unmarshal(src, &v)
		data = v
	}
	if err != nil {
		return err
	}
	umc.MessageComponent = data
	return err
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

// Type is a method to get the type of a component.
func (r ActionsRow) Type() ComponentType {
	return ComponentActionRow
}

// ButtonStyle is style of button.
type ButtonStyle uint

// Button styles.
const (
	// ButtonPrimary is a button with blurple color.
	ButtonPrimary ButtonStyle = 1
	// ButtonSecondary is a button with grey color.
	ButtonSecondary ButtonStyle = 2
	// ButtonSuccess is a button with green color.
	ButtonSuccess ButtonStyle = 3
	// ButtonDanger is a button with red color.
	ButtonDanger ButtonStyle = 4
	// ButtonLink is a special type of button which navigates to a URL. Has grey color.
	ButtonLink ButtonStyle = 5
)

// ComponentEmoji represents a button's emoji, if it has one.
type ComponentEmoji struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}

// Button represents button component.
type Button struct {
	Label    string         `json:"label"`
	Style    ButtonStyle    `json:"style"`
	Disabled bool           `json:"disabled"`
	Emoji    ComponentEmoji `json:"emoji"`

	// NOTE: Only button with ButtonLink style can have link. Also, URL is mutually exclusive with CustomID.
	URL      string `json:"url,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
}

// MarshalJSON is a method for marshaling Button to a JSON object.
func (b Button) MarshalJSON() ([]byte, error) {
	type button Button

	if b.Style == 0 {
		b.Style = ButtonPrimary
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
	return ComponentButton
}

// SelectMenuOption represents an option for a select menu.
type SelectMenuOption struct {
	Label       string         `json:"label,omitempty"`
	Value       string         `json:"value"`
	Description string         `json:"description"`
	Emoji       ComponentEmoji `json:"emoji"`
	// Determines whenever option is selected by default or not.
	Default bool `json:"default"`
}

// SelectMenu represents select menu component.
type SelectMenu struct {
	CustomID string `json:"custom_id,omitempty"`
	// The text which will be shown in the menu if there's no default options or all options was deselected and component was closed.
	Placeholder string `json:"placeholder"`
	// This value determines the minimal amount of selected items in the menu.
	MinValues int `json:"min_values,omitempty"`
	// This value determines the maximal amount of selected items in the menu.
	// If MaxValues or MinValues are greater than one then the user can select multiple items in the component.
	MaxValues int                `json:"max_values,omitempty"`
	Options   []SelectMenuOption `json:"options"`
}

// Type is a method to get the type of a component.
func (SelectMenu) Type() ComponentType {
	return ComponentSelectMenu
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
