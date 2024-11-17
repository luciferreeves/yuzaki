package commands

import (
	"github.com/bwmarrin/discordgo"
)

var (
	manageMessagesPermissions int64 = discordgo.PermissionManageMessages
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "purge",
			Description:              "Purge messages from the current channel",
			DefaultMemberPermissions: &manageMessagesPermissions,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "amount",
					Description: "The amount of messages to purge",
					Required:    true,
				},
			},
		},
	}
)
