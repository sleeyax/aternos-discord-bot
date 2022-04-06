package models

const ServerSettingsTable = "server_settings"

// ServerSettings models the stored settings for a specific discord server.
type ServerSettings struct {
	GuildID       string `bson:"guildID,omitempty"`
	SessionCookie string `bson:"sessionCookie,omitempty"`
	ServerCookie  string `bson:"serverCookie,omitempty"`
	CreatedAt     int64  `bson:"createdAt,omitempty"`
	UpdatedAt     int64  `bson:"updatedAt,omitempty"`
}
