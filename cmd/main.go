package main

import (
	"fmt"
	discordBot "github.com/sleeyax/aternos-discord-bot"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Read auth creds from system env vars.
	token := os.Getenv("DISCORD_TOKEN")
	session := os.Getenv("ATERNOS_SESSION")
	server := os.Getenv("ATERNOS_SERVER")

	if token == "" || session == "" || server == "" {
		log.Fatalln("Missing environment variables!")
	}

	// Create discord bot instance.
	bot := discordBot.Bot{
		DiscordToken:  token,
		SessionCookie: session,
		ServerCookie:  server,
	}
	if err := bot.Start(); err != nil {
		log.Fatalln(err)
	}
	defer bot.Stop()

	// Wait until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-interruptSignal
}
