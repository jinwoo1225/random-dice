package server

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jinwoo1225/random-dice/client"
	randomdicev1 "github.com/jinwoo1225/random-dice/gen/proto/go/randomdice/v1"
	"github.com/jinwoo1225/random-dice/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RandomDiceServer struct {
	randomdicev1.UnimplementedUserServiceServer
	randomdicev1.UnimplementedRoomServiceServer

	cfg *config.Config
	mdb *client.MongoDBClient
}

func NewRandomDiceServer(cfg *config.Config, mdb *client.MongoDBClient) (*RandomDiceServer, error) {
	return &RandomDiceServer{
		cfg: cfg,
		mdb: mdb,
	}, nil
}

func NewGRPCServer(
	cfg *config.Config,
	mdb *client.MongoDBClient,
) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpc_recovery.UnaryServerInterceptor()),
	)

	randomDiceServer, err := NewRandomDiceServer(cfg, mdb)
	if err != nil {
		return nil, err
	}

	randomdicev1.RegisterUserServiceServer(grpcServer, randomDiceServer)
	randomdicev1.RegisterRoomServiceServer(grpcServer, randomDiceServer)
	reflection.Register(grpcServer)

	return grpcServer, nil
}
