// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains custom types, currently only a timestamp wrapper.

package discordgo

import (
	"encoding/json"
	"net/http"
	"time"
)

// Timestamp stores a timestamp, as sent by the Discord API.
type Timestamp string

// Parse parses a timestamp string into a time.Time object.
// The only time this can fail is if Discord changes their timestamp format.
func (t Timestamp) Parse() (time.Time, error) {
	return time.Parse(time.RFC3339, string(t))
}

// NullString is a custom String type to allow for optional modifying strings that can be reset
// the JSON Marshalling encoding is unique:
// if the string is empty, it will be encoded as null
// if the string is not empty, it will be encoded as a string
type NullString string

// MarshalJSON returns the JSON encoding of the NullString
func (s NullString) MarshalJSON() ([]byte, error) {
	if s == "" {
		return json.Marshal(nil)
	}

	return json.Marshal(string(s))
}

// NullBool is a custom Bool type to allow for optional modifying booleans that can have false as value
// the JSON Marshalling encoding is unique:
// if the bool is false, it will be encoded as false
// if the bool is true, it will be encoded as true
type NullBool bool

// MarshalJSON returns the JSON encoding of the NullBool
func (b NullBool) MarshalJSON() ([]byte, error) {
	if b == false {
		return json.Marshal(false)
	}

	return json.Marshal(true)
}

// NullInt is a custom Int type to allow for optional modifying integers that can have 0 as value
// the JSON Marshalling encoding is unique:
// if the integer is 0, it will be encoded as 0
// if the integer is not equal to 0, it will be encoded as an integer
type NullInt int

// MarshalJSON returns the JSON encoding of the NullBool
func (i NullInt) MarshalJSON() ([]byte, error) {
	if i == 0 {
		return json.Marshal(0)
	}

	return json.Marshal(int(i))
}

// RESTError stores error information about a request with a bad response code.
// Message is not always present, there are cases where api calls can fail
// without returning a json message.
type RESTError struct {
	Request      *http.Request
	Response     *http.Response
	ResponseBody []byte

	Message *APIErrorMessage // Message may be nil.
}

func newRestError(req *http.Request, resp *http.Response, body []byte) *RESTError {
	restErr := &RESTError{
		Request:      req,
		Response:     resp,
		ResponseBody: body,
	}

	// Attempt to decode the error and assume no message was provided if it fails
	var msg *APIErrorMessage
	err := json.Unmarshal(body, &msg)
	if err == nil {
		restErr.Message = msg
	}

	return restErr
}

func (r RESTError) Error() string {
	return "HTTP " + r.Response.Status + ", " + string(r.ResponseBody)
}
