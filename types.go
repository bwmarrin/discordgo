// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains custom types, currently only a timestamp wrapper.

package discordgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Timestamp stores a timestamp, as sent by the Discord API.
type Timestamp string

// Parse parses a timestamp string into a time.Time object.
// The only time this can fail is if Discord changes their timestamp format.
func (t Timestamp) Parse() (time.Time, error) {
	return time.Parse(time.RFC3339, string(t))
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
	return fmt.Sprintf("HTTP %s, %s", r.Response.Status, r.ResponseBody)
}

// IDSlice Is a slice of snowflake id's that properly marshals and unmarshals the way discord expects them to
// They unmarshal from string arrays and marshals back to string arrays
type IDSlice []int64

func (ids *IDSlice) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return nil
	}

	// Split and strip away "[" "]"
	split := strings.Split(string(data[1:len(data)-1]), ",")
	*ids = make([]int64, len(split))
	for i, s := range split {
		s = strings.TrimSpace(s)
		if len(s) < 3 {
			// Empty or invalid
			continue
		}

		// Strip away quotes and parse
		parsed, err := strconv.ParseInt(s[1:len(s)-1], 10, 64)
		if err != nil {
			return err
		}
		(*ids)[i] = parsed
	}

	return nil
}

func (ids IDSlice) MarshalJSON() ([]byte, error) {
	// Capacity:
	// 2 brackets
	// each id is:
	//    18 characters currently, but 1 extra added for the future,
	//    1 comma
	//    2 quotes
	if len(ids) < 1 {
		return []byte("[]"), nil
	}

	outPut := make([]byte, 1, 2+(len(ids)*22))
	outPut[0] = '['

	for i, id := range ids {
		if i != 0 {
			outPut = append(outPut, '"', ',', '"')
		} else {
			outPut = append(outPut, '"')
		}
		outPut = append(outPut, []byte(strconv.FormatInt(id, 10))...)
	}

	outPut = append(outPut, '"', ']')
	return outPut, nil
}
