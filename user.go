package discordgo

import (
	"fmt"
)

// A User stores all data for an individual Discord user.
type User struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	Token         string `json:"token"`
	Verified      bool   `json:"verified"`
	MFAEnabled    bool   `json:"mfa_enabled"`
	Bot           bool   `json:"bot"`
}

//String returns a unique identifier of the form username#discriminator
func (u *User) String() string {
	return fmt.Sprintf("%v#%v", u.Username, u.Discriminator)
}
