package discordgo_test

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func ExampleApplication() {

	// Authentication Token pulled from environment variable DG_TOKEN
	Token := os.Getenv("DG_TOKEN")
	if Token == "" {
		return
	}

	// Create a new Discordgo session
	dg, err := discordgo.New(Token)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create an new Application
	ap := &discordgo.Application{}
	ap.Name = "TestApp"
	ap.Description = "TestDesc"
	ap, err = dg.ApplicationCreate(ap)
	fmt.Printf("ApplicationCreate: err: %+v, app: %+v\n", err, ap)

	// Get a specific Application by it's ID
	ap, err = dg.Application(ap.ID)
	fmt.Printf("Application: err: %+v, app: %+v\n", err, ap)

	// Update an existing Application with new values
	ap.Description = "Whooooa"
	ap, err = dg.ApplicationUpdate(ap.ID, ap)
	fmt.Printf("ApplicationUpdate: err: %+v, app: %+v\n", err, ap)

	// create a new bot account for this application
	bot, err := dg.ApplicationBotCreate(ap.ID, "")
	fmt.Printf("BotCreate: err: %+v, bot: %+v\n", err, bot)

	// Get a list of all applications for the authenticated user
	apps, err := dg.Applications()
	fmt.Printf("Applications: err: %+v, apps : %+v\n", err, apps)
	for k, v := range apps {
		fmt.Printf("Applications: %d : %+v\n", k, v)
	}

	// Delete the application we created.
	err = ap.Delete()
	fmt.Printf("Delete: err: %+v\n", err)

	return
}

// This provides an example on converting an existing normal user account
// into a bot account.  You must authentication to Discord using your personal
// username and password then provide the authentication token of the account
// you want converted.
func ExampleApplicationConvertBot() {

	dg, err := discordgo.New("myemail", "mypassword")
	if err != nil {
		fmt.Println(err)
		return
	}

	// create an application
	ap := &discordgo.Application{}
	ap.Name = "Application Name"
	ap.Description = "Application Description"
	ap, err = dg.ApplicationCreate(ap)
	fmt.Printf("ApplicationCreate: err: %+v, app: %+v\n", err, ap)

	// create a bot account
	bot, err := dg.ApplicationBotCreate(ap.ID, "existing bot user account token")
	fmt.Printf("BotCreate: err: %+v, bot: %+v\n", err, bot)

	if err != nil {
		fmt.Printf("You can not login with your converted bot user using the below token\n%s\n", bot.Token)
	}

	return
}
