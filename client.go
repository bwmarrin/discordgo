/******************************************************************************
 * Discordgo by Bruce Marriner <bruce@sqls.net>
 * A Discord API for Golang.
 * See discord.go for more information.
 *
 * This file contains functions for interacting with the Discord API
 * at the lowest level.  See other files for easier methods of access.
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

// RequestToken asks the Discord server for an authentication token
func Login(session *Session, email string, password string) (token string, err error) {

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

// Request makes a REST API GET Request with Discord.
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

// Returns the user details of the given userId
// session : An active session connection to Discord
// user    : A user Id or name
func Users(session *Session, userId string) (user User, err error) {

	body, err := Request(session, fmt.Sprintf("%s/%s", discordApi, "users/%s", userId))
	err = json.Unmarshal(body, &user)

	return
}

// PrivateChannels returns an array of Channel structures for all private
// channels for a user
func PrivateChannels(session *Session, userId string) (channels []Channel, err error) {

	body, err := Request(session, fmt.Sprintf("%s/%s", discordApi, fmt.Sprintf("users/%s/channels", userId)))
	err = json.Unmarshal(body, &channels)

	return
}

// Servers returns an array of Server structures for all servers for a user
func Servers(session *Session, userId string) (servers []Server, err error) {

	body, err := Request(session, fmt.Sprintf("%s/%s", discordApi, fmt.Sprintf("users/%s/guilds", userId)))
	err = json.Unmarshal(body, &servers)

	return
}

// Channels returns an array of Channel structures for all channels of a given
// server.
func Channels(session *Session, serverId int) (channels []Channel, err error) {

	body, err := Request(session, fmt.Sprintf("%s/%s", discordApi, fmt.Sprintf("guilds/%d/channels", serverId)))
	err = json.Unmarshal(body, &channels)

	return
}

// Close ends a session and logs out from the Discord REST API.
func Close(session *Session) (err error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", discordApi, fmt.Sprintf("auth/logout")), bytes.NewBuffer([]byte(fmt.Sprintf(``))))
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}
	return
}
