package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`          // Окружение разработки
	StoragePath string        `yaml:"storage_path" env-required:"true"` // Путь до БД
	TokenTTL    time.Duration `yaml:"token_ttl" env-default:"1h"`       // Время жизни токена
	GRPC        GRPCConfig    `yaml:"grpc"`                             // Конфигурация gRPC
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"44044"`  // Порт
	Timeout time.Duration `yaml:"timeout" env-default:"10h"` // Таймаут
}

// MustLoad загружает и возвращает конфигурацию приложения
func MustLoad() *Config {
	configPath := fetchConfigPath()

	// Проверяем что путь не пустой
	if configPath == "" {
		panic("config path is empty")
	}

	// Проверяем что файл существует
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not found")
	}

	var cfg Config

	// Читаем конфигурационный файл
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

// fetchConfigPath получает путь до конфигурационного файла из командной строки или переменной окружения.
// Приоритет: flag > env > default.
// Если переменная окружения не определена, возвращает пустую строку.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
