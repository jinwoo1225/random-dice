package handler

import (
	"context"
	"errors"
	"github.com/benbjohnson/clock"
	"github.com/jinwoo1225/random-dice/client"
	randomdicev1 "github.com/jinwoo1225/random-dice/gen/proto/go/randomdice/v1"
	"time"
)

var (
	errEmptyName     = errors.New("empty name")
	errEmptyPassword = errors.New("empty password")
	errEmptyEmail    = errors.New("empty email")
)

const (
	MongoDatabaseName       = "random-dice"
	MongoUserCollectionName = "users"
)

type User struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Password  string    `bson:"password"`
	Email     string    `bson:"email"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type CreateUserFunc func(ctx context.Context, req *randomdicev1.CreateUserRequest) (*randomdicev1.CreateUserResponse, error)

func CreateUser(clk clock.Clock, mdb *client.DefaultMongoDBClient) CreateUserFunc {
	return func(ctx context.Context, req *randomdicev1.CreateUserRequest) (*randomdicev1.CreateUserResponse, error) {
		if err := validateCreateUserRequest(req); err != nil {
			return nil, err
		}

		user := &User{
			Name:      req.Name,
			Password:  req.Password,
			Email:     req.Email,
			CreatedAt: clk.Now(),
			UpdatedAt: clk.Now(),
		}

		insertedID, err := mdb.InsertOne(ctx, MongoDatabaseName, MongoUserCollectionName, user)
		if err != nil {
			return nil, err
		}

		user.ID = insertedID.Hex()

		return &randomdicev1.CreateUserResponse{
			User: &randomdicev1.User{
				Id:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				Password: user.Password,
			},
		}, nil
	}
}

func validateCreateUserRequest(req *randomdicev1.CreateUserRequest) error {
	if req.Name == "" {
		return errEmptyName
	}

	if req.Password == "" {
		return errEmptyPassword
	}

	if req.Email == "" {
		return errEmptyEmail
	}

	return nil
}
