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
