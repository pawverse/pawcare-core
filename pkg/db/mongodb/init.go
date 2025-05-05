package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbCloseFunc func(context.Context) error

func clientOptions(connectionString string) *options.ClientOptions {
	clientOptions := options.Client().ApplyURI(connectionString)

	return clientOptions
}

func createClient(ctx context.Context, clientOptions *options.ClientOptions) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ConnectDB(ctx context.Context, connectionString string, databaseName string) (*mongo.Database, DbCloseFunc, error) {
	clientOptions := clientOptions(connectionString)

	client, err := createClient(ctx, clientOptions)
	if err != nil {
		return nil, func(ctx context.Context) error { return nil }, err
	}

	disconnect := func(ctx context.Context) error {
		return client.Disconnect(ctx)
	}

	return client.Database(databaseName), disconnect, nil
}
