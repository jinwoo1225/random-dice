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

type DeleteUserFunc func(ctx context.Context, req *randomdicev1.DeleteUserRequest) (*randomdicev1.DeleteUserResponse, error)

func DeleteUser(
	clk clock.Clock,
	mdb client.MongoDBClient,
) DeleteUserFunc {
	return func(ctx context.Context, req *randomdicev1.DeleteUserRequest) (*randomdicev1.DeleteUserResponse, error) {
		if err := validateDeleteUserRequest(req); err != nil {
			return nil, err
		}

		var user User

		objectID, err := primitive.ObjectIDFromHex(req.Id)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		err = mdb.FindOne(ctx, MongoDatabaseName, MongoUserCollectionName, bson.M{"_id": objectID}).Decode(&user)
		if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Error(codes.NotFound, errors.Wrap(err, errNoUserFound.Error()).Error())
		}

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		res, err := mdb.UpdateOne(ctx, MongoDatabaseName, MongoUserCollectionName, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"deleted_at": clk.Now()}})
		if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		var deleted bool
		if res.ModifiedCount > 0 {
			deleted = true
		}

		return &randomdicev1.DeleteUserResponse{
			Success: deleted,
		}, nil
	}
}

func validateDeleteUserRequest(req *randomdicev1.DeleteUserRequest) error {
	if req.Id == "" {
		return errInvalidID
	}

	if validObjectID := primitive.IsValidObjectID(req.Id); !validObjectID {
		return errInvalidID
	}

	return nil
}
