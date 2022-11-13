package handler

import (
	"context"
	"errors"
	"github.com/jinwoo1225/random-dice/client"
	randomdicev1 "github.com/jinwoo1225/random-dice/gen/proto/go/randomdice/v1"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errInvalidPage  = errors.New("invalid page")
	errInvalidLimit = errors.New("invalid limit")
)

type ListUsersFunc func(ctx context.Context, req *randomdicev1.ListUsersRequest) (*randomdicev1.ListUsersResponse, error)

func ListUser(mdb client.MongoDBClient) ListUsersFunc {
	return func(ctx context.Context, req *randomdicev1.ListUsersRequest) (*randomdicev1.ListUsersResponse, error) {
		if err := validateListUsersRequest(req); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		sortOpt := bson.M{
			"created_at": -1,
		}

		cur, err := mdb.FindMany(ctx, MongoDatabaseName, MongoUserCollectionName, bson.M{}, sortOpt, req.Page, req.Limit)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		var users []User
		if err := cur.All(ctx, &users); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		var userResponses []*randomdicev1.User
		for _, user := range users {
			userResponses = append(userResponses, &randomdicev1.User{
				Id:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				Password: user.Password,
			})
		}

		return &randomdicev1.ListUsersResponse{
			Users: userResponses,
		}, nil
	}
}

func validateListUsersRequest(req *randomdicev1.ListUsersRequest) error {
	if req.Page < 0 {
		return errInvalidPage
	}

	if req.Limit <= 0 {
		return errInvalidLimit
	}

	return nil
}
