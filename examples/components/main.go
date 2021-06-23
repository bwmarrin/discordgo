package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID  = flag.String("guild", "", "Test guild ID")
	BotToken = flag.String("token", "", "Bot access token")
	AppID    = flag.String("app", "", "Application ID")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	// Buttons are part of interactions, so we register InteractionCreate handler
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if i.ApplicationCommandData().Name == "feedback" {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Are you satisfied with Buttons?",
						// Buttons and other components are specified in Components field.
						Components: []discordgo.MessageComponent{
							// ActionRow is a container of all buttons within the same row.
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.Button{
										Label:    "Yes",
										Style:    discordgo.SuccessButton,
										Disabled: false,
										CustomID: "yes_btn",
									},
									discordgo.Button{
										Label:    "No",
										Style:    discordgo.DangerButton,
										Disabled: false,
										CustomID: "no_btn",
									},
									discordgo.Button{
										Label:    "I don't know",
										Style:    discordgo.LinkButton,
										Disabled: false,
										// Link buttons doesn't require CustomID and does not trigger the gateway/HTTP event
										URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
										Emoji: discordgo.ButtonEmoji{
											Name: "🤷",
										},
									},
								},
							},
							// The message may have multiple actions rows.
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.Button{
										Label:    "Discord Developers server",
										Style:    discordgo.LinkButton,
										Disabled: false,
										URL:      "https://discord.gg/discord-developers",
									},
								},
							},
						},
					},
				})
				if err != nil {
					panic(err)
				}
			}
			return
		}
		// Type for button press will be always InteractionButton (3)
		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		content := "Thanks for your feedback "

		// CustomID field contains the same id as when was sent. It's used to identify the which button was clicked.
		switch i.MessageComponentData().CustomID {
		case "yes_btn":
			content += "(yes)"
		case "no_btn":
			content += "(no)"
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Buttons also may update the message which they was attached to.
			// Or may just acknowledge (InteractionResponseDeferredMessageUpdate) that the event was received and not update the message.
			// To update it later you need to use interaction response edit endpoint.
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: content,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Our sponsor",
								Style:    discordgo.LinkButton,
								Disabled: false,
								URL:      "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
								Emoji: discordgo.ButtonEmoji{
									Name: "💠",
								},
							},
						},
					},
				},
			},
		})
	})
	_, err := s.ApplicationCommandCreate(*AppID, *GuildID, &discordgo.ApplicationCommand{
		Name:        "feedback",
		Description: "Give your feedback",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}
