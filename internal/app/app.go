package grpcapplication

import (
	grpcapp "github.com/yokoshima228/sso/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort string, storagePath string, tokenTTL time.Duration) *App {
	grpcApp := grpcapp.New(log, grpcPort)
	return &App{
		GRPCServer: grpcApp,
	}
}
