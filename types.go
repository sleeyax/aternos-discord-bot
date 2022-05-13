package aternos_discord_bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sleeyax/aternos-discord-bot/database"
	"github.com/sleeyax/aternos-discord-bot/worker"
	"net/url"
)

type Bot struct {
	// Database instance that implements the database.Database interface.
	Database database.Database

	// Discord bot token.
	DiscordToken string

	// HTTP(S) proxy to connect to.
	//
	// Format: protocol://username:password@host:port
	//
	// Examples:
	//
	// http://127.0.0.1:8888
	//
	// http://username:password@1.2.3.4:1234
	Proxy *url.URL

	// Current discord bot session.
	discord *discordgo.Session

	// Map of active workers for each discord server.
	workers worker.WorkersMap

	// List of registered discord commands.
	// These can be used to delete them once the bot has been stopped.
	registeredCommands []*discordgo.ApplicationCommand
}
