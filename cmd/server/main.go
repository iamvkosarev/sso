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
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
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

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.RecoveryInterceptor(logger)),
	}

	grpcServer := grpc.NewServer(opts...)
	useCase := usecase.NewUserUseCase(userRepository, cfg.App)
	ssoServer := server.NewServer(useCase, logger)

	pb.RegisterSSOServer(grpcServer, ssoServer)

	grpcAddr := fmt.Sprintf("0.0.0.0%s", cfg.Server.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", grpcAddr, err)
	}

	go func() {
		log.Println("Starting gRPC server on", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	httpMux := http.NewServeMux()

	grpcGatewayTarget := fmt.Sprintf("localhost%s", cfg.Server.GRPCPort)

	gwOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	gwMux := runtime.NewServeMux()
	err = pb.RegisterSSOHandlerFromEndpoint(ctx, gwMux, grpcGatewayTarget, gwOpts)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	httpMux.Handle("/v1/", gwMux)

	httpMux.HandleFunc(
		cfg.Server.RestPrefix+"/v1/", func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, cfg.Server.RestPrefix)
			r2 := new(http.Request)
			*r2 = *r
			r2.URL.Path = path
			gwMux.ServeHTTP(w, r2)
		},
	)

	httpAddr := fmt.Sprintf("0.0.0.0%s", cfg.Server.RESTPort)
	log.Printf("Starting REST gateway on %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, httpMux); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
