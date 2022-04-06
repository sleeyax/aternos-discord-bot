package aternos_discord_bot

import "github.com/bwmarrin/discordgo"

const (
	PingCommand      = "ping"
	ConfigureCommand = "configure"
	StartCommand     = "start"
	StopCommand      = "stop"
	StatusCommand    = "status"
	InfoCommand      = "status"
	PlayersCommand   = "players"
	CookiesOption    = "cookies"
)

// List of available discord commands.
var commands = []*discordgo.ApplicationCommand{
	{
		Name:        ConfigureCommand,
		Description: "Save configuration settings",
		Options: []*discordgo.ApplicationCommandOption{{
			Name:         CookiesOption,
			Description:  "Set the authentication cookies",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
		}},
	},
	{
		Name:        StartCommand,
		Description: "Start the minecraft server",
	},
	{
		Name:        StopCommand,
		Description: "Stop the minecraft server",
	},
	{
		Name:        PingCommand,
		Description: "Check if the discord bot is still alive",
	},
	{
		Name:        StatusCommand,
		Description: "Get the minecraft server status",
	},
	{
		Name:        InfoCommand,
		Description: "Get detailed information about the minecraft server status",
	},
	{
		Name:        PlayersCommand,
		Description: "List active players",
	},
}
