package aternos_discord_bot

import (
	"github.com/bwmarrin/discordgo"
	aternos "github.com/sleeyax/aternos-api"
)

type Connection struct {
	// Discord guild id.
	guildId string

	// Current server status.
	// This is basically only used as a cache for the info command.
	serverInfo *aternos.ServerInfo

	// Current active websocket connection.
	wss *aternos.Websocket

	// Aternos API instance.
	api *aternos.Api
}

type Bot struct {
	// Optional MongoDB database to connect to.
	//
	// If this is set it means the bot is operating on behalf of multiple discord servers.
	// Otherwise, it's only configured for one; in that case SessionCookie and ServerCookie should be set.
	MongoDbURI string

	// Discord bot token.
	DiscordToken string

	// Aternos discord cookie (ATERNOS_SESSION).
	SessionCookie string

	// Aternos server cookie (ATERNOS_SERVER).
	ServerCookie string

	// Current discord session.
	discord *discordgo.Session

	// Aternos api instance.
	api *aternos.Api

	// List of active connections for each discord server.
	connections []*Connection

	// List of registered discord commands.
	// These can be used to delete them once the bot has been stopped or removed from the discord server.
	registeredCommands []*discordgo.ApplicationCommand
}
