// This is an example of using DiscordGo to obtain the
// authentication token for a given user account.
package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {

	// Check for Username and Password CLI arguments.
	if len(os.Args) != 3 {
		fmt.Println("You must provide username and password as arguments. See below example.")
		fmt.Println(os.Args[0], " [email] [password]")
		return
	}

	// Create a New Discord session and login with the provided
	// email and password.
	dg, err := discordgo.New(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Your Authentication Token is:\n\n%s\n", dg.Token)
}
