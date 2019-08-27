package discordgo

import "errors"

var (
	NotACustomEmoji   = errors.New("you can't do this to a custom emoji")
	UnknownEmojiGuild = errors.New("the guild that this emoji comes from is not in the cache")
)

// Emoji struct holds data related to Emoji's
type Emoji struct {
	// The ID of the emoji, this is empty if the emoji is not custom
	ID string `json:"id"`

	// The name of the emoji, this is the unicode character of the emoji if it's not custom
	Name string `json:"name"`

	// A list of roles that is allowed to use this emoji, if it is empty, the emoji is unrestricted.
	Roles []string `json:"roles"`

	// if the emoji is managed by an external service
	Managed bool `json:"managed"`

	// If colons are required to use this emoji in the client
	RequireColons bool `json:"require_colons"`

	// If the emoji is animated
	Animated bool `json:"animated"`

	// The user that created the emoji, his can only be retrieved when fetching the emoji
	User *User `json:"user,omitempty"`

	// The Session to call the API and retrieve other objects
	Session *Session `json:"-"`

	// the guild this emoji belongs to
	Guild *Guild `json:"-"`
}

// IsCustom returns true if the emoji is a custom emoji
func (e *Emoji) IsCustom() bool {
	return e.ID != ""
}

// MessageFormat returns a correctly formatted Emoji for use in Message content and embeds
func (e *Emoji) MessageFormat() string {
	if e.ID != "" && e.Name != "" {
		if e.Animated {
			return "<a:" + e.APIName() + ">"
		}

		return "<:" + e.APIName() + ">"
	}

	return e.APIName()
}

// APIName returns an correctly formatted API name for use in the MessageReactions endpoints.
func (e *Emoji) APIName() string {
	if e.ID != "" && e.Name != "" {
		return e.Name + ":" + e.ID
	}
	if e.Name != "" {
		return e.Name
	}
	return e.ID
}

// RoleObjects returns a slice of role objects,
// formed from the slice of strings that is the Roles attribute of Emoji
func (e *Emoji) RoleObjects() (roles []*Role) {
	for _, r := range e.Guild.Roles {
		if Contains(e.Roles, r.ID) {
			roles = append(roles, r)
		}
	}
	return
}

// Delete deletes the emoji
func (e *Emoji) Delete() error {
	if e.ID == "" {
		return NotACustomEmoji
	}

	if e.Guild == nil {
		return UnknownEmojiGuild
	}

	return e.Session.GuildEmojiDelete(e.Guild.ID, e.ID)
}

// EditName edits the name of the custom emoji
// name :  the new name for the custom emoji
func (e *Emoji) EditName(name string) (edited *Emoji, err error) {
	if e.ID == "" {
		err = NotACustomEmoji
		return
	}

	if e.Guild == nil {
		err = UnknownEmojiGuild
		return
	}

	return e.Session.GuildEmojiEdit(e.Guild.ID, e.ID, name, e.Roles)
}

// LimitRoles limits the use of the emoji to the roles given here,
// leave empty to make the emoji unrestricted
// roles :  the list of roles to make the emoji exclusive to
func (e *Emoji) LimitRoles(roles []*Role) (edited *Emoji, err error) {
	if e.ID == "" {
		err = NotACustomEmoji
		return
	}

	if e.Guild == nil {
		err = UnknownEmojiGuild
		return
	}

	var roleIDs []string
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}

	return e.Session.GuildEmojiEdit(e.Guild.ID, e.ID, e.Name, roleIDs)
}
