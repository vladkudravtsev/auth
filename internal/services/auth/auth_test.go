package auth

import (
	"log/slog"
	"testing"
	"time"

	"github.com/vladkudravtsev/auth/internal/config"
	"github.com/vladkudravtsev/auth/internal/database"
	"github.com/vladkudravtsev/auth/internal/lib/jwt"
	"github.com/vladkudravtsev/auth/internal/lib/slogdiscard"
	"github.com/vladkudravtsev/auth/internal/models"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	auth Service
	db   *gorm.DB
	cfg  config.Config
}

const (
	secret            = "test-secret"
	name, email, pass = "name", "email", "password"
	deltaInSeconds    = 1
)

var appID uint

func (s *Suite) SetupSuite() {
	cfg := config.LoadFromPath("../../../test.env")
	s.cfg = *cfg

	log := slog.New(slogdiscard.NewDiscardHandler())

	db, err := database.New(cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	if err != nil {
		panic("failed to connect database")
	}
	s.db = db

	authService := NewService(db, log, cfg.Auth.TokenTTL)
	s.auth = authService

	app := &models.App{
		Name:   "test-app",
		Secret: secret,
	}

	res := s.db.Create(&app)
	appID = app.ID
	s.Require().NoError(res.Error)
}

func TestRegisterSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TearDownSuite() {
	s.db.Unscoped().Delete(&models.App{}, appID)
	s.db.Unscoped().Where("email = ?", email).Delete(&models.User{})
}

func (s *Suite) SetupTest() {
	s.db.Unscoped().Where("email = ?", email).Delete(&models.User{})
}

func (s *Suite) TestRegister() {

	var user models.User
	err := s.auth.RegisterUser(name, email, pass)
	s.NoError(err)

	resp := s.db.Where("email = ?", email).First(&user)

	s.NoError(resp.Error)

	s.Equal(email, user.Email)
	s.Equal(name, user.Name)
	s.NotEqual(pass, user.PasswordHash)
}

func (s *Suite) TestRegisterError() {
	s.auth.RegisterUser(name, email, pass)

	err := s.auth.RegisterUser(name, email, pass)
	s.ErrorContains(err, ErrUserAlreadyExists.Error())
}

func (s *Suite) TestLogin() {
	s.auth.RegisterUser(name, email, pass)

	loginTime := time.Now()

	token, err := s.auth.Login(email, pass, appID)
	s.NoError(err)
	s.NotEmpty(token)

	claims, err := jwt.ValidateToken("Bearer "+token, secret)

	s.NoError(err)
	s.Equal(appID, claims.AppID)
	s.Equal(email, claims.Email)

	loginexp := loginTime.Add(s.cfg.Auth.TokenTTL).Unix()
	exp := claims.ExpiresAt.Unix()

	s.InDelta(loginexp, exp, deltaInSeconds)
}

func (s *Suite) TestLoginError() {
	s.auth.RegisterUser(name, email, pass)

	tests := []struct {
		email, pass string
		appID       uint
	}{
		{email: "wrong email", pass: pass, appID: appID},
		{email: email, pass: "wrong pass", appID: appID},
		{email: email, pass: pass, appID: appID + 1},
	}

	for _, tt := range tests {
		token, err := s.auth.Login(tt.email, tt.pass, tt.appID)
		s.ErrorContains(err, ErrInvalidCredentials.Error())
		s.Empty(token)
	}
}
