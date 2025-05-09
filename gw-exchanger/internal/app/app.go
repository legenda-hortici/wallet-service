package app

import (
	grpcapp "gw-exchanger/internal/app/grpc"
	"gw-exchanger/internal/services/exchanger"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {

	exchagerService := exchanger.New(log, nil)

	grpcApp := grpcapp.New(log, exchagerService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
