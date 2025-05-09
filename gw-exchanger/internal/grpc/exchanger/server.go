package exchanger

import (
	"context"

	// Сгенерированный код
	walletv1 "github.com/legenda-hortici/protos/gen/go/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Exchange - интерфейс сервиса
type Exchange interface {
	GetExchangeRates(context.Context, *walletv1.Empty) (*walletv1.ExchangeRatesResponse, error)
	GetExchangeRateForCurrency(context.Context, *walletv1.CurrencyRequest) (*walletv1.ExchangeRateResponse, error)
}

// serverAPI - реализация сервиса
type serverAPI struct {
	walletv1.UnimplementedExchangeServiceServer
	exchange Exchange
}

// Register - создание сервиса
func Register(gRPCServer *grpc.Server, exchange Exchange) {
	walletv1.RegisterExchangeServiceServer(gRPCServer, &serverAPI{exchange: exchange})
}

// GetExchangeRates - получение курсов обмена всех валют
func (s *serverAPI) GetExchangeRates(
	ctx context.Context,
	in *walletv1.Empty,
) (*walletv1.ExchangeRatesResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	rates, err := s.exchange.GetExchangeRates(ctx, in)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get exchange rates")
	}

	return rates, nil
}

// GetExchangeRateForCurrency - получение курса обмена для конкретной валюты
func (s *serverAPI) GetExchangeRateForCurrency(
	ctx context.Context,
	req *walletv1.CurrencyRequest,
) (*walletv1.ExchangeRateResponse, error) {
	if req.FromCurrency == "" || req.ToCurrency == "" {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	rate, err := s.exchange.GetExchangeRateForCurrency(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get exchange rate for currency")
	}

	return rate, nil
}
