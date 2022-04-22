package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally.")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", false, "Weather to remove all commands after shutdowning or not.")
	DoInBulk       = flag.Bool("bulk", true, `Weather to create and delete commands in bulk or not. 
		Setting to "true" avoids rate limits of the API but may complicate debugging.`)
	// Note that deleteCommandsOnFailure() will always delete commands in bulk
)

func deleteCommands(s *discordgo.Session, guildId string) {
	for _, cmd := range commands {
		err := s.ApplicationCommandDelete(
			s.State.User.ID,
			guildId,
			cmd.ID,
		)
		if err != nil {
			log.Fatalf("Could not delete command %s: %v\n", cmd.Name, err)
		}
	}
}

func createCommands(s *discordgo.Session, guildId string) {
	defer deleteCommandsOnFailure(s, guildId)

	var err error
	for i, cmd := range commands {
		commands[i], err = s.ApplicationCommandCreate(
			s.State.User.ID,
			guildId,
			cmd,
		)
		if err != nil {
			log.Panicf("Could not create command %s: %s\n", cmd.Name, err)
		}
	}
}

// Discord API has rate limit which limits the amount of request bot can make per time period
// Thus if there are to many commands creation or deletion of them on by one may take a lot of time
// due to said limits.
func createCommandsBulk(s *discordgo.Session, guildId string) {
	defer deleteCommandsOnFailure(s, guildId)

	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildId, commands)
	if err != nil {
		log.Panicf("Could not create commands in bulk")
	}
}

func deleteCommandsBulk(s *discordgo.Session, guildId string) {
	empty := make([]*discordgo.ApplicationCommand, 0)
	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildId, empty)
	if err != nil {
		log.Fatalf("Could not delete bulk of commands %s\n", err)
	}
}

// Commands are almost gauranteed to be deleted on failure so that
// they are not left hanging with nobody to handle them
func deleteCommandsOnFailure(s *discordgo.Session, guildID string) {
	err := recover()
	if err != nil && *RemoveCommands {
		deleteCommandsBulk(s, guildID)
		log.Fatalf("Error occured, commands were deleted: %s\n", err)
	}
}

func main() {

	flag.Parse()

	s, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(HandleInteraction)
	s.AddHandler(Login)

	err = s.Open()
	if err != nil {
		log.Fatalf("Could not open the session: %v", err)
	}
	defer s.Close()

	createCmd := createCommandsBulk
	deleteCmd := deleteCommandsBulk
	if !*DoInBulk {
		createCmd = createCommands
		deleteCmd = deleteCommands
	}

	log.Println("Adding commands...")
	createCmd(s, *GuildID)
	if *RemoveCommands {
		defer deleteCmd(s, *GuildID)
		defer log.Println("Removing commands...")
	}

	log.Println("Bot is up and running")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
}
