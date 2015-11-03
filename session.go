/******************************************************************************
 * Discordgo v0 by Bruce Marriner <bruce@sqls.net>
 * A DiscordApp API for Golang.
 *
 * Currently only the REST API is functional.  I will add on the websocket
 * layer once I get the API section where I want it.
 *
 */

package discordgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Represents a session connection to the Discord REST API.
// I suspect I'll be adding more to this later :)
type Session struct {
	Token string
	Debug bool
}

// RequestToken asks the Rest server for a token by provided email/password
func (session *Session) RequestToken(email string, password string) (token string, err error) {

	var urlStr string = fmt.Sprintf("%s/%s", discordApi, "auth/login")
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password))))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: (20 * time.Second)}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}

	if session.Debug {
		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, body, "", "\t")
		if error != nil {
			fmt.Print("JSON parse error: ", error)
			return
		}
		fmt.Println("requestToken Response:\n", string(prettyJSON.Bytes()))
	}

	temp := &Session{} // TODO Must be a better way
	err = json.Unmarshal(body, &temp)
	token = temp.Token
	return
}

// Identify session user
func (session *Session) Self() (user User, err error) {

	body, err := Request(session, fmt.Sprintf("%s/%s", discordApi, "users/@me"))
	err = json.Unmarshal(body, &user)

	return
}

// Request makes a API GET Request.  This is a general purpose function
// and is used by all API functions.  It is exposed currently so it can
// also be used outside of this library.
func Request(session *Session, urlStr string) (body []byte, err error) {

	req, err := http.NewRequest("GET", urlStr, bytes.NewBuffer([]byte(fmt.Sprintf(``))))
	if err != nil {
		return
	}

	req.Header.Set("authorization", session.Token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: (20 * time.Second)}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}

	if session.Debug {
		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, body, "", "\t")
		if error != nil {
			fmt.Print("JSON parse error: ", error)
			return
		}
		fmt.Println(urlStr+" Response:\n", string(prettyJSON.Bytes()))
	}
	return
}

// Get all of the session user's private channels.
func (session *Session) PrivateChannels() (channels []Channel, err error) {

	body, err := Request(session, fmt.Sprintf("%s/%s", discordApi, fmt.Sprintf("users/@me/channels")))
	err = json.Unmarshal(body, &channels)

	return
}

// Get all of the session user's servers
func (session *Session) Servers() (servers []Server, err error) {

	body, err := Request(session, fmt.Sprintf("%s/%s", discordApi, fmt.Sprintf("users/@me/guilds")))
	err = json.Unmarshal(body, &servers)

	return
}

// Get all channels for the given server
func (session *Session) Channels(serverId int) (channels []Channel, err error) {

	body, err := Request(session, fmt.Sprintf("%s/%s", discordApi, fmt.Sprintf("guilds/%d/channels", serverId)))
	err = json.Unmarshal(body, &channels)

	return
}
