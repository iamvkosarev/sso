package usecase

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/iamvkosarev/sso/internal/config"
	"github.com/iamvkosarev/sso/internal/domain/entity"
	"github.com/iamvkosarev/sso/internal/infrastructure/auth/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type userRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, user entity.User) (int64, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type UserUseCase struct {
	userRepository
	validator *validator.Validate
	app       config.App
}

func NewUserUseCase(repo userRepository, app config.App) *UserUseCase {
	return &UserUseCase{
		userRepository: repo,
		validator:      validator.New(),
		app:            app,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, email, password string) (int64, error) {
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

func (uc *UserUseCase) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := uc.userRepository.GetByEmail(ctx, email)

	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", err
	}

	token, err := jwt.NewToken(user, uc.app.Secret, uc.app.TokenTTL)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (uc *UserUseCase) Verify(_ context.Context, token string) (int64, error) {
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
