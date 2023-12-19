package authgrpc

import (
	"context"
	"errors"
	"log/slog"

	authv1 "local/gorm-example/api/gen/go/auth"
	"local/gorm-example/internal/services/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth auth.Service
	log  *slog.Logger
}

func RegisterServer(gRPCServer *grpc.Server, auth auth.Service, log *slog.Logger) {
	authv1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth, log: log})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	const fn = "authgrpc.Login"
	log := s.log.With("fn", fn)

	if err := req.ValidateAll(); err != nil {
		log.Warn(err.Error())
		return nil, status.Error(codes.InvalidArgument, auth.ErrInvalidCredentials.Error())
	}

	token, err := s.auth.Login(req.GetEmail(), req.GetPassword(), uint(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			log.Warn(err.Error())
			return nil, status.Error(codes.InvalidArgument, auth.ErrInvalidCredentials.Error())
		}
		log.Error(err.Error())
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &authv1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	const fn = "authgrpc.Register"
	log := s.log.With("fn", fn)

	if err := req.ValidateAll(); err != nil {
		log.Warn(err.Error())
		return nil, status.Error(codes.InvalidArgument, auth.ErrInvalidCredentials.Error())
	}

	if err := s.auth.RegisterUser(req.GetName(), req.GetEmail(), req.GetPassword()); err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			log.Warn(err.Error())
			return nil, status.Error(codes.InvalidArgument, auth.ErrUserAlreadyExists.Error())
		}
		log.Error(err.Error())
		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &authv1.RegisterResponse{}, nil
}
