package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/yokoshima228/sso/internal/domain/models"
	"github.com/yokoshima228/sso/internal/lib/jwt"
	"github.com/yokoshima228/sso/internal/lib/logger/handlers/sl"
	"github.com/yokoshima228/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, password []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int64) (models.App, error)
}

func New(log *slog.Logger, userSaver UserSaver,
	userProvider UserProvider, appProvider AppProvider,
	tokenTTL time.Duration) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    userSaver,
		usrProvider: userProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string, appId int64) (string, error) {
	const position = "services.auth.Login"
	log := a.log.With(
		slog.String("position", position),
		slog.String("email", email),
	)
	log.Info("Attempting to login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%v: %v", position, ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%v: %v", position, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%v: %v", position, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		return "", fmt.Errorf("%v: %v", position, err)
	}

	log.Info("user logged sucessfully")
	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%v: %v", position, err)
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, email, password string) (int64, error) {
	const position = "services.auth.Register"
	log := a.log.With(
		slog.String("position", position),
		slog.String("email", email),
	)

	log.Info("Registering user")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error creating hash", sl.Err(err))
		return 0, fmt.Errorf("%v: %v", err.Error(), position)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, hash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("User already exists", sl.Err(err))
			return 0, fmt.Errorf("%v: %v", position, ErrUserExists)
		}
		log.Error("Error creating user", sl.Err(err))
		return 0, fmt.Errorf("%v: %v", err.Error(), position)
	}
	log.Info("user registered")
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const position = "services.auth.IsAdmin"
	log := a.log.With(
		slog.String("position", position),
		slog.Int64("userId", userId),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found", sl.Err(err))
			return false, fmt.Errorf("%v: %v", position, ErrInvalidAppId)
		}
		return false, fmt.Errorf("%v: %v", position, err)
	}

	log.Info("Checked if user is admin", slog.Bool("isAdmin", isAdmin))
	return isAdmin, nil
}
