package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"yuzaki/config"
	"yuzaki/handlers"

	"github.com/bwmarrin/discordgo"
)

var (
	session *discordgo.Session
)

func init() {
	var err error
	if err = config.Load(); err != nil {
		log.Fatalf("loading config: %v", err)
	}

	session, err = discordgo.New("Bot " + config.BotConfig.DiscordToken)
	if err != nil {
		log.Fatalf("creating discord session: %v", err)
	}

	session.Identify.Intents = discordgo.IntentsAll
	session.AddHandler(ready)
	session.AddHandler(handlers.MessageGatewayHandler)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := setupAndRun(ctx); err != nil {
		log.Fatalf("error running bot: %v", err)
	}
}

func setupAndRun(ctx context.Context) error {
	if err := session.Open(); err != nil {
		return err
	}
	defer session.Close()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-shutdown:
		log.Println("Shutting down gracefully...")
		return nil
	}
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s on %d guilds", event.User.String(), len(event.Guilds))
	if err := s.UpdateWatchStatus(0, "all users in Yuzaki's Canyon!"); err != nil {
		log.Printf("error setting status: %v", err)
	}
	log.Println("Bot is now running. Press CTRL-C to exit.")

	// On guild 1009009522767052860, assign all users a role with role id 1307471288415162428
	// var rolesAssigned, rolesSkipped int
	// guild, err := s.Guild("1009009522767052860")
	// if err != nil {
	// 	log.Printf("error getting guild: %v", err)
	// 	return
	// }

	// var lastMemberID string
	// for {
	// 	members, err := s.GuildMembers(guild.ID, lastMemberID, 1000)
	// 	if err != nil {
	// 		log.Printf("error getting guild members: %v", err)
	// 		return
	// 	}
	// 	if len(members) == 0 {
	// 		break
	// 	}

	// 	for _, member := range members {
	// 		// skip if bot
	// 		if member.User.Bot {
	// 			continue
	// 		}

	// 		// Check if member already has the role
	// 		hasRole := false
	// 		for _, roleID := range member.Roles {
	// 			if roleID == "1307471288415162428" {
	// 				hasRole = true
	// 				break
	// 			}
	// 		}

	// 		if hasRole {
	// 			log.Printf("Skipping role assignment for %s (%s) - already has role", member.User.Username, member.User.ID)
	// 			rolesSkipped++
	// 			continue
	// 		}

	// 		if err := s.GuildMemberRoleAdd(guild.ID, member.User.ID, "1307471288415162428"); err != nil {
	// 			log.Printf("error adding role to %s (%s): %v", member.User.Username, member.User.ID, err)
	// 			continue
	// 		}
	// 		log.Printf("Assigned role to %s (%s)", member.User.Username, member.User.ID)
	// 		rolesAssigned++
	// 	}

	// 	if len(members) < 1000 {
	// 		break
	// 	}
	// 	lastMemberID = members[len(members)-1].User.ID
	// }

	// log.Printf("Operation complete: Assigned roles to %d members, skipped %d members", rolesAssigned, rolesSkipped)
}
