package discordgo

import (
	"encoding/json"
	"fmt"
)

// ComponentType is type of component.
type ComponentType uint

// MessageComponent types.
const (
	ActionsRowComponent            ComponentType = 1
	ButtonComponent                ComponentType = 2
	SelectMenuComponent            ComponentType = 3
	TextInputComponent             ComponentType = 4
	UserSelectMenuComponent        ComponentType = 5
	RoleSelectMenuComponent        ComponentType = 6
	MentionableSelectMenuComponent ComponentType = 7
	ChannelSelectMenuComponent     ComponentType = 8
	SectionComponent               ComponentType = 9
	TextDisplayComponent           ComponentType = 10
	ThumbnailComponent             ComponentType = 11
	MediaGalleryComponent          ComponentType = 12
	FileComponent                  ComponentType = 13
	SeparatorComponent             ComponentType = 14
	ContainerComponent             ComponentType = 17
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

	switch v.Type {
	case ActionsRowComponent:
		umc.MessageComponent = &ActionsRow{}
	case ButtonComponent:
		umc.MessageComponent = &Button{}
	case SelectMenuComponent, ChannelSelectMenuComponent, UserSelectMenuComponent,
		RoleSelectMenuComponent, MentionableSelectMenuComponent:
		umc.MessageComponent = &SelectMenu{}
	case TextInputComponent:
		umc.MessageComponent = &TextInput{}
	case SectionComponent:
		umc.MessageComponent = &Section{}
	case TextDisplayComponent:
		umc.MessageComponent = &TextDisplay{}
	case ThumbnailComponent:
		umc.MessageComponent = &Thumbnail{}
	case MediaGalleryComponent:
		umc.MessageComponent = &MediaGallery{}
	case FileComponent:
		umc.MessageComponent = &FileComponentData{}
	case SeparatorComponent:
		umc.MessageComponent = &Separator{}
	case ContainerComponent:
		umc.MessageComponent = &Container{}
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

// ActionsRow is a container for components within one row.
type ActionsRow struct {
	Components []MessageComponent `json:"components"`
}

// MarshalJSON is a method for marshaling ActionsRow to a JSON object.
func (r ActionsRow) MarshalJSON() ([]byte, error) {
	type actionsRow ActionsRow

	return Marshal(struct {
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
	// PremiumButton is a special type of button with a blurple color that links to a SKU.
	PremiumButton ButtonStyle = 6
)

// ComponentEmoji represents button emoji, if it does have one.
type ComponentEmoji struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}

// Button represents button component.
type Button struct {
	Label    string          `json:"label,omitempty"`
	Style    ButtonStyle     `json:"style"`
	Disabled bool            `json:"disabled,omitempty"`
	Emoji    *ComponentEmoji `json:"emoji,omitempty"`

	// NOTE: Only button with LinkButton style can have link. Also, URL is mutually exclusive with CustomID.
	URL      string `json:"url,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
	// Identifier for a purchasable SKU. Only available when using premium-style buttons.
	SKUID string `json:"sku_id,omitempty"`
}

// MarshalJSON is a method for marshaling Button to a JSON object.
func (b Button) MarshalJSON() ([]byte, error) {
	type button Button

	if b.Style == 0 {
		b.Style = PrimaryButton
	}

	return Marshal(struct {
		button
		Type ComponentType `json:"type"`
	}{
		button: button(b),
		Type:   b.Type(),
	})
}

// Type is a method to get the type of a component.
func (Button) Type() ComponentType {
	return ButtonComponent
}

// SelectMenuOption represents an option for a select menu.
type SelectMenuOption struct {
	Label       string          `json:"label"`
	Value       string          `json:"value"`
	Description string          `json:"description,omitempty"`
	Emoji       *ComponentEmoji `json:"emoji,omitempty"`
	// Determines whenever option is selected by default or not.
	Default bool `json:"default,omitempty"`
}

// SelectMenuDefaultValueType represents the type of an entity selected by default in auto-populated select menus.
type SelectMenuDefaultValueType string

// SelectMenuDefaultValue types.
const (
	SelectMenuDefaultValueUser    SelectMenuDefaultValueType = "user"
	SelectMenuDefaultValueRole    SelectMenuDefaultValueType = "role"
	SelectMenuDefaultValueChannel SelectMenuDefaultValueType = "channel"
)

// SelectMenuDefaultValue represents an entity selected by default in auto-populated select menus.
type SelectMenuDefaultValue struct {
	// ID of the entity.
	ID string `json:"id"`
	// Type of the entity.
	Type SelectMenuDefaultValueType `json:"type"`
}

// SelectMenuType represents select menu type.
type SelectMenuType ComponentType

// SelectMenu types.
const (
	StringSelectMenu      = SelectMenuType(SelectMenuComponent)
	UserSelectMenu        = SelectMenuType(UserSelectMenuComponent)
	RoleSelectMenu        = SelectMenuType(RoleSelectMenuComponent)
	MentionableSelectMenu = SelectMenuType(MentionableSelectMenuComponent)
	ChannelSelectMenu     = SelectMenuType(ChannelSelectMenuComponent)
)

// SelectMenu represents select menu component.
type SelectMenu struct {
	// Type of the select menu.
	MenuType SelectMenuType `json:"type,omitempty"`
	// CustomID is a developer-defined identifier for the select menu.
	CustomID string `json:"custom_id"`
	// The text which will be shown in the menu if there's no default options or all options was deselected and component was closed.
	Placeholder string `json:"placeholder,omitempty"`
	// This value determines the minimal amount of selected items in the menu.
	MinValues *int `json:"min_values,omitempty"`
	// This value determines the maximal amount of selected items in the menu.
	// If MaxValues or MinValues are greater than one then the user can select multiple items in the component.
	MaxValues int `json:"max_values,omitempty"`
	// List of default values for auto-populated select menus.
	// NOTE: Number of entries should be in the range defined by MinValues and MaxValues.
	DefaultValues []SelectMenuDefaultValue `json:"default_values,omitempty"`

	Options  []SelectMenuOption `json:"options,omitempty"`
	Disabled bool               `json:"disabled,omitempty"`

	// NOTE: Can only be used in SelectMenu with Channel menu type.
	ChannelTypes []ChannelType `json:"channel_types,omitempty"`
}

// Type is a method to get the type of a component.
func (s SelectMenu) Type() ComponentType {
	if s.MenuType != 0 {
		return ComponentType(s.MenuType)
	}
	return SelectMenuComponent
}

// MarshalJSON is a method for marshaling SelectMenu to a JSON object.
func (s SelectMenu) MarshalJSON() ([]byte, error) {
	type selectMenu SelectMenu

	return Marshal(struct {
		selectMenu
		Type ComponentType `json:"type"`
	}{
		selectMenu: selectMenu(s),
		Type:       s.Type(),
	})
}

// TextInput represents text input component.
type TextInput struct {
	CustomID    string         `json:"custom_id"`
	Label       string         `json:"label"`
	Style       TextInputStyle `json:"style"`
	Placeholder string         `json:"placeholder,omitempty"`
	Value       string         `json:"value,omitempty"`
	Required    bool           `json:"required"`
	MinLength   int            `json:"min_length,omitempty"`
	MaxLength   int            `json:"max_length,omitempty"`
}

// Type is a method to get the type of a component.
func (TextInput) Type() ComponentType {
	return TextInputComponent
}

// MarshalJSON is a method for marshaling TextInput to a JSON object.
func (m TextInput) MarshalJSON() ([]byte, error) {
	type inputText TextInput

	return Marshal(struct {
		inputText
		Type ComponentType `json:"type"`
	}{
		inputText: inputText(m),
		Type:      m.Type(),
	})
}

// TextInputStyle is style of text in TextInput component.
type TextInputStyle uint

// Text styles
const (
	TextInputShort     TextInputStyle = 1
	TextInputParagraph TextInputStyle = 2
)

// Section is a layout component that allows joining text contextually with an accessory.
type Section struct {
	CustomID   *int               `json:"id,omitempty"`
	Components []MessageComponent `json:"components"`          // Array of TextDisplayComponents (1-3)
	Accessory  MessageComponent   `json:"accessory,omitempty"` // Thumbnail or Button
}

// Type returns the component type for Section.
func (c Section) Type() ComponentType {
	return SectionComponent
}

// MarshalJSON marshals the Section component to JSON.
func (c Section) MarshalJSON() ([]byte, error) {
	type section Section
	return Marshal(struct {
		section
		Type ComponentType `json:"type"`
	}{
		section: section(c),
		Type:    c.Type(),
	})
}

// UnmarshalJSON is a helper function to unmarshal Section.
func (s *Section) UnmarshalJSON(data []byte) error {
	type sectionAlias struct {
		CustomID      *int                            `json:"id,omitempty"`
		RawComponents []unmarshalableMessageComponent `json:"components"`
		RawAccessory  *unmarshalableMessageComponent  `json:"accessory,omitempty"`
	}
	var v sectionAlias
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	s.CustomID = v.CustomID
	s.Components = make([]MessageComponent, len(v.RawComponents))
	for i, comp := range v.RawComponents {
		s.Components[i] = comp.MessageComponent
	}
	if v.RawAccessory != nil {
		s.Accessory = v.RawAccessory.MessageComponent
	}

	return nil
}

// TextDisplay is a content component for displaying static text.
type TextDisplay struct {
	CustomID *int   `json:"id,omitempty"`
	Content  string `json:"content"`
}

// Type returns the component type for TextDisplay.
func (c TextDisplay) Type() ComponentType {
	return TextDisplayComponent
}

// MarshalJSON marshals the TextDisplay component to JSON.
func (c TextDisplay) MarshalJSON() ([]byte, error) {
	type textDisplay TextDisplay
	return Marshal(struct {
		textDisplay
		Type ComponentType `json:"type"`
	}{
		textDisplay: textDisplay(c),
		Type:        c.Type(),
	})
}

// UnfurledMediaItem is used within FileComponentData.
// For FileComponent, it only supports attachment references.
type UnfurledMediaItem struct {
	URL         string `json:"url"`                    // e.g., "attachment://filename.ext"
	ProxyURL    string `json:"proxy_url,omitempty"`    // Output only.
	Height      *int   `json:"height,omitempty"`       // Output only.
	Width       *int   `json:"width,omitempty"`        // Output only.
	ContentType string `json:"content_type,omitempty"` // Output only.
}

// FileComponentData displays an attached file.
// Named FileComponentData to avoid conflict with discordgo.File
type FileComponentData struct {
	CustomID *int              `json:"id,omitempty"`
	File     UnfurledMediaItem `json:"file"`
	Spoiler  *bool             `json:"spoiler,omitempty"`
}

// Type returns the component type for FileComponentData.
func (c FileComponentData) Type() ComponentType {
	return FileComponent
}

// MarshalJSON marshals the FileComponentData component to JSON.
func (c FileComponentData) MarshalJSON() ([]byte, error) {
	type fileComponentData FileComponentData
	return Marshal(struct {
		fileComponentData
		Type ComponentType `json:"type"`
	}{
		fileComponentData: fileComponentData(c),
		Type:              c.Type(),
	})
}

// Separator is a layout component that adds vertical padding.
type Separator struct {
	CustomID *int  `json:"id,omitempty"`
	Divider  *bool `json:"divider,omitempty"`
	Spacing  *int  `json:"spacing,omitempty"`
}

// Type returns the component type for Separator.
func (c Separator) Type() ComponentType {
	return SeparatorComponent
}

// MarshalJSON marshals the Separator component to JSON.
func (c Separator) MarshalJSON() ([]byte, error) {
	type separator Separator
	return Marshal(struct {
		separator
		Type ComponentType `json:"type"`
	}{
		separator: separator(c),
		Type:      c.Type(),
	})
}

// Container is a layout component that visually groups a set of components.
type Container struct {
	CustomID    *int               `json:"id,omitempty"`
	Components  []MessageComponent `json:"components"` // ActionRow, TextDisplay, Section, MediaGallery, Separator, or File
	AccentColor *int               `json:"accent_color,omitempty"`
	Spoiler     *bool              `json:"spoiler,omitempty"`
}

// Type returns the component type for Container.
func (c Container) Type() ComponentType {
	return ContainerComponent
}

// MarshalJSON marshals the Container component to JSON.
func (c Container) MarshalJSON() ([]byte, error) {
	type container Container
	return Marshal(struct {
		container
		Type ComponentType `json:"type"`
	}{
		container: container(c),
		Type:      c.Type(),
	})
}

// UnmarshalJSON is a helper function to unmarshal Container.
func (c *Container) UnmarshalJSON(data []byte) error {
	type containerAlias struct {
		CustomID      *int                            `json:"id,omitempty"`
		RawComponents []unmarshalableMessageComponent `json:"components"`
		AccentColor   *int                            `json:"accent_color,omitempty"`
		Spoiler       *bool                           `json:"spoiler,omitempty"`
	}
	var v containerAlias
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	c.CustomID = v.CustomID
	c.AccentColor = v.AccentColor
	c.Spoiler = v.Spoiler
	c.Components = make([]MessageComponent, len(v.RawComponents))
	for i, comp := range v.RawComponents {
		c.Components[i] = comp.MessageComponent
	}
	return nil
}

// MediaGalleryItem defines a single item within a MediaGallery.
type MediaGalleryItem struct {
	Media       UnfurledMediaItem `json:"media"`
	Description string            `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler,omitempty"`
}

// MediaGallery is a component that displays 1-10 media attachments in an organized gallery format.
type MediaGallery struct {
	CustomID *int               `json:"id,omitempty"`
	Items    []MediaGalleryItem `json:"items"` // 1 to 10 media gallery items
}

// Type returns the component type for MediaGallery.
func (c MediaGallery) Type() ComponentType {
	return MediaGalleryComponent
}

// MarshalJSON marshals the MediaGallery component to JSON.
func (c MediaGallery) MarshalJSON() ([]byte, error) {
	type mediaGallery MediaGallery
	return Marshal(struct {
		mediaGallery
		Type ComponentType `json:"type"`
	}{
		mediaGallery: mediaGallery(c),
		Type:         c.Type(),
	})
}

// Thumbnail is a content component that is a small image only usable as an accessory in a section.
type Thumbnail struct {
	CustomID    *int              `json:"id,omitempty"`
	Media       UnfurledMediaItem `json:"media"`
	Description string            `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler,omitempty"`
}

// Type returns the component type for Thumbnail.
func (c Thumbnail) Type() ComponentType {
	return ThumbnailComponent
}

// MarshalJSON marshals the Thumbnail component to JSON.
func (c Thumbnail) MarshalJSON() ([]byte, error) {
	type thumbnail Thumbnail
	return Marshal(struct {
		thumbnail
		Type ComponentType `json:"type"`
	}{
		thumbnail: thumbnail(c),
		Type:      c.Type(),
	})
}
