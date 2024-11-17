package messagehandlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	poketwoID       = "716390085896962058"
	poketwoChannels = []string{"1307335103508254844", "1307335132046168074"}
)

func PoketwoHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	isAllowed := isAllowedChannel(m.ChannelID)

	if m.Author.ID == poketwoID {
		if !isAllowed {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}
		return
	}

	isCommand := strings.HasPrefix(strings.ToLower(m.Content), "p!") ||
		(len(m.Mentions) > 0 && m.Mentions[0].ID == poketwoID)

	if isCommand && !isAllowed {
		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			return
		}

		fmt.Printf("%s#%s sent a Poketwo command in #%s - deleting\n",
			m.Author.Username, m.Author.Discriminator, channel.Name)

		go func() {
			time.Sleep(1 * time.Second)
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}()

		var channelMentions string
		for i, channelID := range poketwoChannels {
			if i == len(poketwoChannels)-1 && i > 0 {
				channelMentions += "or "
			}
			channelMentions += fmt.Sprintf("<#%s>", channelID)
			if i < len(poketwoChannels)-2 {
				channelMentions += ", "
			}
		}

		_, err = s.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf("<@%s>, please use %s for Pokétwo commands. If you need roles, visit <id:customize> to get them.",
				m.Author.ID, channelMentions))

		if err != nil {
			return
		}

		go func() {
			time.Sleep(1 * time.Second)
			messages, err := s.ChannelMessages(m.ChannelID, 5, "", "", "")
			if err != nil {
				return
			}

			for _, msg := range messages {
				if msg.Author.ID == poketwoID {
					s.ChannelMessageDelete(m.ChannelID, msg.ID)
					break
				}
			}
		}()
	}
}

func isAllowedChannel(channelID string) bool {
	for _, channel := range poketwoChannels {
		if channelID == channel {
			return true
		}
	}
	return false
}
