package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Storage - хранилище
type Storage struct {
	db *sql.DB
}

// New - инициализация хранилища
func New(storagePath string) (*Storage, error) {
	const op = "storage.grpc.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// Close - закрытие хранилища
func (s *Storage) Close() error {
	return s.db.Close()
}

// GetRate - получение курса обмена для конкретной валюты
func (s *Storage) GetRate(ctx context.Context, from, to string) (float32, error) {
	const op = "sqlite.grpc.GetRate"

	var rate float32

	err := s.db.QueryRow("SELECT rate FROM rates WHERE from = ? AND to = ?", from, to).Scan(&rate)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return rate, nil
}

// GetRates - получение курсов обмена
func (s *Storage) GetRates(ctx context.Context) (map[string]float32, error) {
	const op = "sqlite.grpc.GetRates"

	rows, err := s.db.Query("SELECT from, to, rate FROM rates")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rates := make(map[string]float32)

	for rows.Next() {
		var from, to string
		var rate float32

		if err := rows.Scan(&from, &to, &rate); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		rates[from+to] = rate
	}

	return rates, nil
}
