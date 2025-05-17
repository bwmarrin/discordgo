package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!components" {
		legacyComponentsExample(s, m)
		actionRowExample(s, m)
		sectionExample(s, m)
		fileExample(s, m)
		// Separator example is visual and best shown with other components.
		// Container example wraps other components.
		combinedExample(s, m)
		mediaGalleryExample(s, m)
	}
}

func legacyComponentsExample(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "This is a message with legacy components",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Click Me",
						Style:    discordgo.PrimaryButton,
						CustomID: "legacy_button_click",
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error sending legacy components example message: %v", err)
	}
}

func actionRowExample(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Flags: discordgo.MessageFlagsIsComponentsV2,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Accept",
						Style:    discordgo.PrimaryButton,
						CustomID: "ar_button_accept",
					},
					discordgo.Button{
						Label: "Learn More",
						Style: discordgo.LinkButton,
						URL:   "http://watchanimeattheoffice.com/",
					},
					discordgo.Button{
						Label:    "Decline",
						Style:    discordgo.DangerButton,
						CustomID: "ar_button_decline",
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error sending action row example message: %v", err)
	}
}

func sectionExample(s *discordgo.Session, m *discordgo.MessageCreate) {
	id := 123 // Example ID
	accessoryThumbnail := discordgo.Thumbnail{
		Media:       discordgo.UnfurledMediaItem{URL: "https://example.com/thumbnail.png"}, // Replace with actual URL or attachment
		Description: "This is a thumbnail description",
	}
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Flags: discordgo.MessageFlagsIsComponentsV2,
		Components: []discordgo.MessageComponent{
			discordgo.Section{
				CustomID: &id,
				Components: []discordgo.MessageComponent{
					discordgo.TextDisplay{
						Content: "# Real Game v7.3",
					},
					discordgo.TextDisplay{
						Content: "This is a section with text and an accessory thumbnail.",
					},
				},
				Accessory: &accessoryThumbnail,
			},
		},
	})
	if err != nil {
		log.Printf("Error sending section example message: %v", err)
	}
}

func fileExample(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Note: To send actual files, you'd use the Files field of MessageSend
	// and refer to them by attachment://filename.ext in UnfurledMediaItem.
	// This example only demonstrates the component structure.
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Flags: discordgo.MessageFlagsIsComponentsV2,
		Components: []discordgo.MessageComponent{
			discordgo.TextDisplay{
				Content: "# New game version released for testing!",
			},
			discordgo.TextDisplay{
				Content: "Grab the game here:",
			},
			discordgo.FileComponentData{
				File: discordgo.UnfurledMediaItem{
					URL: "attachment://game.zip", // Assumes a file named game.zip is attached
				},
			},
			discordgo.TextDisplay{
				Content: "Latest manual artwork here:",
			},
			discordgo.FileComponentData{
				File: discordgo.UnfurledMediaItem{
					URL: "attachment://manual.pdf", // Assumes a file named manual.pdf is attached
				},
			},
		},
		// Example of how you would attach files (content of files not included for brevity)
		// Files: []*discordgo.File{
		// {Name: "game.zip", Reader: ...},
		// {Name: "manual.pdf", Reader: ...},
		// },
	})
	if err != nil {
		log.Printf("Error sending file example message: %v", err)
	}
}

func combinedExample(s *discordgo.Session, m *discordgo.MessageCreate) {
	containerColor := 0x7289DA // Discord Blurple
	containerBorder := true

	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Flags: discordgo.MessageFlagsIsComponentsV2,
		Components: []discordgo.MessageComponent{
			discordgo.Container{
				Color:  &containerColor,
				Border: &containerBorder,
				Components: []discordgo.MessageComponent{
					discordgo.TextDisplay{
						Content: "This is text inside a container.",
					},
					discordgo.Separator{},
					discordgo.Section{
						Components: []discordgo.MessageComponent{
							discordgo.TextDisplay{Content: "Section within Container"},
						},
						Accessory: &discordgo.Button{
							Label:    "Button in Section",
							Style:    discordgo.SuccessButton,
							CustomID: "container_section_button",
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Button in ActionRow in Container",
								Style:    discordgo.PrimaryButton,
								CustomID: "container_ar_button",
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error sending combined V2 components example message: %v", err)
	}
}

func mediaGalleryExample(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Flags: discordgo.MessageFlagsIsComponentsV2,
		Components: []discordgo.MessageComponent{
			discordgo.TextDisplay{
				Content: "Live webcam shots as of 18-04-2025 at 12:00 UTC",
			},
			discordgo.MediaGallery{
				Items: []discordgo.MediaGalleryItem{
					{
						Media:       discordgo.UnfurledMediaItem{URL: "https://livevideofeedconvertedtoimage/webcam1.png"},
						Description: "An aerial view looking down on older industrial complex buildings. The main building is white with many windows and pipes running up the walls.",
					},
					{
						Media:       discordgo.UnfurledMediaItem{URL: "https://livevideofeedconvertedtoimage/webcam2.png"},
						Description: "An aerial view of old broken buildings. Nature has begun to take root in the rooftops. A portion of the middle building's roof has collapsed inward. In the distant haze you can make out a far away city.",
					},
					{
						Media:       discordgo.UnfurledMediaItem{URL: "https://livevideofeedconvertedtoimage/webcam3.png"},
						Description: "A street view of a downtown city. Prominently in photo are skyscrapers and a domed building",
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error sending media gallery example message: %v", err)
	}
}

func init() {
	// Simple way to get token from env. Replace with your preferred method.
	Token = os.Getenv("DISCORD_BOT_TOKEN")
	if Token == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable not set.")
	}
}
