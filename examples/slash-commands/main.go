package main

import (
	dgo "github.com/bwmarrin/discordgo"
	dotenv "github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func loadEnv() {
	err := dotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	loadEnv()
	client, err := dgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	client.AddHandler(func(*dgo.Session, *dgo.Ready) {
		log.Println("Bot is up!")
	})
	client.AddHandler(func(_ *dgo.Session, i *dgo.InteractionCreate) {
		log.Println("Received an interaction: ", i.Interaction)
	})

	if err := client.Open(); err != nil {
		log.Fatal(err)
	}

	cmd, err := client.ApplicationCommandCreate(&dgo.ApplicationCommand{
		Name:        "pingme",
		Description: "Command for pinging a bot",
		Options: []*dgo.ApplicationCommandOption{
			{
				Type:        dgo.ApplicationCommandOptionBoolean,
				Name:        "showcmd",
				Description: "Show the command in the chat or hide it",
				Default:     false,
				Required:    false,
			},
		},
	}, "")

	if err != nil {
		panic(err)
	}

	log.Println("Created command ID: ", cmd.ID)
	log.Println("Old command: ", cmd)
	cmd, err = client.ApplicationCommandEdit(cmd.ID, "", &dgo.ApplicationCommand{
		Name:        "pingme_twice",
		Description: "Twice",
	})

	if err != nil {
		panic(err)
	}

	log.Println("New command: ", cmd)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, os.Interrupt)
	<-ch
	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
