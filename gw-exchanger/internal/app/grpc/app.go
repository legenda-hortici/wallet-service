package grpcapp

import (
	"fmt"
	exchangergrpc "gw-exchanger/internal/grpc"
	exchanger "gw-exchanger/internal/services"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

// App - сервер gRPC
type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New - создание сервера gRPC
func New(
	log *slog.Logger,
	exchangeGrpc *exchanger.ExchangeSrvc,
	port int,
) *App {
	grpcServer := grpc.NewServer()

	exchangergrpc.Register(grpcServer, exchangeGrpc)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
	}
}

// MustRun - запуск приложения. Если возникла ошибка, то приложение завершает работу принудительно
func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// Run - запуск приложения
func (app *App) Run() error {
	const op = "grpcapp.Run"

	log := app.log.With(
		slog.String("op", op),      // операция
		slog.Int("port", app.port), // порт
	)

	log.Info("starting gRPC server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port)) // слушаем tcp соединение
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("address", l.Addr().String()))

	if err := app.gRPCServer.Serve(l); err != nil { // запускаем gRPC-сервер
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop - остановка приложения
func (app *App) Stop() {
	const op = "grpcapp.Stop"

	app.log.With(slog.String("op", op)).Info("stopping gRPC server")

	app.gRPCServer.GracefulStop()
}
