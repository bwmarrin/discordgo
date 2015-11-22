/******************************************************************************
 * A Discord API for Golang.
 *
 * This file has structs and functions specific to a session.
 *
 * A session is a single connection to Discord for a given
 * user and all REST and Websock API functions exist within
 * a session.
 *
 * See the restapi.go and wsapi.go for more information.
 */

package discordgo

import (
	"net"

	"github.com/gorilla/websocket"
)
