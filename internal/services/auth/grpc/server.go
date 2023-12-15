package authgrpc

import (
	"context"

	authv1 "local/gorm-example/api/gen/go/auth"
	"local/gorm-example/internal/services/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth auth.Service
}

func RegisterServer(gRPCServer *grpc.Server, auth auth.Service) {
	authv1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	// TODO add validation
	token, err := s.auth.Login(req.GetEmail(), req.GetPassword(), uint(req.GetAppId()))

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &authv1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	// TODO add validation
	if err := s.auth.RegisterUser(req.GetName(), req.GetEmail(), req.GetPassword()); err != nil {
		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &authv1.RegisterResponse{}, nil
}
