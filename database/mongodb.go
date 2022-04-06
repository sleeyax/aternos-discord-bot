package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDb struct {
	// MongoDB's connection URI string.
	ConnectionURI string

	// Connection timeout.
	//
	// Defaults to 10 seconds.
	ConnectionTimeout time.Duration

	client *mongo.Client
}

func New(uri string) *MongoDb {
	return &MongoDb{ConnectionURI: uri, ConnectionTimeout: time.Second * 10}
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

func (db *MongoDb) ListDatabases() ([]string, error) {
	databases, err := db.client.ListDatabaseNames(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	return databases, nil
}
