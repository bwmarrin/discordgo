package main

import (
	"fmt"
	"strings"
)

//
func admin(line string) (response string) {

	var err error

	// trim any leading or trailing space off the whole line
	line = strings.TrimSpace(line)

	// split the command from the rest
	split := strings.SplitN(line, " ", 2)

	// store the command and payload seperately
	command := strings.ToLower(split[0])
	command = strings.TrimPrefix(command, "~")

	var payload string = ""
	if len(split) > 1 {
		payload = split[1]
	}

	if command == "help" {
		response += fmt.Sprintln("`~help ...............` Display this help text")
		response += fmt.Sprintln("`~username [string]...` Set login username to [string]")
		response += fmt.Sprintln("`~password [string]...` Set login password to [string]")
		response += fmt.Sprintln("`~login ..............` Login to Discord")
		response += fmt.Sprintln("`~listen .............` Start websocket listener")
		response += fmt.Sprintln("`~logout .............` Logout from Discord")
		return
	}

	if command == "username" {
		Username = payload
		response += "Done."
		return
	}

	if command == "password" {
		Password = payload
		response += "Done."
		return
	}

	if command == "login" {
		Session.Token, err = Session.Login(Username, Password)
		if err != nil {
			fmt.Println("Unable to login to Discord.")
			fmt.Println(err)
		}
		response += "Done."
		return
	}

	if command == "listen" {

		// open connection
		err = Session.Open()
		if err != nil {
			fmt.Println(err)
		}

		// Do Handshake? (dumb name)
		err = Session.Handshake()
		if err != nil {
			fmt.Println(err)
		}

		// Now listen for events / messages
		go Session.Listen()
		response += "Done."
		return
	}

	if command == "logout" {
		err = Session.Logout()
		if err != nil {
			fmt.Println("Unable to logout from Discord.")
			fmt.Println(err)
		}
		response += "Done."
		return
	}

	response += "I'm sorry I don't understand that command.  Try ~help"
	return
}
