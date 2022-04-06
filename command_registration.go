package aternos_discord_bot

import (
	"github.com/bwmarrin/discordgo"
)

// registerCommands registers all available Discord commands.
func (ab *Bot) registerCommands() error {
	ab.registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))

	for i, v := range commands {
		cmd, err := ab.discord.ApplicationCommandCreate(ab.discord.State.User.ID, "", v)
		if err != nil {
			return err
		}
		ab.registeredCommands[i] = cmd
	}

	return nil
}

// removeCommands removes all Discord commands that were previously registered using registerCommands.
func (ab *Bot) removeCommands() error {
	for _, v := range ab.registeredCommands {
		if err := ab.discord.ApplicationCommandDelete(ab.discord.State.User.ID, "", v.ID); err != nil {
			return err
		}
	}

	ab.registeredCommands = nil

	return nil
}
