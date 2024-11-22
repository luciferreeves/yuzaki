package admin

import (
	"fmt"
	"log"
	"time"

	"yuzaki/utils"

	"github.com/bwmarrin/discordgo"
)


func PurgeChat(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	option, ok := optionMap["amount"]
	if !ok {
		utils.SendFollowUpMessage(s, i, "You must provide an amount to purge", true)
		return
	}

	amount := option.IntValue()
	if amount < 1 {
		utils.SendFollowUpMessage(s, i, "Amount must be greater than 0", true)
		return
	}

	channel := i.ChannelID
	remaining := amount
	deletedCount := int64(0)

	for remaining > 0 {
		batchSize := int64(100)
		if remaining < 100 {
			batchSize = remaining
		}

		messages, err := s.ChannelMessages(channel, int(batchSize), "", "", "")
		if err != nil {
			utils.SendFollowUpMessage(s, i, "Failed to fetch messages", true)
			return
		}

		if len(messages) == 0 {
			break
		}

		messageIDs := make([]string, len(messages))
		for i, message := range messages {
			messageIDs[i] = message.ID
		}

		if len(messageIDs) > 1 {
			err = s.ChannelMessagesBulkDelete(channel, messageIDs)
			if err != nil {
				for _, msgID := range messageIDs {
					err = s.ChannelMessageDelete(channel, msgID)
					if err != nil {
						continue
					}
					deletedCount++
					time.Sleep(time.Millisecond * 100)
				}
			} else {
				deletedCount += int64(len(messageIDs))
			}
		} else if len(messageIDs) == 1 {
			err = s.ChannelMessageDelete(channel, messageIDs[0])
			if err == nil {
				deletedCount++
			}
			time.Sleep(time.Millisecond * 100)
		}

		remaining -= batchSize
	}

	utils.SendFollowUpMessage(s, i, fmt.Sprintf("Successfully purged %d messages!", deletedCount), false)
	_, err = s.ChannelMessageSend(channel, fmt.Sprintf("<@%s> has purged %d messages from this channel", i.Member.User.ID, deletedCount))
	if err != nil {
		log.Printf("Failed to send message after purging messages. Successfully purged %d messages from channel %s. Purged By: %s\n", deletedCount, channel, i.Member.DisplayName())
	} else {
		log.Printf("Successfully purged %d messages from channel %s. Purged By: %s\n", deletedCount, channel, i.Member.DisplayName())
	}
}
