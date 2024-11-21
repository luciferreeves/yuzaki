package commands

import (
	"github.com/bwmarrin/discordgo"
)

var (
	manageMessagesPermissions int64 = discordgo.PermissionManageMessages
	PermissionKickMembers     int64 = discordgo.PermissionKickMembers
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
		{
			Name:                     "kick",
			Description:              "Kick a member from the server",
			DefaultMemberPermissions: &PermissionKickMembers,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "target",
					Description: "The member to kick",
					Required:    true,
				},
			},
		},
	}
)
