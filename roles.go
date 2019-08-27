package discordgo

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"
)

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
	Guild *Guild `json:"-"`

	// The Session to call the API and retrieve other objects
	Session *Session `json:"-"`
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

// A RoleMove stores the information needed to change a role's position in the role hierarchy
type RoleMove struct {
	// The role's ID
	ID string `json:"id"`

	// The position of this role in the guild's role hierarchy.
	Position int `json:"position"`
}

// GetID returns the ID of the Role
func (r *Role) GetID() string {
	return r.ID
}

// CreatedAt returns the roles creation time in UTC
func (r *Role) CreatedAt() (creation time.Time, err error) {
	return SnowflakeToTime(r.ID)
}

// Mention returns a string which mentions the role
func (r *Role) Mention() string {
	return fmt.Sprintf("<@&%s>", r.ID)
}

// IsDefault checks if the Role is the default (@everyone) role
func (r *Role) IsDefault() bool {
	return r.ID == r.Guild.ID
}

// GetMembers returns a slice with all members in the guild with this role
func (r *Role) GetMembers() (members []*Member, err error) {
	allMembers := r.Guild.Members
	for _, m := range allMembers {
		for _, roleID := range m.Roles {
			if roleID == r.ID {
				members = append(members, m)
			}
		}
	}
	return
}

// Edit updates the Role with new values
// name      : The name of the Role.
// color     : The color of the role (decimal, not hex).
// hoist     : Whether to display the role's users separately.
// perm      : The permissions for the role.
// mention   : Whether this role is mentionable
func (r *Role) Edit(name string, color int, hoist bool, perm int, mention bool) (edited *Role, err error) {
	return r.Session.GuildRoleEdit(r.Guild.ID, r.ID, name, color, hoist, perm, mention)
}

// EditComplex updates the Role with new values
// data      : data to send to the API
func (r *Role) EditComplex(data *RoleEdit) (edited *Role, err error) {
	return r.Session.GuildRoleEditComplex(r.Guild.ID, r.ID, data)
}

// Move changes the position of the role in the role hierarchy
// position    : the new position of the role
func (r *Role) Move(position int) (err error) {
	if position <= 0 {
		return errors.New("role position can not be 0 or lower")
	}

	if r.IsDefault() {
		return errors.New("can't move the default role")
	}

	var editedRoles Roles
	var edits []*RoleMove
	min := int(math.Min(float64(position), float64(r.Position)))
	max := int(math.Max(float64(position), float64(r.Position)))

	for _, role := range r.Guild.Roles {
		if role.ID != r.ID && role.Position <= max && role.Position >= min {
			editedRoles = append(editedRoles, role)
		}
	}

	sort.Sort(sort.Reverse(editedRoles))

	if position == min {
		editedRoles = append(Roles{r}, editedRoles...)
	} else {
		editedRoles = append(editedRoles, r)
	}

	for p, i := min, 0; p <= max && i < len(editedRoles); p, i = p+1, i+1 {
		editedRoles[i].Position = p
		edits = append(edits, &RoleMove{editedRoles[i].ID, editedRoles[i].Position})
	}

	_, err = r.Session.GuildRoleReorder(r.Guild.ID, edits)
	if err != nil {
		return
	}
	return
}

// Delete deletes the role
func (r *Role) Delete() (err error) {
	return r.Session.GuildRoleDelete(r.Guild.ID, r.ID)
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

// ContainsID checks if the slice of Role objects contains a role with the given ID
// ID     : the ID to search for
func (r Roles) ContainsID(ID string) bool {
	for _, role := range r {
		if role.ID == ID {
			return true
		}
	}
	return false
}
