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
	configuration *config.Config
	session       *discordgo.Session
)

func init() {
	var err error
	configuration, err = config.Load()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	session, err = discordgo.New("Bot " + configuration.DiscordToken)
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
}
