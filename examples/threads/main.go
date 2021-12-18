package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	BotToken = flag.String("token", "", "Bot token")
)

var games map[string]time.Time = make(map[string]time.Time)

func main() {
	flag.Parse()
	s, _ := discordgo.New("Bot " + *BotToken)
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Ready")
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if strings.Contains(m.Content, "ping") {
			if ch, err := s.State.Channel(m.ChannelID); err != nil || !ch.IsThread() {
				thread, err := s.StartMessageThreadComplex(m.ChannelID, m.ID, &discordgo.ThreadStart{
					Name:                "Pong game with " + m.Author.Username,
					AutoArchiveDuration: 60,
					Invitable:           false,
					RateLimitPerUser:    10,
				})
				if err != nil {
					panic(err)
				}
				_, _ = s.ChannelMessageSend(thread.ID, "Pong")
				m.ChannelID = thread.ID
			} else {
				_, _ = s.ChannelMessageSendReply(m.ChannelID, "Pong", m.Reference())
			}
			games[m.ChannelID] = time.Now()
			<-time.After(time.Second * 10)
			if time.Since(games[m.ChannelID]) >= time.Second*10 {
				s.ChannelEditComplex(m.ChannelID, &discordgo.ChannelEdit{
					Archived: true,
					Locked:   true,
				})
			}
		}
	})
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

}
