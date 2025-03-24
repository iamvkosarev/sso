package grpc

import (
	"context"
	"errors"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/sso/internal/domain/entity"
	"github.com/iamvkosarev/sso/internal/infrastructure/auth/jwt"
	pb "github.com/iamvkosarev/sso/pkg/proto/sso/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type UserUseCase interface {
	Register(ctx context.Context, email string, password string) (int64, error)
	Login(ctx context.Context, email string, password string) (string, error)
	Verify(ctx context.Context, token string) (int64, error)
}

type Server struct {
	pb.UnimplementedSSOServer
	*slog.Logger
	userUseCase UserUseCase
}

func NewServer(userUseCase UserUseCase, logger *slog.Logger) *Server {
	return &Server{userUseCase: userUseCase, Logger: logger}
}
func (s *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	const op = "grpc.Server.RegisterUser"

	log := s.Logger.With(
		slog.String("op", op),
	)

	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := s.userUseCase.Register(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		default:
			log.Error("failed to register", err.Error(), slog.String("email", req.Email))
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	log.Info("User registered", slog.String("email", req.Email), slog.Int64("id", id))
	return &pb.RegisterUserResponse{UserId: id}, nil
}

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	const op = "grpc.Server.LoginUser"

	log := s.Logger.With(
		slog.String("op", op),
	)

	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.userUseCase.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			log.Error("failed to login", err.Error(), slog.String("email", req.Email))
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &pb.LoginUserResponse{Token: token}, nil
}

func (s *Server) VerifyToken(ctx context.Context, _ *emptypb.Empty) (*pb.VerifyTokenResponse, error) {
	const op = "grpc.Server.VerifyToken"

	log := s.Logger.With(
		slog.String("op", op),
	)

	token, err := jwt.GetTokenFormContext(ctx)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNoMetadata):
			return nil, status.Error(codes.PermissionDenied, "failed to extract token: empty metadata")
		case errors.Is(err, entity.ErrNoAuthHeader):
			return nil, status.Error(
				codes.PermissionDenied, "failed to extract token: there is no header \"authorization\"",
			)
		case errors.Is(err, entity.ErrInvalidAuthHeader):
			return nil, status.Error(
				codes.PermissionDenied,
				"failed to extract token: not correct \"authorization\" "+
					"value format: correct format is \"Bearer your_token\"",
			)
		default:
			log.Error("failed to extract", err.Error())
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	userID, err := s.userUseCase.Verify(ctx, token)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrTokenIsInvalid):
			return nil, status.Error(codes.PermissionDenied, "invalid token")
		case errors.Is(err, entity.ErrTokenExpired):
			return nil, status.Error(codes.PermissionDenied, "token expired")
		default:
			log.Error("failed to verify token", sl.Err(err))
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &pb.VerifyTokenResponse{UserId: userID}, nil
}
