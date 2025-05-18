package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Mongo struct {
	client *mongo.Client
	// collection represents a MongoDB collection used for performing database operations.
	collection *mongo.Collection
}

func New(uri string, db string, collection string) (*Mongo, error) {
	const op = "storage.mongo.New"

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Mongo{
		client:     client,
		collection: client.Database(db).Collection(collection),
	}, nil
}

func (m *Mongo) Disconnect(ctx context.Context) error {
	const op = "storage.mongo.Disconnect"
	if err := m.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
