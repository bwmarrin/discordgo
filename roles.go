package discordgo

import "fmt"

// A Role stores information about Discord guild member roles.
type Role struct {
	// The ID of the role.
	ID string `json:"id"`

	// The name of the role.
	Name string `json:"name"`

	// Whether this role is managed by an integration, and
	// thus cannot be manually added to, or taken from, members.
	Managed bool `json:"managed"`

	// Whether this role is mentionable.
	Mentionable bool `json:"mentionable"`

	// Whether this role is hoisted (shows up separately in member list).
	Hoist bool `json:"hoist"`

	// The hex color of this role.
	Color int `json:"color"`

	// The position of this role in the guild's role hierarchy.
	Position int `json:"position"`

	// The permissions of the role on the guild (doesn't include channel overrides).
	// This is a combination of bit masks; the presence of a certain permission can
	// be checked by performing a bitwise AND between this int and the permission.
	Permissions int `json:"permissions"`

	// ID of the guild this role belongs to
	GuildID string `json:"guild_id,omitempty"`

	// The Session to call the API and retrieve other objects
	Session *Session `json:"session,omitempty"`
}

// A RoleEdit stores information used to edit a Role
type RoleEdit struct {
	// The role's name (overwrites existing)
	Name string `json:"name"`

	// The color the role should have (as a decimal, not hex)
	Color int `json:"color"`

	// Whether to display the role's users separately (overwrites existing)
	Hoist bool `json:"hoist"`

	// The overall permissions number of the role (overwrites existing)
	Permissions int `json:"permissions"`

	// Whether this role is mentionable (overwrites existing)
	Mentionable bool `json:"mentionable"`
}

// GetID returns the ID of the Role
func (r *Role) GetID() string {
	return r.ID
}

// Mention returns a string which mentions the role
func (r *Role) Mention() string {
	return fmt.Sprintf("<@&%s>", r.ID)
}

// GetGuild returns the Guild struct this role belongs to
func (r *Role) GetGuild() (g *Guild, err error) {
	return r.Session.State.Guild(r.GuildID)
}

// Edit updates the Role with new values
// name      : The name of the Role.
// color     : The color of the role (decimal, not hex).
// hoist     : Whether to display the role's users separately.
// perm      : The permissions for the role.
// mention   : Whether this role is mentionable
func (r *Role) Edit(name string, color int, hoist bool, perm int, mention bool) (edited *Role, err error) {
	return r.Session.GuildRoleEdit(r.GuildID, r.ID, name, color, hoist, perm, mention)
}

// EditComplex updates the Role with new values
// data      : data to send to the API
func (r *Role) EditComplex(data *RoleEdit) (edited *Role, err error) {
	return r.Session.GuildRoleEditComplex(r.GuildID, r.ID, data)
}

// Delete deletes the role
func (r *Role) Delete() (err error) {
	return r.Session.GuildRoleDelete(r.GuildID, r.ID)
}

// Roles are a collection of Role
type Roles []*Role

func (r Roles) Len() int {
	return len(r)
}

func (r Roles) Less(i, j int) bool {
	return r[i].Position > r[j].Position
}

func (r Roles) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
