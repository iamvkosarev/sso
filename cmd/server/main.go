package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/sso/internal/app/usecase"
	"github.com/iamvkosarev/sso/internal/config"
	"github.com/iamvkosarev/sso/internal/infrastructure/database/postgres"
	server "github.com/iamvkosarev/sso/internal/infrastructure/grpc"
	"github.com/iamvkosarev/sso/internal/infrastructure/grpc/interceptor"
	"github.com/iamvkosarev/sso/internal/infrastructure/http/middleware"
	sqlRepository "github.com/iamvkosarev/sso/internal/infrastructure/repository/postgres"
	"github.com/iamvkosarev/sso/internal/otel/tracing"
	pb "github.com/iamvkosarev/sso/pkg/proto/sso/v1"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	notifyCtx, notifyCtxCancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM, os.Interrupt,
	)
	defer notifyCtxCancel()

	mainCtx, cancelWithCause := context.WithCancelCause(notifyCtx)
	defer cancelWithCause(nil)

	var err error
	shutdownFunc := make([]func(ctx context.Context) error, 0)

	if err = godotenv.Load(); err != nil {
		return err
	}

	cfg := config.MustLoad()
	logger, err := sl.SetupLogger(cfg.Env)
	if err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	shutdown := func(outErr error) error {
		notifyCtxCancel()
		log.Println("Shutting down...")
		shuttingDownCtx, cancel := context.WithTimeout(context.Background(), cfg.App.ShuttingDownTimeout)
		defer cancel()
		if outErr != nil {
			err = errors.Join(outErr, err)
		}
		for _, f := range shutdownFunc {
			if funcErr := f(shuttingDownCtx); funcErr != nil {
				err = errors.Join(funcErr, err)
			}
		}
		log.Println("Shutting down gracefully")
		return err
	}

	log.Println("Starting setup tracer provider")
	tracerProvider, err := tracing.SetupTracerProvider(mainCtx, cfg.OTel.Tracing)
	if err != nil {
		return fmt.Errorf("failed to setup tracer provider: %w", err)
	}
	shutdownFunc = append(shutdownFunc, tracerProvider.Shutdown)

	dns := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SERVICE_NAME"),
		os.Getenv("DB_PORT_INTERNAL"), os.Getenv("DB_NAME"),
	)
	log.Println("Starting Postgres pool connection")
	pool, err := postgres.NewPostgresPool(mainCtx, dns)
	if err != nil {
		return shutdown(fmt.Errorf("failed to connect to postgres: %w", err))
	}
	shutdownFunc = append(
		shutdownFunc, func(ctx context.Context) error {
			pool.Close()
			log.Println("Pool shutdown gracefully")
			return nil
		},
	)

	userRepository := sqlRepository.NewUserRepository(pool)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.RecoveryInterceptor(logger)),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}

	grpcServer := grpc.NewServer(opts...)
	shutdownFunc = append(
		shutdownFunc, func(ctx context.Context) error {
			grpcServer.GracefulStop()
			log.Println("gRPC server shutdown gracefully")
			return nil
		},
	)

	useCase := usecase.NewUserUseCase(userRepository, cfg.App)
	ssoServer := server.NewServer(useCase, logger)

	pb.RegisterSSOServer(grpcServer, ssoServer)

	lis, err := net.Listen("tcp", cfg.Server.GRPCAddress)
	if err != nil {
		return shutdown(fmt.Errorf("failed to listen on %s: %w", cfg.Server.GRPCAddress, err))
	}

	go func() {
		log.Println("Starting gRPC server on", cfg.Server.GRPCAddress)
		if err = grpcServer.Serve(lis); err != nil {
			cancelWithCause(fmt.Errorf("failed to start gRPC server: %w", err))
		}
	}()

	httpMux := http.NewServeMux()

	gwOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	gwMux := runtime.NewServeMux()
	err = pb.RegisterSSOHandlerFromEndpoint(mainCtx, gwMux, cfg.Server.GRPCAddress, gwOpts)
	if err != nil {
		return shutdown(fmt.Errorf("failed to register gateway: %w", err))
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

	server := &http.Server{
		Addr:    cfg.Server.HTTPAddress,
		Handler: otelhttp.NewHandler(middleware.CorsWithOptions(httpMux, cfg.Server.CorsOptions), "http-gateway"),
	}
	shutdownFunc = append(
		shutdownFunc, func(ctx context.Context) error {
			if err = server.Shutdown(ctx); err != nil {
				return shutdown(fmt.Errorf("failed to shutdown HTTP server: %w", err))
			}
			log.Println("HTTP server shutdown gracefully")
			return nil
		},
	)

	go func() {
		log.Printf("Starting HTTP server on %s\n", server.Addr)
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			cancelWithCause(fmt.Errorf("failed to start HTTP server: %w", err))
		}
	}()

	select {
	case <-notifyCtx.Done():
		return shutdown(nil)
	case <-mainCtx.Done():
		return shutdown(mainCtx.Err())
	}
	return nil
}
