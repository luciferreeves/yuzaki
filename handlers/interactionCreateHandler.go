package handlers

import (
	"fmt"
	"yuzaki/commands/admin"

	"github.com/bwmarrin/discordgo"
)

var (
	SlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"purge": admin.PurgeChat,
		"kick": admin.KickMember,
	}
)

func InteractionCreateHandler(s *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		if handler, ok := SlashCommandHandlers[interaction.ApplicationCommandData().Name]; ok {
			handler(s, interaction)
		}
	case discordgo.InteractionMessageComponent:
		// Detect what type of message component interaction it is.
		switch interaction.MessageComponentData().ComponentType {
		case discordgo.ButtonComponent:
			fmt.Println("Button interaction detected.")
		case discordgo.SelectMenuComponent:
			fmt.Println("Select menu interaction detected.")
		default:
			fmt.Println("Unknown message component interaction detected.")
		}
	}
}
