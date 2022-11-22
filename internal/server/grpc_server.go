package server

import (
	"context"

	"github.com/benbjohnson/clock"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jinwoo1225/random-dice/client"
	randomdicev1 "github.com/jinwoo1225/random-dice/gen/proto/go/randomdice/v1"
	"github.com/jinwoo1225/random-dice/internal/config"
	"github.com/jinwoo1225/random-dice/internal/server/handler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RandomDiceServer struct {
	randomdicev1.UnimplementedUserServiceServer
	randomdicev1.UnimplementedRoomServiceServer

	clk clock.Clock
	cfg *config.Config
	mdb client.MongoDBClient
}

func (s *RandomDiceServer) CreateUser(ctx context.Context, req *randomdicev1.CreateUserRequest) (*randomdicev1.CreateUserResponse, error) {
	return handler.CreateUser(s.clk, s.mdb)(ctx, req)
}

func (s *RandomDiceServer) ListUsers(ctx context.Context, req *randomdicev1.ListUsersRequest) (*randomdicev1.ListUsersResponse, error) {
	return handler.ListUser(s.mdb)(ctx, req)
}

func (s *RandomDiceServer) GetUser(ctx context.Context, req *randomdicev1.GetUserRequest) (*randomdicev1.GetUserResponse, error) {
	return handler.GetUser(s.mdb)(ctx, req)
}

func (s *RandomDiceServer) UpdateUser(ctx context.Context, req *randomdicev1.UpdateUserRequest) (*randomdicev1.UpdateUserResponse, error) {
	return handler.UpdateUser(s.clk, s.mdb)(ctx, req)
}

func (s *RandomDiceServer) DeleteUser(ctx context.Context, req *randomdicev1.DeleteUserRequest) (*randomdicev1.DeleteUserResponse, error) {
	return handler.DeleteUser(s.clk, s.mdb)(ctx, req)
}

func NewRandomDiceServer(
	cfg *config.Config,
	clk clock.Clock,
	mdb client.MongoDBClient,
) (*RandomDiceServer, error) {
	return &RandomDiceServer{
		clk: clk,
		cfg: cfg,
		mdb: mdb,
	}, nil
}

func NewGRPCServer(
	cfg *config.Config,
	logger *zap.Logger,
	clk clock.Clock,
	mdb client.MongoDBClient,
) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(),
		),
	)

	randomDiceServer, err := NewRandomDiceServer(cfg, clk, mdb)
	if err != nil {
		return nil, err
	}

	randomdicev1.RegisterUserServiceServer(grpcServer, randomDiceServer)
	randomdicev1.RegisterRoomServiceServer(grpcServer, randomDiceServer)
	reflection.Register(grpcServer)

	return grpcServer, nil
}
