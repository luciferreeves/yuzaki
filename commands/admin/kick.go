package admin

import (
	//"fmt"
	"log"
	//"time"

	"yuzaki/utils"

	"github.com/bwmarrin/discordgo"
)

func KickMember(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %s. Interaction: purge. Interaction By: %s\n", err, i.Member.DisplayName())
		return
	}

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}

	option, ok := optionMap["target"]
	if !ok {
		utils.SendFollowUpMessage(s, i, "You must mention the user to kick!", true)
		return
	}

	user := option.UserValue(s);
	if user == nil {
		utils.SendFollowUpMessage(s, i, "Invalid user", true)
		return
	}
		err = s.GuildMemberDeleteWithReason(i.GuildID, user.ID, "Kicked by an admin")
	if err != nil {
		utils.SendFollowUpMessage(s, i, "Could not kick user", true)
		return
	}

	utils.SendFollowUpMessage(s, i, "User kicked successfully!", true)
	
}