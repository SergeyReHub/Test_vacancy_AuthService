package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	DB         DB     `env-prefix:"DB_"`
	HTTPServer Server `env-prefix:"HTTP_"`
	Logger     Logger `env-prefix:"LOGGER_"`
}

type Logger struct {
	Level string `env:"LEVEL"`
}

type Server struct {
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}

type DB struct {
	PortDB       string `env:"PORT"`
	HostDB       string `env:"HOST"`
	NameDB       string `env:"NAME"`
	PasswordDB   string `env:"PASSWORD"`
	UserDB       string `env:"USER"`
	SSLMode      string `env:"SSLMODE"`
	MaxOpenConns int32  `env:"MAX_CONNS"`
}

func (c *DB) ConnStr() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.UserDB, c.PasswordDB, c.HostDB, c.PortDB, c.NameDB, c.SSLMode,
	)
}

func (c *Server) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

var (
	cfg  *Config
	once sync.Once
	mu   sync.Mutex
)

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}

		// Определяем путь к конфигам
		configDir := os.Getenv("CONFIG_DIR")
		if configDir == "" {
			configDir = "internal/config"
		}

		// 1. Загружаем .env из той же папки
		envPath := filepath.Join(configDir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err != nil {
				log.Printf("⚠️ Failed to load .env: %v", err)
			}
		} else {
			log.Printf("ℹ️ .env file not found at %s, using defaults", envPath)
		}

		// Записываем из переменных окружения
		if err := cleanenv.ReadEnv(cfg); err != nil {
			log.Printf("⚠️ Env vars warning: %v", err)
		}
	})
	return cfg
}

// ResetConfig сбрасывает синглтон для тестов
func ResetConfig() {
	mu.Lock()
	defer mu.Unlock()
	once = sync.Once{} // Сбрасываем once
	cfg = nil          // Сбрасываем конфиг
}
