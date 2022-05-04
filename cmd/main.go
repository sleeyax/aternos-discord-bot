package main

import (
	"fmt"
	discordBot "github.com/sleeyax/aternos-discord-bot"
	"github.com/sleeyax/aternos-discord-bot/database"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Read configuration settings from environment variables
	token := os.Getenv("DISCORD_TOKEN")
	session := os.Getenv("ATERNOS_SESSION")
	server := os.Getenv("ATERNOS_SERVER")
	mongoDbUri := os.Getenv("MONGO_DB_URI")
	proxy := os.Getenv("PROXY")

	// Validate values
	if token == "" || (mongoDbUri == "" && (session == "" || server == "")) {
		log.Fatalln("Missing environment variables!")
	}

	bot := discordBot.Bot{
		DiscordToken: token,
	}

	if mongoDbUri != "" {
		bot.Database = database.NewMongo(mongoDbUri)
	} else {
		bot.Database = database.NewInMemory(session, server)
	}

	if proxy != "" {
		u, err := url.Parse(proxy)
		if err != nil {
			log.Fatalln(err)
		}
		bot.Proxy = u
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
