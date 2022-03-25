package commands

import "github.com/bwmarrin/discordgo"

var integerOptionMinValue = 1.0

/* ---------------------------------------------------------------------------------------------------------------------

												NEED 2 KNOW
											-------------------

	1. All commands and options must have a description, Commands/options without description will fail the registration
	2. For ApplicationCommandOption: required options must be listed before optional parameters. The same concept applies
		to Discord's Slash-commands API
	3. For Subcommands: When a command has subcommands/subcommand groups, it must not have top-level options, they aren't
		accesible in the UI. In this case (at least not yet), so if a command has subcommands/subcommand any groups
		registering top-level options will cause the registration of the command to fail
	4. Also, you can create both subcommand groups and subcommands in the command at the same time. But, there's some
		limits to nesting, count of subcommands (top level and nested) and options. Read the intro of slash-commands
		docs on Discord dev portal to get more information


	ApplicationCommandOption Types:
		(ApplicationCommandOptionString),
		(ApplicationCommandOptionInteger: MinValue, MaxValue, Choices),
		(ApplicationCommandOptionNumber: MinValue, MaxValue),
		(ApplicationCommandOptionBoolean),
		(ApplicationCommandOptionChannel: ChannelTypes),
		(ApplicationCommandOptionUser),
		(ApplicationCommandOptionRole)

--------------------------------------------------------------------------------------------------------------------- */

func DefineCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "basic-command",
			Description: "Basic command",
		},
		// -------------------------------------------------------------------------------------------------------------
		{
			Name:        "basic-command-with-files",
			Description: "Basic command with files",
		},
		// -------------------------------------------------------------------------------------------------------------
		{
			Name:        "options",
			Description: "Command for demonstrating options",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "string-option",
					Description: "String option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "integer-option",
					Description: "Integer option",
					MinValue:    &integerOptionMinValue,
					MaxValue:    10,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel-option",
					Description: "Channel option",
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
						discordgo.ChannelTypeGuildVoice,
					},
					Required: false,
				},
			},
		},
		// -------------------------------------------------------------------------------------------------------------
		{
			Name:        "subcommands",
			Description: "Subcommands and command groups example",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "scmd-grp",
					Description: "Subcommands group",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "nst-subcmd",
							Description: "Nested subcommand",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommandGroup,
				},
				{
					Name:        "subcmd",
					Description: "Top-level subcommand",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		// -------------------------------------------------------------------------------------------------------------
		{
			Name:        "responses",
			Description: "Interaction responses testing initiative",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "resp-type",
					Description: "Response type",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Channel message with source",
							Value: 4,
						},
						{
							Name:  "Deferred response With Source",
							Value: 5,
						},
					},
					Required: true,
				},
			},
		},
		// -------------------------------------------------------------------------------------------------------------
		{
			Name:        "followups",
			Description: "Followup messages",
		},
	}
}
