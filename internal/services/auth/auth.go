package auth

import (
	"errors"
	"local/gorm-example/internal/lib/jwt"
	"local/gorm-example/internal/services/auth/models"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authService struct {
	db       *gorm.DB
	log      *slog.Logger
	tokenTTL time.Duration
}

func NewService(db *gorm.DB, log *slog.Logger, tokenTTL time.Duration) Service {
	return &authService{db: db, log: log, tokenTTL: tokenTTL}
}

func (s *authService) RegisterUser(name, email, password string) error {
	const fn = "internal.services.auth.RegisterUser"

	log := s.log.With(slog.String("fn", fn))

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}

	if result := s.db.Create(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			log.Warn("user already exists", slog.Any("err", result.Error))
			return ErrUserAlreadyExists
		}
		s.log.Error("user query error", slog.Any("err", result.Error))
		return result.Error
	}

	return nil
}

func (s *authService) Login(email, password string, appID uint) (string, error) {
	const fn = "internal.services.auth.Login"

	log := s.log.With(slog.String("fn", fn))

	var user models.User

	if result := s.db.Where("email = ?", email).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Warn("user not found", slog.Any("err", result.Error))
			return "", ErrInvalidCredentials
		}
		log.Error("user query error", slog.Any("err", result.Error))
		return "", result.Error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Warn("failed to compare password", slog.Any("err", err))
		return "", ErrInvalidCredentials
	}

	var app models.App

	if result := s.db.Where("id = ?", appID).First(&app); result.Error != nil {
		log.Warn("app not found", slog.Any("err", result.Error))
		return "", ErrInvalidCredentials
	}

	claims := &jwt.Claims{
		UserID: user.ID,
		Email:  user.Email,
		AppID:  app.ID,
	}

	token, err := jwt.NewToken(claims, app.Secret, s.tokenTTL)

	if err != nil {
		log.Error("failed to generate token", slog.Any("err", err))
		return "", err
	}

	return token, nil
}
