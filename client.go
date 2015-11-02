/******************************************************************************

Known API Commands:

Login        - POST http://discordapp.com/api/auth/login
Send Message - POST http://discordapp.com/api/channels/107877361818570752/messages

About Self   - GET  http://discordapp.com/api/users/@me
Guild List   - GET  http://discordapp.com/api/users/90975935880241152/guilds
Channel List - GET  http://discordapp.com/api/guilds/107877361818570752/channels
Get Messages - GET  http://discordapp.com/api/channels/107877361818570752/messages
Get PM Channels - GET  http://discordapp.com/api/users/@me/channels
Get Guild Members - GET http://discordapp.com/api/guilds/107877361818570752/members


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

var discordUrl = "http://discordapp.com/api"

type RestClient struct {
	Url     string
	Session *Session
	client  *http.Client
	Debug   bool
}

type Session struct {
	Id       string
	Email    string
	Password string
	Token    string
}

type Guild struct {
	Afk_timeout int
	Joined_at   string
	// Afk_channel_id int `json:",string"`
	Id   int `json:",string"`
	Icon string
	Name string
	//	Roles          []Role
	Region string
	//Embed_channel_id int `json:",string"`
	//	Embed_channel_id string
	//	Embed_enabled    bool
	Owner_id int `json:",string"`
}

type Role struct {
	Permissions int
	Id          int `json:",string"`
	Name        string
}

type Channel struct {
	Guild_id        int `json:",string"`
	Id              int `json:",string"`
	Name            string
	Last_message_id string
	Is_private      string

	//	Permission_overwrites string
	//	Position              int `json:",string"`
	//	Type                  string
}

type Message struct {
	Attachments      []Attachment
	Tts              bool
	Embeds           []Embed
	Timestamp        string
	Mention_everyone bool
	Id               int `json:",string"`
	Edited_timestamp string
	Author           *Author
	Content          string
	Channel_id       int `json:",string"`
	Mentions         []Mention
}

type Mention struct {
}

type Attachment struct {
}

type Embed struct {
}

type Author struct {
	Username      string
	Discriminator int `json:",string"`
	Id            int `json:",string"`
	Avatar        string
}

// Create takes an email and password then prepares a RestClient with the given data,
// which is a simple object used for future requests.
func Create(email string, password string) (restClient *RestClient, err error) {
	if len(email) < 3 {
		err = errors.New("email too short")
		return
	}
	if len(password) < 3 {
		err = errors.New("password too short")
		return
	}
	session := &Session{"0", email, password, ""}
	httpClient := &http.Client{Timeout: (20 * time.Second)}
	restClient = &RestClient{discordUrl, session, httpClient, false}
	restClient.Session.Token, err = requestToken(restClient)
	if err != nil {
		return
	}
	restClient.Session.Id, err = requestSelf(restClient)
	if err != nil {
		return
	}
	return
}

// RequestToken asks the Rest server for a token by provided email/password
func requestToken(restClient *RestClient) (token string, err error) {

	if restClient == nil {
		err = errors.New("Empty restClient, Create() one first")
		return
	}

	if restClient.Session == nil || len(restClient.Session.Email) == 0 || len(restClient.Session.Password) == 0 {
		err = errors.New("Empty restClient.Session data, Create() to set email/password")
		return
	}

	var urlStr string = fmt.Sprintf("%s/%s", restClient.Url, "auth/login")
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(fmt.Sprintf(`{"email":"%s", "password":"%s"}`, restClient.Session.Email, restClient.Session.Password))))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := restClient.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}
	session := &Session{}

	err = json.Unmarshal(body, &session)
	token = session.Token
	return
}

// Identify user himself
func requestSelf(restClient *RestClient) (clientId string, err error) {

	body, err := Request(restClient, fmt.Sprintf("%s/%s", restClient.Url, "users/@me"))
	session := &Session{} // what's this for?
	err = json.Unmarshal(body, &session)
	clientId = session.Id
	return
}

func ListGuilds(restClient *RestClient) (guilds []Guild, err error) {

	body, err := Request(restClient, fmt.Sprintf("%s/%s", restClient.Url, fmt.Sprintf("users/%s/guilds", restClient.Session.Id)))
	err = json.Unmarshal(body, &guilds)

	return
}

func ListChannels(restClient *RestClient, guildId int) (channels []Channel, err error) {

	body, err := Request(restClient, fmt.Sprintf("%s/%s", restClient.Url, fmt.Sprintf("guilds/%d/channels", guildId)))
	err = json.Unmarshal(body, &channels)

	body, err = Request(restClient, fmt.Sprintf("%s/%s", restClient.Url, fmt.Sprintf("users/@me/channels", guildId)))
	err = json.Unmarshal(body, &channels)

	return
}

func GetMessages(restClient *RestClient, channelId int, before int, limit int) (messages []Message, err error) {
	// var urlStr = fmt.Sprintf("%s/%s", restClient.Url, fmt.Sprintf("channels/%d/messages?limit=%d&after=%d", channelId,limit,before))

	var urlStr = fmt.Sprintf("%s/%s", restClient.Url, fmt.Sprintf("channels/%d/messages", channelId))

	if limit > 0 {
		urlStr = urlStr + fmt.Sprintf("?limit=%d", limit)
	} else {
		urlStr = urlStr + "?limit=1"
	}

	if before > 0 {
		urlStr = urlStr + fmt.Sprintf("&after=%d", before)

	}

	body, err := Request(restClient, urlStr)
	err = json.Unmarshal(body, &messages)

	return
}

func CreateChannelUser(restClient *RestClient, userId int) (channelId int, err error) {

	var urlStr string = fmt.Sprintf("%s/%s", restClient.Url, fmt.Sprintf("users/@me/channels"))
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(fmt.Sprintf(`{"recipient_id":"%d"}`, userId))))
	if err != nil {
		return
	}
	req.Header.Set("authorization", restClient.Session.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := restClient.client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}

	// something something if debug
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		fmt.Print("JSON parse error: ", error)
		return
	}
	fmt.Println(urlStr+" Response:\n", string(prettyJSON.Bytes()))

	// err = json.Unmarshal(body, &responseMessage)
	return

	return
}

func SendMessage(restClient *RestClient, channelId int, message string) (responseMessage Message, err error) {
	var urlStr string = fmt.Sprintf("%s/%s", restClient.Url, fmt.Sprintf("channels/%d/messages", channelId))
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(fmt.Sprintf(`{"content":"%s"}`, message))))
	if err != nil {
		return
	}
	req.Header.Set("authorization", restClient.Session.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := restClient.client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode: %d, %s", resp.StatusCode, string(body)))
		return
	}

	// something something if debug
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		fmt.Print("JSON parse error: ", error)
		return
	}
	fmt.Println(urlStr+" Response:\n", string(prettyJSON.Bytes()))

	err = json.Unmarshal(body, &responseMessage)
	return
}

func Close(restClient *RestClient) (err error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", restClient.Url, fmt.Sprintf("auth/logout")), bytes.NewBuffer([]byte(fmt.Sprintf(``))))
	if err != nil {
		return
	}
	req.Header.Set("authorization", restClient.Session.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := restClient.client.Do(req)
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

func ReadLoop(restClient *RestClient, channelId int) {

	var lastMessageId int = 0

	var i int = 0

	for i < 1000 {

		messages, err := GetMessages(restClient, channelId, lastMessageId, 10)
		if err != nil {
			fmt.Println(err)
			return
		}

		var i int = len(messages) - 1
		// fmt.Println("loop ", i, " ", lastMessageId);

		if i > -1 { // seems poorly wrote..

			for i >= 0 {
				var message Message = messages[i]
				fmt.Println("\n", message.Id, ":", message.Timestamp, ":\n", message.Author.Username, " > ", message.Content)
				lastMessageId = message.Id
				i--
			}
		}

		time.Sleep(2000 * time.Millisecond)
		i++
	}
}

// Request makes a API GET Request.  This is a general purpose function
// and is used by all API functions.  It is exposed currently so it can
// also be used outside of this library.
func Request(restClient *RestClient, urlStr string) (body []byte, err error) {

	req, err := http.NewRequest("GET", urlStr, bytes.NewBuffer([]byte(fmt.Sprintf(``))))
	if err != nil {
		return
	}

	req.Header.Set("authorization", restClient.Session.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := restClient.client.Do(req)
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

	if restClient.Debug {
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
