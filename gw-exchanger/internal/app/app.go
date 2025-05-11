package app

import (
	grpcapp "gw-exchanger/internal/app/grpc"
	exchanger "gw-exchanger/internal/services"
	"gw-exchanger/internal/sqlite"
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

	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	exchagerService := exchanger.New(log, nil, storage) // создание сервисного слоя

	grpcApp := grpcapp.New(log, exchagerService, grpcPort) // создание gRPC-приложения

	// возвращаем экземпляр приложения
	return &App{
		GRPCSrv: grpcApp,
	}
}
