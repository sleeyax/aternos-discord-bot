package database

import (
	"context"
	"github.com/sleeyax/aternos-discord-bot/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDb struct {
	// MongoDB's connection URI string.
	ConnectionURI string

	// Name of the database to use.
	// Defaults to `aternos-discord-bot`
	DatabaseName string

	// Connection timeout.
	//
	// Defaults to 10 seconds.
	ConnectionTimeout time.Duration

	client *mongo.Client
}

func NewMongo(uri string) *MongoDb {
	return &MongoDb{ConnectionURI: uri, ConnectionTimeout: time.Second * 10, DatabaseName: "aternos-discord-bot"}
}

func (db *MongoDb) Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(db.ConnectionURI))
	if err != nil {
		return err
	}

	db.client = client

	ctx, cancel := context.WithTimeout(context.Background(), db.ConnectionTimeout)
	defer cancel()

	return client.Connect(ctx)
}

func (db *MongoDb) Disconnect() error {
	return db.client.Disconnect(context.Background())
}

func (db *MongoDb) ReadServerSettings(guildId string) (models.ServerSettings, error) {
	collection := db.client.Database(db.DatabaseName).Collection(models.ServerSettingsTable)
	ctx := context.Background()

	var settings models.ServerSettings
	err := collection.FindOne(ctx, models.ServerSettings{GuildID: guildId}).Decode(&settings)

	if err == mongo.ErrNoDocuments {
		return models.ServerSettings{}, ErrDataNotFound
	}
	if err != nil {
		return models.ServerSettings{}, err
	}

	return settings, nil
}

func (db *MongoDb) UpdateServerSettings(settings *models.ServerSettings) error {
	if settings.UpdatedAt == 0 {
		settings.UpdatedAt = time.Now().UnixMilli()
	}

	collection := db.client.Database(db.DatabaseName).Collection(models.ServerSettingsTable)

	var updatedSettings models.ServerSettings

	ctx := context.Background()
	filter := models.ServerSettings{GuildID: settings.GuildID}
	update := bson.D{{"$set", settings}, {"$setOnInsert", models.ServerSettings{CreatedAt: time.Now().UnixMilli()}}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedSettings)

	// ErrNoDocuments means that the filter did not match any documents in the collection.
	// In our case this means the document was updated instead.
	if err != mongo.ErrNoDocuments {
		return err
	}

	return nil
}

func (db *MongoDb) DeleteServerSettings(guildId string) error {
	collection := db.client.Database(db.DatabaseName).Collection(models.ServerSettingsTable)
	_, err := collection.DeleteOne(context.Background(), models.ServerSettings{GuildID: guildId})
	return err
}
