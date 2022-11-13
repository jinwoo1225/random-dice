package handler

import (
	"context"
	"errors"
	"github.com/jinwoo1225/random-dice/client"
	randomdicev1 "github.com/jinwoo1225/random-dice/gen/proto/go/randomdice/v1"
)

var (
	errEmptyName     = errors.New("empty name")
	errEmptyPassword = errors.New("empty password")
	errEmptyEmail    = errors.New("empty email")
)

const (
	MongoDatabaseName   = "random-dice"
	MongoCollectionName = "users"
)

type User struct {
	ID       string `bson:"_id,omitempty"`
	Name     string `bson:"name"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
}

type CreateUserFunc func(ctx context.Context, req *randomdicev1.CreateUserRequest) (*randomdicev1.CreateUserResponse, error)

func CreateUser(mdb *client.MongoDBClient) CreateUserFunc {
	return func(ctx context.Context, req *randomdicev1.CreateUserRequest) (*randomdicev1.CreateUserResponse, error) {
		if err := validateCreateUserRequest(req); err != nil {
			return nil, err
		}

		user := &User{
			Name:     req.Name,
			Password: req.Password,
			Email:    req.Email,
		}

		insertedID, err := mdb.InsertOne(ctx, MongoDatabaseName, MongoCollectionName, user)
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
