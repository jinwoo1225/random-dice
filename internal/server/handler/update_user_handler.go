package handler

import (
	"context"

	"github.com/pkg/errors"

	"github.com/benbjohnson/clock"
	"github.com/jinwoo1225/random-dice/client"
	randomdicev1 "github.com/jinwoo1225/random-dice/gen/proto/go/randomdice/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errNameNotFound     = errors.New("name not found")
	errEmailNotFound    = errors.New("email not found")
	errPasswordNotFound = errors.New("password not found")
)

type UpdateUserFunc func(ctx context.Context, req *randomdicev1.UpdateUserRequest) (*randomdicev1.UpdateUserResponse, error)

func UpdateUser(clk clock.Clock, mdb client.MongoDBClient) UpdateUserFunc {
	return func(ctx context.Context, req *randomdicev1.UpdateUserRequest) (*randomdicev1.UpdateUserResponse, error) {
		if err := validateUpdateUserRequest(req); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		id := req.GetId()
		name := req.GetName()
		email := req.GetEmail()
		password := req.GetPassword()

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		var user User

		err = mdb.FindOne(ctx, MongoDatabaseName, MongoUserCollectionName, bson.M{"_id": objectID}).Decode(&user)
		if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Error(codes.NotFound, errors.Wrap(err, errNoUserFound.Error()).Error())
		}

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		_, err = mdb.UpdateOne(ctx, MongoDatabaseName, MongoUserCollectionName, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"name": name, "email": email, "password": password, "updated_at": clk.Now()}})
		if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &randomdicev1.UpdateUserResponse{
			User: &randomdicev1.User{
				Id:       id,
				Name:     name,
				Email:    email,
				Password: password,
			},
		}, nil
	}
}

func validateUpdateUserRequest(req *randomdicev1.UpdateUserRequest) error {
	if req.GetId() == "" {
		return errIDNotFound
	}
	if req.GetName() == "" {
		return errNameNotFound
	}
	if req.GetEmail() == "" {
		return errEmailNotFound
	}
	if req.GetPassword() == "" {
		return errPasswordNotFound
	}
	if validObjectID := primitive.IsValidObjectID(req.Id); !validObjectID {
		return errInvalidID
	}
	return nil
}
