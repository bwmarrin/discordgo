/******************************************************************************
 * Discordgo by Bruce Marriner <bruce@sqls.net>
 * A Discord API for Golang.
 */

package discordgo

// A Session represents a connection to the Discord REST API.
// Token : The authentication token returned from Discord
// Debug : If set to ture debug logging will be displayed.
type Session struct {
	Token string
	Debug bool
}

/******************************************************************************
 * The below functions are "shortcut" methods for functions in client.go
 * Reference the client.go file for more documentation.
 */
func (session *Session) Login(email string, password string) (token string, err error) {
	token, err = Login(session, email, password)
	return
}

func (session *Session) Self() (user User, err error) {
	user, err = Users(session, "@me")
	return
}

func (session *Session) PrivateChannels() (channels []Channel, err error) {
	channels, err = PrivateChannels(session, "@me")
	return
}

func (session *Session) Servers() (servers []Server, err error) {
	servers, err = Servers(session, "@me")
	return
}

// Close ends a session and logs out from the Discord REST API.
func (session *Session) Close() (err error) {
	err = Close(session)
	return
}
