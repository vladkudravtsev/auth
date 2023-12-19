package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/vladkudravtsev/auth/internal/config"
	authgrpc "github.com/vladkudravtsev/auth/internal/grpc"
	"github.com/vladkudravtsev/auth/internal/services/auth"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(authService auth.Service, log *slog.Logger, cfg *config.GRPCServer) *App {
	recoveryOptions := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOptions...),
		),
	)

	authgrpc.RegisterServer(gRPCServer, authService, log)

	return &App{
		gRPCServer: gRPCServer,
		log:        log,
		port:       cfg.Port,
	}
}

func (a *App) Run() error {
	const fn = "grpcapp.Run"

	log := a.log.With(slog.String("fn", fn))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("gRPC server is running", slog.String("port", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	const fn = "grpcapp.stop"

	log := a.log.With(slog.String("fn", fn))

	log.Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()

	log.Info("grpc server stopped")
}
