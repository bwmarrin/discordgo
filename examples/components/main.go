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
										Style:    discordgo.ButtonSuccess,
										Disabled: false,
										CustomID: "yes_btn",
									},
									discordgo.Button{
										Label:    "No",
										Style:    discordgo.ButtonDanger,
										Disabled: false,
										CustomID: "no_btn",
									},
									discordgo.Button{
										Label:    "I don't know",
										Style:    discordgo.ButtonLink,
										Disabled: false,
										// Link buttons don't require CustomID and do not trigger the gateway/HTTP event
										URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
										Emoji: discordgo.ComponentEmoji{
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
										Style:    discordgo.ButtonLink,
										Disabled: false,
										URL:      "https://discord.gg/discord-developers",
									},
								},
							},
							// If a select menu is used in an action row, buttons cannot be used in the same row,
							// and only one select menu is allowed per row.
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.SelectMenu{
										CustomID:    "select_menu_rating",
										Placeholder: "Or do you like Select Menus?",
										Options: []discordgo.SelectOption{
											{
												Label:       "Yes",
												Value:       "select_menu_yes",
												Description: "I like them more than buttons",
											},
											{
												Label:       "I like both",
												Value:       "select_menu_both",
												Description: "Both select menus and buttons are good",
											},
											{
												Label:       "No",
												Value:       "select_menu_no",
												Description: "Buttons for life",
											},
										},
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

		// Values contain the values currently selected from the select menu component.
		if len(i.MessageComponentData().Values) > 0 {
			// The select menu is set to only allow one value to be selected.
			// It can be configured to allow more than one value to be selected at a time.
			switch i.MessageComponentData().Values[0] {
			case "select_menu_yes":
				content += "(yes to select menus)"
			case "select_menu_both":
				content += "(likes both buttons and select menus)"
			case "select_menu_no":
				content += "(no to select menus)"
			}
		} else {
			// CustomID field contains the same id as when was sent. It's used to identify the which button was clicked.
			switch i.MessageComponentData().CustomID {
			case "yes_btn":
				content += "(yes to buttons)"
			case "no_btn":
				content += "(no to buttons)"
			}

		}

		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Buttons also may update the message which to which they are attached.
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
								Style:    discordgo.ButtonLink,
								Disabled: false,
								URL:      "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
								Emoji: discordgo.ComponentEmoji{
									Name: "💠",
								},
							},
						},
					},
				},
			},
		})
		if e != nil {
			panic(e)
		}
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
