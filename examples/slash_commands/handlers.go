package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Login(ds *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", ds.State.User.Username, ds.State.User.Discriminator)
}

func HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmdName := i.ApplicationCommandData().Name

	h, ok := commandHandlers[cmdName]
	if !ok {
		log.Panicf("No handler found for command: %s", cmdName)
	}

	defer removeCommandsOnFailure(s, i.GuildID)
	// This way commands will be deleted even if panic occurs inside one of the handlers
	h(s, i)
}

func removeCommandsOnFailure(ds *discordgo.Session, guildID string) {
	r := recover()
	if r != nil {
		deleteCommands(ds, guildID)
		log.Fatal("Error occured inside command handler")
	}
}

func Basic(s *discordgo.Session, i *discordgo.InteractionCreate) {
	panic("test")
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey there! Congratulations, you've just executed your first slash command!",
		},
	})
}

func BasicWithFile(s *discordgo.Session, i *discordgo.InteractionCreate) {

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey there! Congratulations, you just executed your first slash command with a file in the response",
			Files: []*discordgo.File{
				{
					ContentType: "text/plain",
					Name:        "test.txt",
					Reader:      strings.NewReader("Hello Discord!!"),
				},
			},
		},
	})
}

func Localized(s *discordgo.Session, i *discordgo.InteractionCreate) {
	responses := map[discordgo.Locale]string{
		discordgo.ChineseCN: "你好！ 这是一个本地化的命令",
	}
	response := "Hi! This is a localized message"
	if r, ok := responses[i.Locale]; ok {
		response = r
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		panic(err)
	}
}

func Options(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// This example stores the provided arguments in an []interface{}
	// which will be used to format the bot's response
	margs := make([]interface{}, 0, len(options))
	msgformat := "You learned how to use command options! " +
		"Take a look at the value(s) you entered:\n"

	// Get the value from the option map.
	// When the option exists, ok = true
	if option, ok := optionMap["string-option"]; ok {
		// Option values must be type asserted from interface{}.
		// Discordgo provides utility functions to make this simple.
		margs = append(margs, option.StringValue())
		msgformat += "> string-option: %s\n"
	}

	if opt, ok := optionMap["integer-option"]; ok {
		margs = append(margs, opt.IntValue())
		msgformat += "> integer-option: %d\n"
	}

	if opt, ok := optionMap["number-option"]; ok {
		margs = append(margs, opt.FloatValue())
		msgformat += "> number-option: %f\n"
	}

	if opt, ok := optionMap["bool-option"]; ok {
		margs = append(margs, opt.BoolValue())
		msgformat += "> bool-option: %v\n"
	}

	if opt, ok := optionMap["channel-option"]; ok {
		margs = append(margs, opt.ChannelValue(nil).ID)
		msgformat += "> channel-option: <#%s>\n"
	}

	if opt, ok := optionMap["user-option"]; ok {
		margs = append(margs, opt.UserValue(nil).ID)
		msgformat += "> user-option: <@%s>\n"
	}

	if opt, ok := optionMap["role-option"]; ok {
		margs = append(margs, opt.RoleValue(nil, "").ID)
		msgformat += "> role-option: <@&%s>\n"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
}

func Subcommands(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options
	content := ""

	// As you can see, names of subcommands (nested, top-level)
	// and subcommand groups are provided through the arguments.
	switch options[0].Name {
	case "subcommand":
		content = "The top-level subcommand is executed. Now try to execute the nested one."
	case "subcommand-group":
		options = options[0].Options
		switch options[0].Name {
		case "nested-subcommand":
			content = "Nice, now you know how to execute nested commands too"
		default:
			content = "Oops, something went wrong.\n" +
				"Hol' up, you aren't supposed to see this message."
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func Responses(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Responses to command are very important.
	// First of all, because you need to react to the interaction
	// by sending the response in 3 seconds after receiving, otherwise
	// interaction will be considered invalid and you can no longer
	// use the interaction token and ID for responding to the user's request

	// As you can see, the response type names used here are pretty self-explanatory,
	// but for those who want more information see the official documentation

	content := ""

	switch i.ApplicationCommandData().Options[0].IntValue() {
	case int64(discordgo.InteractionResponseChannelMessageWithSource):
		content =
			"You've just responded to user inte interaction. User input is supposed to be shown as well but non were given."
	case int64(discordgo.InteractionResponseDeferredChannelMessageWithSource):
		content =
			"You've just deferedly responded to user interaction. Congratulations!"
	default:
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
		})

		if err != nil {
			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went wrong.",
			})
		}
		return
	}
	content +=
		"\nAlso... you can edit response after it was sent. Wait for 5 seconds and this message will be changed."

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})

	if err != nil {
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong.",
		})
		return
	}

	time.AfterFunc(time.Second*5, func() {
		_, err = s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
			Content: content + "\n\nWell, now you know how to create and edit responses. " +
				"But you still don't know how to delete them... so... wait 10 seconds and this " +
				"message will be deleted.",
		})

		if err != nil {
			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went wrong.",
			})
			return
		}

		time.Sleep(time.Second * 10)
		s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
	})
}

func Followups(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Followup messages are basically regular messages (you can create as many of them as you wish)
	// but work as they are created by webhooks and their functionality
	// is for handling additional messages after sending a response.

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Note: this isn't documented, but you can use that if you want to.
			// This flag allows you to create messages visible only for the caller of the command
			// (user who triggered the command)
			Flags:   1 << 6,
			Content: "Surprise!",
		},
	})

	msg, err := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
		Content: "Followup message has been created, after 5 seconds it will be edited.",
	})

	if err != nil {
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong.",
		})
		return
	}
	time.Sleep(time.Second * 5)

	s.FollowupMessageEdit(s.State.User.ID, i.Interaction, msg.ID, &discordgo.WebhookEdit{
		Content: "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted.",
	})

	time.Sleep(time.Second * 10)

	s.FollowupMessageDelete(s.State.User.ID, i.Interaction, msg.ID)

	s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
		Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
			"take a unicorn :unicorn: as reward!\n" +
			"Also, as bonus... look at the original interaction response :D",
	})
}
