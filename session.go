/******************************************************************************
 * A Discord API for Golang.
 */

package discordgo

import "github.com/gorilla/websocket"

// A Session represents a connection to the Discord REST API.
// Token : The authentication token returned from Discord
// Debug : If set to ture debug logging will be displayed.
type Session struct {
	Token     string          // Authentication token for this session
	Gateway   string          // Websocket Gateway for this session
	Debug     bool            // Debug for printing JSON request/responses
	Cache     int             // number in X to cache some responses
	Websocket *websocket.Conn // TODO: use this
	//TODO, add bools for like.
	// are we connnected to websocket?
	// have we authenticated to login?
	// lets put all the general session
	// tracking and infos here.. clearly
}

/******************************************************************************
 * The below functions are "shortcut" methods for functions in restapi.go
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

// Logout ends a session and logs out from the Discord REST API.
func (session *Session) Logout() (err error) {
	err = Logout(session)
	return
}
