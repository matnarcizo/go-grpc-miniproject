package main

import (
	"go-grpc-miniproject/internal/svc"
	greeting_v1 "go-grpc-miniproject/pkg/pb/greeting/v1"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-grpc-miniproject/internal/svc"
	greeting_v1 "github.com/go-grpc-miniproject/pkg/pb/greeting/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	listener net.Listener
	server   *grpc.Server
	logger   *zap.Logger
)

func main() {
	logger, _ = zap.NewProduction()
	defer logger.Sync()

	initListener()

	server = grpc.NewServer()

	greeting_v1.RegisterGreeterServiceServer(server, &svc.GreeterService{})
	logger.Info("Handlers registered")

	go signalsListener(server)

	logger.Info("Starting gRPC server...")
	if err := server.Serve(listener); err != nil {
		logger.Panic("Failed to start gRPC server", zap.Error(err))
	}
}

func initListener() {
	var err error
	address := "localhost:50051"

	listener, err = net.Listen("tcp", address)
	if err != nil {
		logger.Panic("Failed to listen", zap.String("address", address), zap.Error(err))
	}

	logger.Info("Started listening...", zap.String("address", address))

	return
}

func signalsListener(server *grpc.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGNIT)
	_ = <-sigs

	logger.Info("Gracefully stopping server...")
	server.GracefulStop()
}
