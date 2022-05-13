package database

import (
	"errors"
	"github.com/sleeyax/aternos-discord-bot/database/models"
)

// ErrDataNotFound is a reusable error that can be used by different database implementations to denote that a row or document was not found.
var ErrDataNotFound = errors.New("database: row or document data not found")

type Database interface {
	Connect() error
	Disconnect() error
	ReadServerSettings(guildId string) (models.ServerSettings, error)
	UpdateServerSettings(settings *models.ServerSettings) error
	DeleteServerSettings(guildId string) error
}
