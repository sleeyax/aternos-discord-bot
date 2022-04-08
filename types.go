package aternos_discord_bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sleeyax/aternos-discord-bot/database"
	"github.com/sleeyax/aternos-discord-bot/worker"
)

type Bot struct {
	// Database instance that implements the database.Database interface.
	Database database.Database

	// Discord bot token.
	DiscordToken string

	// Current discord bot session.
	discord *discordgo.Session

	// Map of active workers for each discord server.
	workers worker.WorkersMap

	// List of registered discord commands.
	// These can be used to delete them once the bot has been stopped or removed from the discord server.
	registeredCommands []*discordgo.ApplicationCommand
}
