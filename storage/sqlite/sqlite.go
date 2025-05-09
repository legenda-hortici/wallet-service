package sqlite

import (
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
	const op = "sqlite.New"

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

// SaveRate - сохранение курса обмена
func (s *Storage) SaveRate(from, to string, rate float64) (int64, error) {
	const op = "sqlite.SaveRate"

	row, err := s.db.Exec("INSERT INTO rates (from, to, rate) VALUES (?, ?, ?)", from, to, rate)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, err := row.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// GetRate - получение курса обмена для конкретной валюты
func (s *Storage) GetRate(from, to string) (float64, error) {
	const op = "sqlite.GetRate"

	var rate float64

	err := s.db.QueryRow("SELECT rate FROM rates WHERE from = ? AND to = ?", from, to).Scan(&rate)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return rate, nil
}

// GetRates - получение курсов обмена
func (s *Storage) GetRates() (map[string]float64, error) {
	const op = "sqlite.GetRates"

	rows, err := s.db.Query("SELECT from, to, rate FROM rates")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rates := make(map[string]float64)

	for rows.Next() {
		var from, to string
		var rate float64

		if err := rows.Scan(&from, &to, &rate); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		rates[from+to] = rate
	}

	return rates, nil
}
