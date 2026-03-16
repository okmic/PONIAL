package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppHost string
	AppEnv  string
	AppPort string
	AppMode string

	ServerHost         string
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration
	AdminSecret        string
	LogLevel           string
	LogFormat          string

	CORSAllowOrigins string
	CORSAllowMethods string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		AppName: getEnv("APP_NAME", "ponial"),
		AppHost: getEnv("APP_HOST", "0.0.0.0"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),
		AppMode: getEnv("APP_MOD", "debug"),

		AdminSecret:        getEnv("ADMIN_SECRET", "TEST"),
		ServerHost:         getEnv("SERVER_HOST", "0.0.0.0"),
		ServerReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second),
		ServerWriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),

		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "text"),

		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "*"),
		CORSAllowMethods: getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "ponial_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}, nil
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func (c *Config) GetCORSAllowOrigins() []string {
	if c.CORSAllowOrigins == "*" {
		return []string{"*"}
	}
	return strings.Split(c.CORSAllowOrigins, ",")
}

func (c *Config) GetCORSAllowMethods() []string {
	return strings.Split(c.CORSAllowMethods, ",")
}
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
