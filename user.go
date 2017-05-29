package discordgo

import (
	"fmt"
	"strings"
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

// String returns a unique identifier of the form username#discriminator
func (u *User) String() string {
	return fmt.Sprintf("%s#%s", u.Username, u.Discriminator)
}

// Mention return a string which mentions the user
func (u *User) Mention() string {
	return fmt.Sprintf("<@%s>", u.ID)
}

// AvatarURL returns a URL to the user's avatar.
//		size:     The size of the user's avatar as a power of two
func (u *User) AvatarURL(size string) string {
	return UserAvatarURL(u.ID, u.Avatar, size)
}

// UserAvatarURL returns a URL to the requested user's avatar
//		userID:   The ID of the user to get
//		avatarID: The ID of the user's avatar
//		size:     The size of the user's avatar as a power of two
func UserAvatarURL(userID, avatarID, size string) string {
	extension := ".jpg"
	if strings.HasPrefix(avatarID, "a_") {
		extension = ".gif"
	}
	return EndpointCDNAvatars + userID + "/" + avatarID + extension + "?size=" + size
}
