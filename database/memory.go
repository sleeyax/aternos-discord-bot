package database

import "github.com/sleeyax/aternos-discord-bot/database/models"

// MemoryDatabase is a simple in-memory database.
// Useful if you just want to self-host a bot for one discord server.
type MemoryDatabase struct {
	// Aternos discord cookie (ATERNOS_SESSION).
	SessionCookie string

	// Aternos server cookie (ATERNOS_SERVER).
	ServerCookie string
}

func NewInMemory(session string, server string) *MemoryDatabase {
	return &MemoryDatabase{session, server}
}

func (m *MemoryDatabase) Connect() error {
	return nil
}

func (m *MemoryDatabase) Disconnect() error {
	return nil
}

func (m *MemoryDatabase) ReadServerSettings(guildId string) (models.ServerSettings, error) {
	return models.ServerSettings{
		GuildID:       guildId,
		SessionCookie: m.SessionCookie,
		ServerCookie:  m.ServerCookie,
	}, nil
}

func (m *MemoryDatabase) UpdateServerSettings(settings *models.ServerSettings) error {
	m.SessionCookie = settings.SessionCookie
	m.ServerCookie = settings.ServerCookie
	return nil
}

func (m *MemoryDatabase) DeleteServerSettings(guildId string) error {
	m.SessionCookie = ""
	m.ServerCookie = ""
	return nil
}
