package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string
	LogLevel    string
	VaultConfig VaultConfig
	Database    DatabaseConfig
	ExternalAPI ExternalAPIConfig
	Redis       RedisConfig
	JWT         JWTConfig
	FileStorage FileStorageConfig
	Events      EventsConfig
}

type VaultConfig struct {
	Address string
	Token   string
	Path    string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ExternalAPIConfig struct {
	BaseURL string
	APIKey  string
	Timeout int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Enabled  bool
}

type JWTConfig struct {
	SecretKey string
	Issuer    string
	ExpiryHours int
}

type FileStorageConfig struct {
	Provider    string // "local", "gcs", "s3"
	BucketName  string
	LocalPath   string
	MaxFileSize int64
}

type EventsConfig struct {
	Provider string // "redis", "pubsub", "webhook"
	Topic    string
	WebhookURL string
}

func Load() *Config {
	// Cargar variables de entorno desde .env si existe
	_ = godotenv.Load()

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		VaultConfig: VaultConfig{
			Address: getEnv("VAULT_ADDR", "http://localhost:8200"),
			Token:   getEnv("VAULT_TOKEN", ""),
			Path:    getEnv("VAULT_PATH", "secret/microservice"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "messaging_service"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			Enabled:  getEnvAsBool("REDIS_ENABLED", true),
		},
		JWT: JWTConfig{
			SecretKey:   getEnv("JWT_SECRET", "your-secret-key"),
			Issuer:      getEnv("JWT_ISSUER", "messaging-service"),
			ExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		},
		FileStorage: FileStorageConfig{
			Provider:    getEnv("FILE_STORAGE_PROVIDER", "local"),
			BucketName:  getEnv("FILE_STORAGE_BUCKET", "messaging-attachments"),
			LocalPath:   getEnv("FILE_STORAGE_LOCAL_PATH", "./uploads"),
			MaxFileSize: getEnvAsInt64("FILE_STORAGE_MAX_SIZE", 10*1024*1024), // 10MB
		},
		Events: EventsConfig{
			Provider:   getEnv("EVENTS_PROVIDER", "redis"),
			Topic:      getEnv("EVENTS_TOPIC", "message.events"),
			WebhookURL: getEnv("EVENTS_WEBHOOK_URL", ""),
		},
		ExternalAPI: ExternalAPIConfig{
			BaseURL: getEnv("EXTERNAL_API_URL", "https://api.example.com"),
			APIKey:  getEnv("EXTERNAL_API_KEY", ""),
			Timeout: getEnvAsInt("EXTERNAL_API_TIMEOUT", 30),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}