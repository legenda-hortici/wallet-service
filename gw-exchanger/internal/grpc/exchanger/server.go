package exchanger

import (
	"context"

	// Сгенерированный код
	exchangev1 "github.com/legenda-hortici/protos/gen/go/exchange"
	"google.golang.org/grpc"
)

// Exchange - интерфейс сервиса
type Exchange interface {
	GetExchangeRates(context.Context, *exchangev1.Empty) (*exchangev1.ExchangeRatesResponse, error)
	GetExchangeRateForCurrency(context.Context, *exchangev1.CurrencyRequest) (*exchangev1.ExchangeRateResponse, error)
}

// serverAPI - реализация сервиса
type serverAPI struct {
	exchangev1.UnimplementedExchangeServiceServer
	exchange Exchange
}

// Register - создание сервиса
func Register(gRPCServer *grpc.Server, exchange Exchange) {
	exchangev1.RegisterExchangeServiceServer(gRPCServer, &serverAPI{exchange: exchange})
}

// GetExchangeRates - получение курсов обмена всех валют
func (s *serverAPI) GetExchangeRates(ctx context.Context, in *exchangev1.Empty) (*exchangev1.ExchangeRatesResponse, error) {
	if in == nil {
		return nil, nil
	}

	rates, err := s.exchange.GetExchangeRates(ctx, in)
	if err != nil {
		return nil, err
	}

	return rates, nil
}

// GetExchangeRateForCurrency - получение курса обмена для конкретной валюты
func (s *serverAPI) GetExchangeRateForCurrency(
	ctx context.Context,
	req *exchangev1.CurrencyRequest,
) (*exchangev1.ExchangeRateResponse, error) {
	if req.FromCurrency == "" || req.ToCurrency == "" {
		return nil, nil
	}

	rate, err := s.exchange.GetExchangeRateForCurrency(ctx, req)
	if err != nil {
		return nil, err
	}

	return rate, nil
}
