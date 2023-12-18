package tests

import (
	"context"
	authv1 "local/gorm-example/api/gen/go/auth"
	"local/gorm-example/internal/config"
	"local/gorm-example/internal/database"
	"local/gorm-example/internal/lib/jwt"
	"local/gorm-example/internal/lib/slogdiscard"
	"local/gorm-example/internal/services/auth"
	authgrpc "local/gorm-example/internal/services/auth/grpc"
	"local/gorm-example/internal/services/auth/models"
	"log/slog"
	"net"
	"time"

	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	Cfg    *config.Config
	client authv1.AuthClient
	db     *gorm.DB
}

const (
	grpcHost       = "localhost"
	secret         = "test-secret"
	bufSize        = 1024 * 1024
	deltaInSeconds = 1
)

var appID uint

var user *authv1.RegisterRequest

func (s *Suite) SetupSuite() {
	cfg := config.LoadFromPath("../test.env")
	s.Cfg = cfg

	log := slog.New(slogdiscard.NewDiscardHandler())

	db, err := database.New(cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)

	if err != nil {
		panic("failed to connect database")
	}
	s.db = db

	lis := bufconn.Listen(bufSize)
	baseServer := grpc.NewServer()

	authService := auth.NewService(db, log, cfg.Auth.TokenTTL)
	authgrpc.RegisterServer(baseServer, authService)

	go func() {
		if err := baseServer.Serve(lis); err != nil {
			s.Failf("grpc server start failed: %v", err.Error())
		}
	}()

	cc, err := grpc.DialContext(context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		s.Failf("grpc server connection failed: %v", err.Error())
	}

	s.client = authv1.NewAuthClient(cc)

	app := &models.App{
		Name:   "test-app",
		Secret: secret,
	}

	res := s.db.Create(&app)
	s.Require().NoError(res.Error)
	appID = app.ID

	user = &authv1.RegisterRequest{
		Name: "user123", Email: "testuser123@gmail.com", Password: "vlad12345",
	}
}

func (s *Suite) TearDownSuite() {
	s.db.Exec("TRUNCATE users")
	s.db.Exec("TRUNCATE apps")
}

func (s *Suite) SetupTest() {
	s.db.Exec("TRUNCATE users")
}

func TestRegisterSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestRegister() {
	_, err1 := s.client.Register(context.Background(), user)
	_, err2 := s.client.Register(context.Background(), user)

	s.NoError(err1)
	s.Error(err2)
}

func (s *Suite) TestLogin() {
	_, err := s.client.Register(context.Background(), user)
	s.Require().NoError(err)

	login := &authv1.LoginRequest{
		Email:    user.Email,
		Password: user.Password,
		AppId:    uint32(appID),
	}

	loginTime := time.Now()

	resp, err := s.client.Login(context.Background(), login)
	s.NoError(err)

	claims, err := jwt.ValidateToken("Bearer "+resp.GetToken(), secret)

	s.NoError(err)
	s.Equal(appID, claims.AppID)
	s.Equal(user.Email, claims.Email)

	loginexp := loginTime.Add(s.Cfg.Auth.TokenTTL).Unix()
	exp := claims.ExpiresAt.Unix()

	s.InDelta(loginexp, exp, deltaInSeconds)
}

func (s *Suite) TestLogin_InvalidCredentials() {
	_, err := s.client.Register(context.Background(), user)
	s.Require().NoError(err)

	tests := []struct {
		input authv1.LoginRequest
	}{
		{
			input: authv1.LoginRequest{
				Email:    user.Email,
				Password: "wrong password",
				AppId:    uint32(appID),
			},
		},
		{
			input: authv1.LoginRequest{
				Email:    "wrong email",
				Password: user.Email,
				AppId:    uint32(appID),
			},
		},
		{
			input: authv1.LoginRequest{
				Email:    user.Email,
				Password: user.Email,
				AppId:    uint32(appID) + 1,
			},
		},
	}
	for _, tt := range tests {
		token, err := s.client.Login(context.Background(), &tt.input)

		s.Error(err, "should be invalid credentials error")
		s.Empty(token, "token should be empty")
	}
}
