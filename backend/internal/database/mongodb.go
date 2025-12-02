package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
)

func ConnectMongoDB(uri, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database
	if err = Client.Ping(ctx, nil); err != nil {
		return err
	}

	Database = Client.Database(dbName)
	log.Println("Connected to MongoDB successfully!")
	return nil
}

func Close() error {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return Client.Disconnect(ctx)
	}
	return nil
}

// Helper function to get a collection
func GetCollection(name string) *mongo.Collection {
	if Database == nil {
		log.Fatal("Database connection is not initialized")
	}
	return Database.Collection(name)
}
