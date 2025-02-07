package grpcapp

import (
	"fmt"
	authgrpc "github.com/yokoshima228/sso/internal/grpc/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(log *slog.Logger, port string) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const position = "internal.app.grpc.Run"
	log := a.log.With(slog.String("position", position), slog.String("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return fmt.Errorf("%v: %v", position, err)
	}

	log.Info("gRPCServer is running", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%v: %v", position, err)
	}

	return nil
}

func (a *App) Stop() {
	const position = "internal.app.grpc.Stop"

	a.log.With(slog.String("position", position)).
		Info("stopping gRPC server", slog.String("port", a.port))
	a.gRPCServer.GracefulStop()
}
