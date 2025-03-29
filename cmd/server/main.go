package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/sso/internal/app/usecase"
	"github.com/iamvkosarev/sso/internal/config"
	"github.com/iamvkosarev/sso/internal/infrastructure/database/postgres"
	server "github.com/iamvkosarev/sso/internal/infrastructure/grpc"
	"github.com/iamvkosarev/sso/internal/infrastructure/grpc/interceptor"
	sqlRepository "github.com/iamvkosarev/sso/internal/infrastructure/repository/postgres"
	pb "github.com/iamvkosarev/sso/pkg/proto/sso/v1"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found or failed to load")
	}

	cfg := config.MustLoad()
	logger, err := sl.SetupLogger(cfg.Env)
	if err != nil {
		log.Fatalf("error setting up logger: %v\n", err)
	}

	ctx := context.Background()

	dns := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SERVICE_NAME"),
		os.Getenv("DB_PORT_INTERNAL"), os.Getenv("DB_NAME"),
	)
	pool, err := postgres.NewPostgresPool(ctx, dns)
	if err != nil {
		log.Fatalf("error setting up postgres: %v\n", err)
	}
	userRepository := sqlRepository.NewUserRepository(pool)

	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(interceptor.RecoveryInterceptor(logger)))

	if cfg.Server.TLSEnabled {
		certFile := os.Getenv("CERT_FILE")
		keyFile := os.Getenv("KEY_FILE")
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)
	useCase := usecase.NewUserUseCase(userRepository, cfg.App)
	ssoServer := server.NewServer(useCase, logger)

	pb.RegisterSSOServer(grpcServer, ssoServer)

	lis, err := net.Listen("tcp", cfg.Server.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("Starting gRPC server on", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	httpMux := http.NewServeMux()
	gwMux := runtime.NewServeMux()

	err = pb.RegisterSSOHandlerFromEndpoint(
		ctx, gwMux, cfg.Server.GRPCPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	httpMux.Handle(cfg.Server.RestPrefix+"/", http.StripPrefix(cfg.Server.RestPrefix, gwMux))
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("Starting REST gateway on", cfg.Server.RESTPort)
	if cfg.Server.TLSEnabled {
		certFile := os.Getenv("CERT_FILE")
		keyFile := os.Getenv("KEY_FILE")
		if err := http.ListenAndServeTLS(
			cfg.Server.RESTPort, certFile, keyFile, httpMux,
		); err != nil {
			log.Fatalf("failed to serve HTTPS: %v", err)
		}
	} else {
		if err := http.ListenAndServe(cfg.Server.RESTPort, httpMux); err != nil {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}
}
