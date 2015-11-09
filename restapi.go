/******************************************************************************
 * A Discord API for Golang.
 * See discord.go for more information.
 *
 * This file contains functions for interacting with the Discord HTTP REST API
 * at the lowest level.
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

// Request makes a (GET/POST/?) Requests to Discord REST API.
// All the other functions in this file use this function.
func Request(session *Session, method, urlStr, body string) (response []byte, err error) {

	if session.Debug {
		fmt.Println("REQUEST  :: " + method + " " + urlStr + "\n" + body)
	}

	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return
	}

	// Not used on initial login..
	if session.Token != "" {
		req.Header.Set("authorization", session.Token)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: (20 * time.Second)}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(response)))
		return
	}

	if session.Debug {
		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, response, "", "\t")
		if error != nil {
			fmt.Print("JSON parse error: ", error)
			return
		}
		fmt.Println("RESPONSE ::\n" + string(prettyJSON.Bytes()))
	}
	return
}

// Login asks the Discord server for an authentication token
func Login(session *Session, email string, password string) (token string, err error) {

	var urlStr string = fmt.Sprintf("%s/%s", discordApi, "auth/login")

	response, err := Request(session, "POST", urlStr, fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password))

	var temp map[string]interface{}
	err = json.Unmarshal(response, &temp)
	token = temp["token"].(string)

	return
}

// Returns the user details of the given userId
// session : An active session connection to Discord
// user    : A user Id or name
func Users(session *Session, userId string) (user User, err error) {

	body, err := Request(session, "GET", fmt.Sprintf("%s/users/%s", discordApi, userId), ``)
	err = json.Unmarshal(body, &user)
	return
}

// USERS could pull users channels, servers, settings and so forth too?
// you know, pull all the data for the user.  update the user strut
// to house that data.  Seems reasonable.

// PrivateChannels returns an array of Channel structures for all private
// channels for a user
func PrivateChannels(session *Session, userId string) (channels []Channel, err error) {

	body, err := Request(session, "GET", fmt.Sprintf("%s/users/%s/channels", discordApi, userId), ``)
	err = json.Unmarshal(body, &channels)

	return
}

// Servers returns an array of Server structures for all servers for a user
func Servers(session *Session, userId string) (servers []Server, err error) {

	body, err := Request(session, "GET", fmt.Sprintf("%s/users/%s/guilds", discordApi, userId), ``)
	err = json.Unmarshal(body, &servers)

	return
}

// add one to get specific server by ID, or enhance the above with an ID field.
// GET http://discordapp.com/api/guilds/ID#

// Members returns an array of Member structures for all members of a given
// server.
func Members(session *Session, serverId int) (members []Member, err error) {

	body, err := Request(session, "GET", fmt.Sprintf("%s/guilds/%d/members", discordApi, serverId), ``)
	err = json.Unmarshal(body, &members)

	return
}

// Channels returns an array of Channel structures for all channels of a given
// server.
func Channels(session *Session, serverId int) (channels []Channel, err error) {

	body, err := Request(session, "GET", fmt.Sprintf("%s/guilds/%d/channels", discordApi, serverId), ``)
	err = json.Unmarshal(body, &channels)

	return
}

// update above or add a way to get channel by ID.  ChannelByName could be handy
// too you know.
// http://discordapp.com/api/channels/ID#

// Messages returns an array of Message structures for messaages within a given
// channel.  limit, beforeId, and afterId can be used to control what messages
// are returned.
func Messages(session *Session, channelId int, limit int, beforeId int, afterId int) (messages []Message, err error) {

	var urlStr string

	if limit > 0 {
		urlStr = fmt.Sprintf("%s/channels/%d/messages?limit=%d", discordApi, channelId, limit)
	}

	if afterId > 0 {
		if urlStr != "" {
			urlStr = urlStr + fmt.Sprintf("&after=%d", afterId)
		} else {
			urlStr = fmt.Sprintf("%s/channels/%d/messages?after=%d", discordApi, channelId, afterId)
		}
	}

	if beforeId > 0 {
		if urlStr != "" {
			urlStr = urlStr + fmt.Sprintf("&before=%d", beforeId)
		} else {
			urlStr = fmt.Sprintf("%s/channels/%d/messages?after=%d", discordApi, channelId, beforeId)
		}
	}

	if urlStr == "" {
		urlStr = fmt.Sprintf("%s/channels/%d/messages", discordApi, channelId)
	}

	body, err := Request(session, "GET", urlStr, ``)
	err = json.Unmarshal(body, &messages)

	return
}

// SendMessage sends a message to the given channel.
func SendMessage(session *Session, channelId int, content string) (message Message, err error) {

	var urlStr string = fmt.Sprintf("%s/channels/%d/messages", discordApi, channelId)
	response, err := Request(session, "POST", urlStr, fmt.Sprintf(`{"content":"%s"}`, content))
	err = json.Unmarshal(response, &message)

	return
}

// Returns the a websocket Gateway address
// session : An active session connection to Discord
func Gateway(session *Session) (gateway string, err error) {

	response, err := Request(session, "GET", fmt.Sprintf("%s/gateway", discordApi), ``)

	var temp map[string]interface{}
	err = json.Unmarshal(response, &temp)
	gateway = temp["url"].(string)
	return
}

// Close ends a session and logs out from the Discord REST API.
// This does not seem to actually invalidate the token.  So you can still
// make API calls even after a Logout.  So, it seems almost pointless to
// even use.
func Logout(session *Session) (err error) {

	urlStr := fmt.Sprintf("%s/auth/logout", discordApi)
	_, err = Request(session, "POST", urlStr, fmt.Sprintf(`{"token": "%s"}`, session.Token))

	return
}
