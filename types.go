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

type sleepCT struct {
	d     time.Duration // desired duration between targets
	t     time.Time     // last time target
	wake  time.Time     // last wake time
	drift int64         // last wake drift microseconds
}

func newSleepCT(d time.Duration) sleepCT {

	s := sleepCT{}

	s.d = d
	s.t = time.Now()

	return s
}

func (s *sleepCT) sleepNext() int64 {

	now := time.Now()

	// if target is zero safety net
	if s.t.IsZero() {
		s.t = now.Add(-s.d)
	}

	// Move forward the sleep target by the duration
	s.t = s.t.Add(s.d)

	// Compute the desired sleep time to reach the target
	d := time.Until(s.t)

	// Sleep
	time.Sleep(d)

	// record the wake time
	s.wake = time.Now()
	s.drift = s.wake.Sub(s.t).Microseconds()

	// return the drift for monitoring purposes
	return s.drift
}
