package discordgo

import "strings"

// A User stores all data for an individual Discord user.
type User struct {
	// The ID of the user.
	ID string `json:"id"`

	// The email of the user. This is only present when
	// the application possesses the email scope for the user.
	Email string `json:"email"`

	// The user's username.
	Username string `json:"username"`

	// The hash of the user's avatar. Use Session.UserAvatar
	// to retrieve the avatar itself.
	Avatar string `json:"avatar"`

	// The user's chosen language option.
	Locale string `json:"locale"`

	// The discriminator of the user (4 numbers after name).
	Discriminator string `json:"discriminator"`

	// The token of the user. This is only present for
	// the user represented by the current session.
	Token string `json:"token"`

	// Whether the user's email is verified.
	Verified bool `json:"verified"`

	// Whether the user has multi-factor authentication enabled.
	MFAEnabled bool `json:"mfa_enabled"`

	// Whether the user is a bot.
	Bot bool `json:"bot"`

	// dm channel with the user, call CreateDM if it doesn't exist
	DMChannel *Channel `json:"dm_channel,omitempty"`

	// The Session to call the API and retrieve other objects
	Session *Session `json:"session,omitempty"`
}

// String returns a unique identifier of the form username#discriminator
func (u *User) String() string {
	return u.Username + "#" + u.Discriminator
}

// Mention return a string which mentions the user
func (u *User) Mention() string {
	return "<@!" + u.ID + ">"
}

func (u *User) GetID() string {
	return u.ID
}

// AvatarURL returns a URL to the user's avatar.
//    size:    The size of the user's avatar as a power of two
//             if size is an empty string, no size parameter will
//             be added to the URL.
func (u *User) AvatarURL(size string) string {
	var URL string
	if u.Avatar == "" {
		URL = EndpointDefaultUserAvatar(u.Discriminator)
	} else if strings.HasPrefix(u.Avatar, "a_") {
		URL = EndpointUserAvatarAnimated(u.ID, u.Avatar)
	} else {
		URL = EndpointUserAvatar(u.ID, u.Avatar)
	}

	if size != "" {
		return URL + "?size=" + size
	}
	return URL
}

func (u *User) CreateDM() (err error) {
	channel, err := u.Session.UserChannelCreate(u.ID)
	if err == nil {
		u.DMChannel = channel
	}
	return
}

func (u *User) SendMessage(content string, embed *MessageEmbed, files []*File) (message *Message, err error) {
	if u.DMChannel == nil {
		err = u.CreateDM()
		if err != nil {
			return
		}
	}

	return u.DMChannel.SendMessage(content, embed, files)
}

func (u *User) SendMessageComplex(data *MessageSend) (message *Message, err error) {
	if u.DMChannel == nil {
		err = u.CreateDM()
		if err != nil {
			return
		}
	}

	return u.DMChannel.SendMessageComplex(data)
}

func (u *User) EditMessage(message *Message) (edited *Message, err error) {
	if u.DMChannel == nil {
		err = u.CreateDM()
		if err != nil {
			return
		}
	}

	return u.DMChannel.EditMessage(message)
}

func (u *User) EditMessageComplex(data *MessageEdit) (edited *Message, err error) {
	if u.DMChannel == nil {
		err = u.CreateDM()
		if err != nil {
			return
		}
	}

	return u.DMChannel.EditMessageComplex(data)
}

func (u *User) FetchMessage(id string) (message *Message, err error) {
	if u.DMChannel == nil {
		err = u.CreateDM()
		if err != nil {
			return
		}
	}

	return u.DMChannel.FetchMessage(id)
}
