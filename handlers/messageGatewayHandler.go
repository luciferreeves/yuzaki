package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func MessageGatewayHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	// if says "hello @bot" <- case insensitive <- bot responds with "Hello @user!"
	if strings.EqualFold(strings.TrimSpace(m.Content), "hello "+s.State.User.Mention()) {
		s.ChannelMessageSendReply(m.ChannelID, "Hello "+m.Author.Mention()+"!", m.Reference())
	}
}
