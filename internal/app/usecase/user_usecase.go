package usecase

import (
	"context"
	"errors"
	"github.com/iamvkosarev/sso/internal/config"
	"github.com/iamvkosarev/sso/internal/domain/entity"
	"github.com/iamvkosarev/sso/internal/infrastructure/auth/jwt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const userUseCaseTraceName = "usecase.user"

type userRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, user entity.User) (entity.UserId, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type UserUseCase struct {
	userRepository
	app    config.App
	tracer trace.Tracer
}

func NewUserUseCase(repo userRepository, app config.App) *UserUseCase {
	return &UserUseCase{
		userRepository: repo,
		app:            app,
		tracer:         otel.Tracer(userUseCaseTraceName),
	}
}

func (uc *UserUseCase) Register(ctx context.Context, email, password string) (entity.UserId, error) {
	ctx, span := uc.tracer.Start(ctx, "Register")
	defer span.End()

	user := entity.User{
		Email: email,
	}

	exists, err := uc.userRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, entity.ErrUserAlreadyExists
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user.PassHash = passHash

	id, err := uc.userRepository.Save(ctx, user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (uc *UserUseCase) Login(ctx context.Context, email string, password string) (string, entity.UserId, error) {
	ctx, span := uc.tracer.Start(ctx, "Login")
	defer span.End()

	user, err := uc.userRepository.GetByEmail(ctx, email)

	if err != nil {
		return "", 0, err
	}
	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", 0, err
	}

	token, err := jwt.NewToken(user, uc.app.Secret, uc.app.TokenTTL)
	if err != nil {
		return "", 0, err
	}
	return token, user.Id, nil
}

func (uc *UserUseCase) Verify(ctx context.Context, token string) (int64, error) {
	ctx, span := uc.tracer.Start(ctx, "Verify")
	defer span.End()

	tokenClaims, err := jwt.ParseToken(token, uc.app.Secret)

	if err != nil {
		if errors.Is(err, entity.ErrTokenExpired) {
			return 0, err
		}
		return 0, err
	}
	if tokenClaims.Exp.Before(time.Now()) {
		return 0, err
	}
	return tokenClaims.UserID, nil
}
