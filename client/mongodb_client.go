package client

import (
	"context"
	"github.com/jinwoo1225/random-dice/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
