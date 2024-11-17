package handlers

import (
	messagehandlers "yuzaki/handlers/messageHandlers"

	"github.com/bwmarrin/discordgo"
)

func MessageGatewayHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	messagehandlers.PoketwoHandler(s, m)

	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

}
