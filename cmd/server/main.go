package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/sso/internal/app/usecase"
	"github.com/iamvkosarev/sso/internal/config"
	"github.com/iamvkosarev/sso/internal/infrastructure/database/sqlite"
	server "github.com/iamvkosarev/sso/internal/infrastructure/grpc"
	sqlRepository "github.com/iamvkosarev/sso/internal/infrastructure/repository/sqlite"
	pb "github.com/iamvkosarev/sso/pkg/proto/sso/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

func main() {
	cfg := config.MustLoad()
	logger, err := sl.SetupLogger(cfg.Env)
	if err != nil {
		log.Fatalf("error setting up logger: %v\n", err)
	}

	db, err := sqlite.NewDB(cfg.StoragePath)
	if err != nil {
		log.Fatalf("error setting up sqlite: %v\n", err)
	}
	userRepository := sqlRepository.NewUserRepository(db)

	lis, err := net.Listen("tcp", cfg.Server.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	useCase := usecase.NewUserUseCase(userRepository, cfg.App)
	ssoServer := server.NewServer(useCase, logger)

	pb.RegisterSSOServer(grpcServer, ssoServer)

	go func() {
		log.Println("Starting gRPC server on", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux()
	err = pb.RegisterSSOHandlerFromEndpoint(
		ctx, mux, cfg.Server.GRPCPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("Starting REST gateway on", cfg.Server.RESTPort)
	if err := http.ListenAndServe(cfg.Server.RESTPort, mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
