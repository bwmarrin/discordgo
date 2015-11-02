/******************************************************************************
 * Discordgo v0 by Bruce Marriner <bruce@sqls.net>
 * A DiscordApp API for Golang.
 *
 * Currently only the REST API is functional.  I will add on the websocket
 * layer once I get the API section where I want it.
 *
 */

package discordgo

// Define known API URL paths as global constants
const (
	discordUrl = "http://discordapp.com"
	discordApi = discordUrl + "/api/"
	servers    = discordApi + "guilds"
	channels   = discordApi + "channels"
	users      = discordApi + "users"
)

// possible all-inclusive strut..
type Discord struct {
	Session
	User    User
	Servers []Server
}

// Create a new connection to Discord API.  Returns a client session handle.
// this is a all inclusive type of easy setup command that will return
// a connection, user information, and available channels.
// This is probably the most common way to use the library but you
// can use the "manual" functions below instead.
func New(email string, password string) (discord *Discord, err error) {

	session := Session{}

	session.Token, err = session.RequestToken(email, password)
	if err != nil {
		return
	}

	user, err := session.Self()
	if err != nil {
		return
	}

	servers, err := session.Servers()

	discord = &Discord{session, user, servers}

	return
}
