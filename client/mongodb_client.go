package client

import (
	"context"
	"errors"
	"github.com/jinwoo1225/random-dice/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	ErrFailedToConvertToObjectId = errors.New("failed to convert to object id")
)

type MongoDBClient struct {
	client *mongo.Client
}

func NewMongoDBClient(ctx context.Context, conf *config.Config) (*MongoDBClient, func(), error) {
	opts := options.Client()
	opts.ApplyURI(conf.MongoDB.Host)
	opts.SetAuth(options.Credential{
		AuthMechanism:           "",
		AuthMechanismProperties: nil,
		AuthSource:              "",
		Username:                conf.MongoDB.Username,
		Password:                conf.MongoDB.Password,
		PasswordSet:             true,
	})

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, func() {}, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, func() {}, err
	}

	cleanUpFunc := func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Println(err)
		}
	}

	return &MongoDBClient{client: client}, cleanUpFunc, nil
}

func (c *MongoDBClient) InsertOne(ctx context.Context, database string, collection string, data interface{}) (*primitive.ObjectID, error) {
	res, err := c.client.Database(database).Collection(collection).InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	v, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, ErrFailedToConvertToObjectId
	}

	return &v, nil
}
