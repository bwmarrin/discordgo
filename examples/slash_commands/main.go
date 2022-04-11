package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Weather to remove all commands after shutdowning or not")
)

func deleteCommands(ds *discordgo.Session, guildId string) {
	for _, cmd := range commands {
		err := ds.ApplicationCommandDelete(
			ds.State.User.ID,
			guildId,
			cmd.ID,
		)
		if err != nil {
			log.Fatalf("Could not delete command %s: %v\n", cmd.Name, err)
		}
	}
}

func createCommands(ds *discordgo.Session, guildId string) {
	var err error
	for i, cmd := range commands {
		commands[i], err = ds.ApplicationCommandCreate(
			ds.State.User.ID,
			guildId,
			cmd,
		)
		if err != nil {
			log.Panicf("Could not create command %s: %s\n", cmd.Name, err)
		}
	}
}

func main() {

	flag.Parse()

	*BotToken = "OTE1NjE2NDcyMTU4NTg0ODYz.YaeMSg.vr_uULSBq7PL1iqwAe8qwOm_EgM"
	*GuildID = "891737479412060191"
	*RemoveCommands = true

	ds, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	ds.AddHandler(HandleInteraction)
	ds.AddHandler(Login)

	err = ds.Open()
	if err != nil {
		log.Fatalf("Could not open the session: %v", err)
	}
	defer ds.Close()

	log.Println("Adding commands...")
	createCommands(ds, *GuildID)
	if *RemoveCommands {
		defer deleteCommands(ds, *GuildID)
		defer log.Println("Removing commands...")
	}

	log.Println("Bot is up and running")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
}
