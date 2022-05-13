package aternos_discord_bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	aternos "github.com/sleeyax/aternos-api"
	"github.com/sleeyax/aternos-discord-bot/worker"
	"log"
	"net/http"
	"time"
)

func (ab *Bot) setupHandlers() {
	ab.discord.AddHandler(ab.handleCommands)
	ab.discord.AddHandler(ab.handleJoinServer)
	ab.discord.AddHandler(ab.handleLeaveServer)
}

func (ab *Bot) handleJoinServer(s *discordgo.Session, e *discordgo.GuildCreate) {
	// The GuildCreate event also fires after the bot has been restarted, so we have to check whether we joined recently or not.
	if time.Now().Sub(e.JoinedAt).Minutes() <= 2 {
		log.Printf("Joined new server %s (ID: %s)\n", e.Name, e.ID)
	}
}

func (ab *Bot) handleLeaveServer(s *discordgo.Session, e *discordgo.GuildDelete) {
	log.Printf("Left server %s (ID: %s)\n", e.BeforeDelete.Name, e.BeforeDelete.ID)
	ab.Database.DeleteServerSettings(e.BeforeDelete.ID)
}

func (ab *Bot) Start() error {
	if err := ab.Database.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %e", err)
	}

	ab.workers = make(map[string]*worker.Worker)

	session, err := discordgo.New("Bot " + ab.DiscordToken)
	if err != nil {
		return err
	}

	ab.discord = session
	ab.setupHandlers()

	if err = ab.discord.Open(); err != nil {
		return err
	}

	return ab.registerCommands()
}

func (ab *Bot) Stop() error {
	if err := ab.removeCommands(); err != nil {
		return err
	}

	if err := ab.Database.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect database: %e", err)
	}

	return ab.discord.Close()
}

// registerCommands registers all available Discord commands.
func (ab *Bot) registerCommands() error {
	var err error

	ab.registeredCommands, err = ab.discord.ApplicationCommandBulkOverwrite(ab.discord.State.User.ID, "", commands)

	return err
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

// getWorker returns the active worker for the specified guildId from the pool or creates a new one if it doesn't exist yet.
func (ab *Bot) getWorker(guildId string) (*worker.Worker, error) {
	w, ok := ab.workers[guildId]

	opts, err := ab.createOptions(guildId)
	if err != nil {
		return nil, err
	}

	if !ok {
		w = worker.New(guildId, opts)
		ab.workers[guildId] = w
	} else {
		w.Reconfigure(opts)
	}

	return w, nil
}

// createOptions creates new aternos configuration options.
func (ab *Bot) createOptions(guildId string) (*aternos.Options, error) {
	options := &aternos.Options{
		Cookies: []*http.Cookie{
			{
				Name:  "ATERNOS_LANGUAGE",
				Value: "en",
			},
		},
		Proxy:              ab.Proxy,
		InsecureSkipVerify: true,
	}

	settings, err := ab.Database.ReadServerSettings(guildId)
	if err != nil {
		return nil, err
	}

	options.Cookies = append(options.Cookies,
		&http.Cookie{
			Name:  "ATERNOS_SESSION",
			Value: settings.SessionCookie,
		},
		&http.Cookie{
			Name:  "ATERNOS_SERVER",
			Value: settings.ServerCookie,
		},
	)

	return options, nil
}
