/******************************************************************************
 * A Discord API for Golang.
 *
 * Currently only the REST API is functional.  I will add on the websocket
 * layer once I get the API section where I want it.
 *
 * The idea is that this file is where we pull together the wsapi, and
 * restapi to create a single do-it-all struct
 *
 * NOTE!!! Currently this file has no purpose, it is here for future
 * access methods. EVERYTHING HERE will just go away or be changed
 * substantially in the future.
 */

package discordgo

// A Discord structure represents a all-inclusive (hopefully) structure to
// access the Discord REST API for a given authenticated user.
type Discord struct {
	Session *Session
	User    User
	Servers []Server
}

// New creates a new connection to Discord and returns a Discord structure.
// This provides an easy entry where most commonly needed information is
// automatically fetched.
// TODO add websocket code in here too
func New(email string, password string) (d *Discord, err error) {

	session := Session{}

	session.Token, err = session.Login(email, password)
	if err != nil {
		return
	}

	user, err := session.Self()
	if err != nil {
		return
	}

	servers, err := session.Servers()

	d = &Discord{session, user, servers}

	return
}

// Renew essentially reruns the New command without creating a new session.
// This will update all the user, server, and channel information that was
// fetched with the New command.  This is not an efficient way of doing this
// but if used infrequently it does provide convenience.
func (d *Discord) Renew() (err error) {

	d.User, err = Users(&d.Session, "@me")
	d.Servers, err = Servers(&d.Session, "@me")

	return
}
