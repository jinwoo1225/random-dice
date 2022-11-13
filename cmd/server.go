package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinwoo1225/random-dice/client"
	"github.com/jinwoo1225/random-dice/internal/config"
	"github.com/jinwoo1225/random-dice/internal/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	mdb, mdbCleanup, err := client.NewMongoDBClient(ctx, cfg)
	if err != nil {
		return err
	}
	defer mdbCleanup()

	grpcServer, err := server.NewGRPCServer(cfg, mdb)
	if err != nil {
		return err
	}

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Starting gRPC server... port:", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalln(err)
		}
	}()

	cancelCh := make(chan os.Signal, 1)
	signal.Notify(cancelCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(cancelCh)

	<-cancelCh

	grpcServer.GracefulStop()

	return nil
}
