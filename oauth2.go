// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains functions related to Discord OAuth2 applications

package discordgo

import (
	"fmt"
)

// An Application struct stores values for a Discord OAuth2 Application
type Application struct {
	ID           string    `json:"id,omitempty"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Icon         string    `json:"icon,omitempty"`
	Secret       string    `json:"secret,omitempty"`
	RedirectURIs *[]string `json:"redirect_uris,omitempty"`

	// Concept.. almost guarenteed to be removed.
	// Imagine that it's just not even here at all.
	ses *Session
}

// Application returns an Application structure of a specific Application
//   appID : The ID of an Application
func (s *Session) Application(appID string) (st *Application, err error) {

	body, err := s.Request("GET", APPLICATION(appID), nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	st.ses = s
	return
}

// Applications returns all applications for the authenticated user
func (s *Session) Applications() (st []*Application, err error) {

	body, err := s.Request("GET", APPLICATIONS, nil)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	for k, _ := range st {
		st[k].ses = s
	}
	return
	// TODO ..
}

// ApplicationCreate creates a new Application
//    name : Name of Application / Bot
//    uris : Redirect URIs (Not required)
func (s *Session) ApplicationCreate(ap *Application) (st *Application, err error) {

	data := struct {
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		RedirectURIs *[]string `json:"redirect_uris,omitempty"`
	}{ap.Name, ap.Description, ap.RedirectURIs}

	body, err := s.Request("POST", APPLICATIONS, data)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	st.ses = s
	return
}

// ApplicationEdit edits an existing Application
//   var : desc
func (s *Session) ApplicationUpdate(appID string, ap *Application) (st *Application, err error) {

	data := struct {
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		RedirectURIs *[]string `json:"redirect_uris,omitempty"`
	}{ap.Name, ap.Description, ap.RedirectURIs}

	body, err := s.Request("PUT", APPLICATION(appID), data)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	st.ses = s
	return
}

// ApplicationDelete deletes an existing Application
//   appID : The ID of an Application
func (s *Session) ApplicationDelete(appID string) (err error) {

	_, err = s.Request("DELETE", APPLICATION(appID), nil)
	if err != nil {
		return
	}

	return
}

//////////////////////////////////////////////////////////////////////////////
// Below two functions are experimental ideas, they will absolutely change
// one way or another and may be deleted entirely.

// Delete is a concept helper function, may be removed.
// this func depends on the Application.ses pointer
// pointing to the Discord session that the application
// came from.  This "magic" makes some very very nice helper
// functions possible.
func (a *Application) Delete() (err error) {
	if a.ses == nil {
		return fmt.Errorf("ses is nil.")
	}
	return a.ses.ApplicationDelete(a.ID)
}

// Delete is a concept helper function, may be removed.
// this one doesn't depend on the "magic" of adding the ses
// pointer to each Application
func (a *Application) DeleteB(s *Session) (err error) {
	return s.ApplicationDelete(a.ID)
}
