package main

import (
	"flag"
	dgo "github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	botToken  = flag.String("bot_tok", "", "Bot token")
	testGuild = flag.String("test_guild", "", "Guild for testing the bot")
)

func setup() *dgo.Session {
	client, err := dgo.New("Bot " + *botToken)
	if err != nil {
		log.Fatal(err)
	}

	client.AddHandler(func(*dgo.Session, *dgo.Ready) {
		log.Println("Bot is up!")
	})

	return client
}

func main() {
	flag.Parse()

	client := setup()

	client.AddHandler(func(s *dgo.Session, i *dgo.InteractionCreate) {
		responseType := dgo.InteractionResponseType(i.Interaction.Data.Options[0].Value.(float64))
		responseData := &dgo.InteractionApplicationCommandResponseData{
			TTS:     false,
			Content: "here we go",
			// Flags: 1 << 6,
		}
		if responseType == dgo.InteractionResponseACKWithSource || responseType == dgo.InteractionResponseAcknowledge {
			responseData = nil
		}
		log.Println("response parameters:", responseType, responseData)
		err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
			Type: responseType,
			Data: responseData,
		})
		log.Println("response", err)
		time.Sleep(time.Second * 2)
		err = s.InteractionResponseEdit("", i.Interaction, &dgo.WebhookEdit{
			Content: "here we go!",
		})
		log.Println("response edit", err)
		time.Sleep(time.Second * 2)
		err = s.InteractionResponseDelete("", i.Interaction)
		log.Println("response delete", err)
		err = s.InteractionResponseDelete("", i.Interaction)
		log.Println("response delete 2", err)

		followupMessage, err := s.FollowupMessageCreate("", i.Interaction, true, &dgo.WebhookParams{
			Content: "followup messages rule!",
		})
		log.Println("followup message create", followupMessage, err)
		time.Sleep(time.Second * 3)
		err = s.FollowupMessageEdit("", i.Interaction, followupMessage.ID, &dgo.WebhookEdit{
			Content: "that's true",
		})
		log.Println("followup message edit", err)
		time.Sleep(time.Second * 3)
		err = s.FollowupMessageDelete("", i.Interaction, followupMessage.ID)
		log.Println("followup message delete", err)
	})

	if err := client.Open(); err != nil {
		log.Fatal(err)
	}

	cmd, err := client.ApplicationCommandCreate("", &dgo.ApplicationCommand{
		Name:        "test-slashes",
		Description: "Command for testing application commands",
		Options: []*dgo.ApplicationCommandOption{
			{
				Type:        dgo.ApplicationCommandOptionInteger,
				Name:        "typ",
				Description: "Response type",
				Choices: []*dgo.ApplicationCommandOptionChoice{
					{
						Name:  "ACK",
						Value: dgo.InteractionResponseAcknowledge,
					},
					{
						Name:  "ACK with source",
						Value: dgo.InteractionResponseACKWithSource,
					},
					{
						Name:  "Channel message",
						Value: dgo.InteractionResponseChannelMessage,
					},
					{
						Name:  "Channel message with source",
						Value: dgo.InteractionResponseChannelMessageWithSource,
					},
				},
				Required: false,
			},
		},
	}, *testGuild)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Old command: ", cmd)
	log.Println("Created command: ", cmd)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, os.Interrupt)
	<-ch
	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
