package database

import "github.com/sleeyax/aternos-discord-bot/database/models"

type Database interface {
	Connect() error
	Disconnect() error
	GetServerSettings() ([]models.ServerSettings, error)
	SaveServerSettings(settings *models.ServerSettings) error
}
