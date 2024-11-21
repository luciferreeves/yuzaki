package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"yuzaki/commands"
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
	session.AddHandler(handlers.InteractionCreateHandler)
	session.AddHandler(handlers.MemberAdd)
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

	log.Printf("Adding commands to %d guilds", len(session.State.Guilds))
	addApplicationCommands(session)

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

func addApplicationCommands(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		registeredCommands, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, guild.ID, commands.Commands)
		if err != nil {
			log.Printf("error adding commands to guild %s: %v", guild.ID, err)
			continue
		}
		log.Printf("added %d commands to guild %s", len(registeredCommands), guild.ID)
		for _, command := range registeredCommands {
			log.Printf("registered command %s", command.Name)
		}
	}
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s on %d guilds", event.User.String(), len(event.Guilds))
	if err := s.UpdateWatchStatus(0, "all users in Yuzaki's Canyon!"); err != nil {
		log.Printf("error setting status: %v", err)
	}
	log.Println("Bot is now running. Press CTRL-C to exit.")
}
