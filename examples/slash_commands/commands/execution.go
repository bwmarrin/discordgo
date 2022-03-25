package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

func DefineCommandExecution() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command":            ExecBasicCommand,
		"basic-command-with-files": ExecBasicCommandWithFiles,
		"options":                  ExecOption,
		"subcommands":              ExecSubcommands,
		"responses":                ExecResponses,
		"followups":                ExecFollowUps,
	}
}

func ExecBasicCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey there! Congratulations, you just executed your first slash command",
		},
	})
}

func ExecBasicCommandWithFiles(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

func ExecOption(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Here we need to convert raw interface{} value to wanted type.
	// Also, as you can see, here is used utility functions to convert the value
	// to particular type. Yeah, you can use just switch type,
	// but this is much simpler
	margs := []interface{}{
		i.ApplicationCommandData().Options[0].StringValue(),
		i.ApplicationCommandData().Options[1].IntValue(),
		i.ApplicationCommandData().Options[2].FloatValue(),
		i.ApplicationCommandData().Options[3].BoolValue(),
	}

	msgformat :=
		` Now you just learned how to use command options. Take a look to the value of which you've just entered:
				> string_option: %s
				> integer_option: %d
				> number_option: %f
				> bool_option: %v
		`
	if len(i.ApplicationCommandData().Options) >= 5 {
		margs = append(margs, i.ApplicationCommandData().Options[4].ChannelValue(nil).ID)
		msgformat += "> channel-option: <#%s>\n"
	}
	if len(i.ApplicationCommandData().Options) >= 6 {
		margs = append(margs, i.ApplicationCommandData().Options[5].UserValue(nil).ID)
		msgformat += "> user-option: <@%s>\n"
	}
	if len(i.ApplicationCommandData().Options) >= 7 {
		margs = append(margs, i.ApplicationCommandData().Options[6].RoleValue(nil, "").ID)
		msgformat += "> role-option: <@&%s>\n"
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, we'll discuss them in "responses" part
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
}

func ExecSubcommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Here we need to convert raw interface{} value to wanted type.
	// Also, as you can see, here is used utility functions to convert the value
	// to particular type. Yeah, you can use just switch type,
	// but this is much simpler
	margs := []interface{}{
		i.ApplicationCommandData().Options[0].StringValue(),
		i.ApplicationCommandData().Options[1].IntValue(),
		i.ApplicationCommandData().Options[2].FloatValue(),
		i.ApplicationCommandData().Options[3].BoolValue(),
	}

	msgformat :=
		` Now you just learned how to use command options. Take a look to the value of which you've just entered:
				> string_option: %s
				> integer_option: %d
				> number_option: %f
				> bool_option: %v
		`

	if len(i.ApplicationCommandData().Options) >= 5 {
		margs = append(margs, i.ApplicationCommandData().Options[4].ChannelValue(nil).ID)
		msgformat += "> channel-option: <#%s>\n"
	}
	if len(i.ApplicationCommandData().Options) >= 6 {
		margs = append(margs, i.ApplicationCommandData().Options[5].UserValue(nil).ID)
		msgformat += "> user-option: <@%s>\n"
	}
	if len(i.ApplicationCommandData().Options) >= 7 {
		margs = append(margs, i.ApplicationCommandData().Options[6].RoleValue(nil, "").ID)
		msgformat += "> role-option: <@&%s>\n"
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, we'll discuss them in "responses" part
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
}

func ExecResponses(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Responses to a command are very important.
	// First of all, because you need to react to the interaction
	// by sending the response in 3 seconds after receiving, otherwise
	// interaction will be considered invalid and you can no longer
	// use the interaction token and ID for responding to the user's request

	content := ""
	// As you can see, the response type names used here are pretty self-explanatory,
	// but for those who want more information see the official documentation
	switch i.ApplicationCommandData().Options[0].IntValue() {
	case int64(discordgo.InteractionResponseChannelMessageWithSource):
		content =
			"You just responded to an interaction, sent a message and showed the original one. " +
				"Congratulations!"
		content +=
			"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
	default:
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
		})
		if err != nil {
			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went wrong",
			})
		}
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
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
				Content: "Something went wrong",
			})
			return
		}
		time.Sleep(time.Second * 10)
		s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
	})
}

func ExecFollowUps(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Followup messages are basically regular messages (you can create as many of them as you wish)
	// but work as they are created by webhooks and their functionality
	// is for handling additional messages after sending a response.

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Note: this isn't documented, but you can use that if you want to.
			// This flag just allows you to create messages visible only for the caller of the command
			// (user who triggered the command)
			Flags:   1 << 6,
			Content: "Surprise!",
		},
	})
	msg, err := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
		Content: "Followup message has been created, after 5 seconds it will be edited",
	})
	if err != nil {
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
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
