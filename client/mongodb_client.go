package client

import (
	"context"
	"errors"
	"log"

	"github.com/jinwoo1225/random-dice/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrFailedToConvertToObjectID = errors.New("failed to convert to object id")
)

type MongoDBClient interface {
	InsertOne(ctx context.Context, database string, collection string, data interface{}) (*primitive.ObjectID, error)
	FindOne(ctx context.Context, database string, collection string, filter interface{}) (*mongo.SingleResult, error)
	FindMany(ctx context.Context, database string, collection string, filter interface{}, orderBy interface{}, page int64, limit int64) (*mongo.Cursor, error)
}

type DefaultMongoDBClient struct {
	client *mongo.Client
}

func NewMongoDBClient(ctx context.Context, conf *config.Config) (*DefaultMongoDBClient, func(), error) {
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

	return &DefaultMongoDBClient{client: client}, cleanUpFunc, nil
}

func (c *DefaultMongoDBClient) InsertOne(ctx context.Context, database string, collection string, data interface{}) (*primitive.ObjectID, error) {
	res, err := c.client.Database(database).Collection(collection).InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	v, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, ErrFailedToConvertToObjectID
	}

	return &v, nil
}

func (c *DefaultMongoDBClient) FindOne(ctx context.Context, database string, collection string, filter interface{}) (*mongo.SingleResult, error) {
	return c.client.Database(database).Collection(collection).FindOne(ctx, filter, nil), nil
}

func (c *DefaultMongoDBClient) FindMany(
	ctx context.Context,
	database string,
	collection string,
	filter interface{},
	orderBy interface{},
	page int64,
	limit int64,
) (*mongo.Cursor, error) {
	opts := options.Find()
	opts.Limit = &limit
	opts.Skip = &page

	opts.SetSort(orderBy)

	res, err := c.client.Database(database).Collection(collection).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	return res, nil
}
