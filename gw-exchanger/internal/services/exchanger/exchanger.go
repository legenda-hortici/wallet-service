package exchanger

import (
	"context"
	"errors"
	"fmt"
	"gw-exchanger/internal/lib/logger/sl"
	"gw-exchanger/internal/storage"
	"log/slog"

	walletv1 "github.com/legenda-hortici/protos/gen/go/exchange"
)

// ExchangeService - сервис обмена валют
type ExchangeService struct {
	log          *slog.Logger
	rateProvider RatesProvider
}

// RatesProvider - сервис получения курсов обмена
type RatesProvider interface {
	GetExchangeRateForCurrency(ctx context.Context, req *walletv1.CurrencyRequest) (map[string]float32, error)
	GetExchangeRates(ctx context.Context, _ *walletv1.Empty) (map[string]float32, error)
}

// New - создание сервиса
func New(
	log *slog.Logger,
	rateProvider RatesProvider,
) *ExchangeService {
	return &ExchangeService{
		log:          log,
		rateProvider: rateProvider,
	}
}

// GetExchangeRates реализует grpc-интерфейс
func (e *ExchangeService) GetExchangeRates(ctx context.Context, _ *walletv1.Empty) (*walletv1.ExchangeRatesResponse, error) {
	const op = "exchanger.GetExchangeRates"

	ratesMap, err := e.rateProvider.GetExchangeRates(ctx, &walletv1.Empty{})
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			e.log.Warn("rates not found", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		e.log.Error("failed to get rates", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &walletv1.ExchangeRatesResponse{
		Rates: ratesMap,
	}, nil
}

// GetExchangeRateForCurrency реализует grpc-интерфейс
func (e *ExchangeService) GetExchangeRateForCurrency(ctx context.Context, req *walletv1.CurrencyRequest) (*walletv1.ExchangeRateResponse, error) {
	const op = "exchanger.GetExchangeRateForCurrency"

	if req.FromCurrency == "" || req.ToCurrency == "" {
		return nil, fmt.Errorf("%s: from/to currency is empty", op)
	}

	rateMap, err := e.rateProvider.GetExchangeRateForCurrency(ctx, req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			e.log.Warn("rate not found", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		e.log.Error("failed to get rate", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &walletv1.ExchangeRateResponse{
		Rate: rateMap[req.FromCurrency+req.ToCurrency],
	}, nil
}
