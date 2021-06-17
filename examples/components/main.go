package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"sort"
	"sync"
)

// Bot parameters
var (
	GuildID   = flag.String("guild", "", "Test guild ID")
	ChannelID = flag.String("channel", "", "Test channel ID")
	BotToken  = flag.String("token", "", "Bot access token")
	AppID     = flag.String("app", "", "Application ID")
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

type MessageVoteStats struct {
	PeopleVoted map[string]string
	Votes       map[string]int
	Author      string
}

var pollVotes map[string]MessageVoteStats
var pollVotesMtx sync.RWMutex

func constructVotesEmbed(stats MessageVoteStats) *discordgo.MessageEmbed {
	var votes []struct {
		label   string
		count   int
		percent float32
	}
	if len(stats.PeopleVoted) == 0 {
		return &discordgo.MessageEmbed{
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Go", Value: "Percentage: 0.0%\nPeople voted: 0"},
				{Name: "JS", Value: "Percentage: 0.0%\nPeople voted: 0"},
				{Name: "Python", Value: "Percentage: 0.0%\nPeople voted: 0"},
			},
			Color:  0xFFA1F0,
			Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("%d people voted", len(stats.PeopleVoted))},
		}
	}
	for k, v := range stats.Votes {
		votes = append(votes, struct {
			label   string
			count   int
			percent float32
		}{label: k, count: v, percent: (100.0 / float32(len(stats.PeopleVoted))) * float32(v)})
	}
	sort.Slice(votes, func(i, j int) bool { return votes[i].count > votes[j].count })

	var fields []*discordgo.MessageEmbedField

	for _, v := range votes {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  v.label,
			Value: fmt.Sprintf("Percentage: %.1f%%\nPeople voted: %d", v.percent, v.count),
		})
	}

	return &discordgo.MessageEmbed{
		Fields: fields,
		Color:  0xFFA1F0,
		Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("%d people voted", len(stats.PeopleVoted))},
	}
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
		pollVotesMtx.Lock()
		defer pollVotesMtx.Unlock()
		pollVotes = make(map[string]MessageVoteStats)
	})
	// Buttons are part of interactions, so we register InteractionCreate handler
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			var data *discordgo.InteractionResponseData
			switch i.ApplicationCommandData().Name {
			case "buttons":
				data = &discordgo.InteractionResponseData{
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
									CustomID: "yes",
								},
								discordgo.Button{
									Label:    "No",
									Style:    discordgo.DangerButton,
									Disabled: false,
									CustomID: "no",
								},
								discordgo.Button{
									Label:    "I don't know",
									Style:    discordgo.LinkButton,
									Disabled: false,
									// Link buttons doesn't require CustomID and does not trigger the gateway/HTTP event
									Link: "https://discord.dev/interactions/message-components",
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
									Label:    "Ask the question in #buttons on Discord Developers server",
									Style:    discordgo.LinkButton,
									Disabled: false,
									Link:     "https://discord.gg/discord-developers",
								},
							},
						},
					},
				}
			case "poll":
				data = &discordgo.InteractionResponseData{
					TTS:     false,
					Content: "What's your favorite most beloved programming languages?",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.SelectMenu{
									CustomID:    "choice",
									Placeholder: "Select your favorite language",
									Options: []discordgo.SelectMenuOption{
										{
											Label:       "Go",
											Value:       "Go",
											Description: "Go programming language.",
											Emoji:       discordgo.ComponentEmoji{},
											Default:     false,
										},
										{
											Label:       "JS",
											Value:       "JS",
											Description: "JavaScript programming language.",
											Emoji:       discordgo.ComponentEmoji{},
											Default:     false,
										},
										{
											Label:       "Python",
											Value:       "Python",
											Description: "Python programming language.",
											Emoji:       discordgo.ComponentEmoji{Name: ""},
											Default:     false,
										},
									},
								},
							},
						},
					},
					Flags: 0,
				}
			}
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: data,
			})
			if err != nil {
				panic(err)
			}
			if i.ApplicationCommandData().Name == "poll" {
				msg, err := s.InteractionResponse(*AppID, i.Interaction)
				if err != nil {
					panic(err)
				}
				pollVotesMtx.Lock()
				defer pollVotesMtx.Unlock()
				pollVotes[msg.ID] = MessageVoteStats{
					Votes:       map[string]int{"Go": 0, "JS": 0, "Python": 0},
					Author:      i.Member.User.ID,
					PeopleVoted: make(map[string]string),
				}
				fmt.Println(pollVotes[msg.ID])
			}
			return
		}
		// Type for all components will be always InteractionMessageComponent
		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		// CustomID field contains the same id as when was sent. It's used to identify the which button was clicked.
		switch i.MessageComponentData().CustomID {
		case "yes", "no":
			content := "Thanks for your feedback " + i.MessageComponentData().CustomID
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Buttons also may update the message which they was attached to.
				// Or may just acknowledge (InteractionResponseDeferredMessageUpdate) that the event was received and not update the message.
				// To update it later you need to use interaction response edit endpoint.
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					TTS:     false,
					Content: content,
					Flags:   1 << 6, // Ephemeral message
				},
			})
		case "end_poll":
			if i.Member.User.ID != pollVotes[i.Message.ID].Author {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You didn't start this poll, so you can't end it",
						Flags:   1 << 6,
					},
				})
				return
			}
			pollVotesMtx.Lock()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds:     []*discordgo.MessageEmbed{constructVotesEmbed(pollVotes[i.Message.ID])},
					Components: []discordgo.MessageComponent{},
				},
			})
			delete(pollVotes, i.Message.ID)
			pollVotesMtx.Unlock()

		case "choice":
			pollVotesMtx.Lock()
			defer pollVotesMtx.Unlock()
			stats := pollVotes[i.Message.ID]
			fmt.Println(stats)
			if len(i.MessageComponentData().Values) == 0 {
				stats.Votes[stats.PeopleVoted[i.Member.User.ID]]--
				delete(stats.PeopleVoted, i.Member.User.ID)
				pollVotes[i.Message.ID] = stats
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							constructVotesEmbed(pollVotes[i.Message.ID]),
						},
						Components: []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.SelectMenu{
										CustomID:    "choice",
										Placeholder: "Select your favorite language",
										Options: []discordgo.SelectMenuOption{
											{
												Label:       "Go",
												Value:       "Go",
												Description: "Go programming language.",
												Emoji:       discordgo.ComponentEmoji{},
												Default:     false,
											},
											{
												Label:       "JS",
												Value:       "javascript",
												Description: "JavaScript programming language.",
												Emoji:       discordgo.ComponentEmoji{},
												Default:     false,
											},
											{
												Label:       "Python",
												Value:       "py",
												Description: "Python programming language.",
												Emoji:       discordgo.ComponentEmoji{Name: ""},
												Default:     false,
											},
										},
									},
								},
							},
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.Button{
										Label:    "End the poll",
										CustomID: "end_poll",
										Style:    discordgo.DangerButton,
										Disabled: false,
									},
								},
							},
						},
					},
				})
				if err != nil {
					panic(err)
				}
				return
			}
			if _, ok := stats.PeopleVoted[i.Member.User.ID]; ok {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You already voted",
						Flags:   1 << 6,
					},
				})
				return
			}
			stats.PeopleVoted[i.Member.User.ID] = i.MessageComponentData().Values[0]
			if stats.Votes == nil {
				stats.Votes = map[string]int{"Go": 0, "Python": 0, "JS": 0}
			}

			stats.Votes[i.MessageComponentData().Values[0]]++
			pollVotes[i.Message.ID] = stats
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						constructVotesEmbed(pollVotes[i.Message.ID]),
					},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.SelectMenu{
									CustomID:    "choice",
									Placeholder: "Select your favorite language",
									Options: []discordgo.SelectMenuOption{
										{
											Label:       "Go",
											Value:       "Go",
											Description: "Go programming language.",
											Emoji:       discordgo.ComponentEmoji{},
											Default:     false,
										},
										{
											Label:       "JS",
											Value:       "JS",
											Description: "JavaScript programming language.",
											Emoji:       discordgo.ComponentEmoji{},
											Default:     false,
										},
										{
											Label:       "Python",
											Value:       "Python",
											Description: "Python programming language.",
											Emoji:       discordgo.ComponentEmoji{Name: ""},
											Default:     false,
										},
									},
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label:    "End the poll",
									CustomID: "end_poll",
									Style:    discordgo.DangerButton,
									Disabled: false,
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
	})
	_, err := s.ApplicationCommandCreate(*AppID, *GuildID, &discordgo.ApplicationCommand{
		Name:        "buttons",
		Description: "Test the buttons if you got courage",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}
	_, err = s.ApplicationCommandCreate(*AppID, *GuildID, &discordgo.ApplicationCommand{
		Name:        "poll",
		Description: "Lo and behold: dropdowns are coming",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}
