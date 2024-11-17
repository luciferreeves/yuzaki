package messagehandlers

import (
	"fmt"
	"strings"
	"time"
	"yuzaki/config"

	"github.com/bwmarrin/discordgo"
)

const poketwoID = "716390085896962058" // Pokétwo's user ID; TODO: Move to config?

func PoketwoHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	isAllowed := isAllowedChannel(m.ChannelID)

	if m.Author.ID == poketwoID {
		if !isAllowed {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}
		return
	}

	if isPoketwoCommand(m) && !isAllowed {
		handleUnauthorizedCommand(s, m)
	}
}

func isPoketwoCommand(m *discordgo.MessageCreate) bool {
	return strings.HasPrefix(strings.ToLower(m.Content), "p!") ||
		(len(m.Mentions) > 0 && m.Mentions[0].ID == poketwoID)
}

func handleUnauthorizedCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return
	}

	fmt.Printf("%s#%s sent a Poketwo command in #%s - deleting\n",
		m.Author.Username, m.Author.Discriminator, channel.Name)

	sendRedirectMessage(s, m.Author.ID, m.ChannelID)

	go deleteMessageWithDelay(s, m.ChannelID, m.ID)
	go cleanupPoketwoResponse(s, m.ChannelID)
}

func sendRedirectMessage(s *discordgo.Session, authorID, channelID string) {
	channelMentions := formatChannelMentions()
	message := fmt.Sprintf("<@%s>, please use %s for Pokétwo commands. If you need roles, visit <id:customize> to get them.",
		authorID, channelMentions)

	s.ChannelMessageSend(channelID, message)
}

func formatChannelMentions() string {
	var mentions []string
	for _, channelID := range config.BotConfig.ConfiguredChannels.PoketwoSpawns {
		mentions = append(mentions, fmt.Sprintf("<#%s>", channelID))
	}

	if len(mentions) > 1 {
		lastMention := mentions[len(mentions)-1]
		mentions = mentions[:len(mentions)-1]
		return strings.Join(mentions, ", ") + " or " + lastMention
	}
	return mentions[0]
}

func cleanupPoketwoResponse(s *discordgo.Session, channelID string) {
	time.Sleep(1 * time.Second)
	messages, err := s.ChannelMessages(channelID, 5, "", "", "")
	if err != nil {
		return
	}

	for _, msg := range messages {
		if msg.Author.ID == poketwoID {
			s.ChannelMessageDelete(channelID, msg.ID)
			break
		}
	}
}

func deleteMessageWithDelay(s *discordgo.Session, channelID, messageID string) {
	time.Sleep(2 * time.Second)
	s.ChannelMessageDelete(channelID, messageID)
}

func isAllowedChannel(channelID string) bool {
	for _, channel := range config.BotConfig.ConfiguredChannels.PoketwoSpawns {
		if channelID == channel {
			return true
		}
	}
	return false
}
