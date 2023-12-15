package grpcapp

import (
	"fmt"
	"local/gorm-example/internal/config"
	"local/gorm-example/internal/services/product"
	productgrpc "local/gorm-example/internal/services/product/grpc"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	gRPCServer *grpc.Server
	address    int
	log        *slog.Logger
}

func New(product product.Service, log *slog.Logger, cfg *config.GRPCServer) *App {
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

	productgrpc.RegisterServer(gRPCServer, product)

	return &App{
		gRPCServer: gRPCServer,
		address:    cfg.Port,
		log:        log,
	}
}

func (a *App) Run() error {
	const fn = "grpcapp.Run"
	log := a.log.With(slog.String("fn", fn))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.address))

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

	log.Info("stopping gRPC server", slog.Int("port", a.address))

	a.gRPCServer.GracefulStop()

	log.Info("grpc server stopped")
}
