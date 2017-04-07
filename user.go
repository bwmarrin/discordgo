package discordgo

import (
	"fmt"
)

//String returns a unique identifier of the form username#discriminator
func (u *User) String() string {
	return fmt.Sprintf("%v#%v", u.Username, u.Discriminator)
}
