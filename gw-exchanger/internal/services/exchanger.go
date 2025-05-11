package exchanger

import (
	"context"
	"errors"
	"fmt"
	"gw-exchanger/internal/lib/logger/sl"
	"gw-exchanger/internal/sqlite"
	"gw-exchanger/internal/storage"
	"log/slog"

	walletv1 "github.com/legenda-hortici/protos/gen/go/exchange"
)

// ExchangeSrvc - структура сервисного слоя обмена валют
type ExchangeSrvc struct {
	log          *slog.Logger
	rateProvider RatesProvider
	storage      sqlite.Storage
}

// RatesProvider - интерфейс получения курсов обмена
type RatesProvider interface {
	GetExchangeRateForCurrency(ctx context.Context, req *walletv1.CurrencyRequest) (map[string]float32, error)
	GetExchangeRates(ctx context.Context, _ *walletv1.Empty) (map[string]float32, error)
}

// New - создание сервисного слоя
func New(
	log *slog.Logger,
	rateProvider RatesProvider,
	storage *sqlite.Storage,
) *ExchangeSrvc {
	return &ExchangeSrvc{
		log:          log,
		rateProvider: rateProvider,
		storage:      *storage,
	}
}

// GetExchangeRates реализует grpc-интерфейс RatesProvider
func (e *ExchangeSrvc) GetExchangeRates(ctx context.Context, _ *walletv1.Empty) (*walletv1.ExchangeRatesResponse, error) {
	const op = "exchanger.GetExchangeRates"

	ratesMap, err := e.storage.GetRates(ctx)
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

// GetExchangeRateForCurrency реализует grpc-интерфейс RatesProvider
func (e *ExchangeSrvc) GetExchangeRateForCurrency(ctx context.Context, req *walletv1.CurrencyRequest) (*walletv1.ExchangeRateResponse, error) {
	const op = "exchanger.GetExchangeRateForCurrency"

	if req.FromCurrency == "" || req.ToCurrency == "" {
		return nil, fmt.Errorf("%s: from/to currency is empty", op)
	}

	rate, err := e.storage.GetRate(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			e.log.Warn("rate not found", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		e.log.Error("failed to get rate", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &walletv1.ExchangeRateResponse{
		Rate: rate,
	}, nil
}
