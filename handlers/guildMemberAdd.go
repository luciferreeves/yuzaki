package handlers

import (
	"log"
	"github.com/bwmarrin/discordgo"
)

func MemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if m.GuildID == "1009009522767052860" { // Yuzaki Guild
		err := s.GuildMemberRoleAdd(m.GuildID, m.User.ID, "1307471288415162428") // Members Role
		if err != nil {
			log.Println(err)
		}
	}


}
