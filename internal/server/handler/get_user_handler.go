package handler

import (
	"context"
	"errors"

	"github.com/jinwoo1225/random-dice/client"
	randomdicev1 "github.com/jinwoo1225/random-dice/gen/proto/go/randomdice/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errIDNotFound  = errors.New("id not found")
	errInvalidID   = errors.New("invalid ID")
	errNoUserFound = errors.New("no user found with provided ID")
)

type GetUserFunc func(ctx context.Context, req *randomdicev1.GetUserRequest) (*randomdicev1.GetUserResponse, error)

func GetUser(mdb client.MongoDBClient) GetUserFunc {
	return func(ctx context.Context, req *randomdicev1.GetUserRequest) (*randomdicev1.GetUserResponse, error) {
		if err := validateGetUserRequest(req); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		userID := req.GetId()
		userObjectId, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		var userData User

		user, err := mdb.FindOne(ctx, MongoDatabaseName, MongoUserCollectionName, bson.M{"_id": userObjectId})
		err = user.Decode(&userData)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, status.Error(codes.NotFound, errNoUserFound.Error())
			}
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &randomdicev1.GetUserResponse{
			User: &randomdicev1.User{
				Id:       userData.ID,
				Name:     userData.Name,
				Email:    userData.Email,
				Password: userData.Password,
			},
		}, nil
	}
}

func validateGetUserRequest(req *randomdicev1.GetUserRequest) error {
	if req.GetId() == "" {
		return errIDNotFound
	}
	if isObjectID := primitive.IsValidObjectID(req.GetId()); !isObjectID {
		return errInvalidID
	}
	return nil
}
