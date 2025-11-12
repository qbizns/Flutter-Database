package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Session  SessionConfig
	RateLimit RateLimitConfig
	Logging  LoggingConfig
	CORS     CORSConfig
	Worker   WorkerConfig
	Upload   UploadConfig
	Email    EmailConfig
	SMS      SMSConfig
	Metrics  MetricsConfig
}

type ServerConfig struct {
	APIPort  string
	GRPCPort string
	Env      string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	SSLRootCert     string
	SSLCert         string
	SSLKey          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	QueryTimeout    time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret                 string
	AccessTokenDuration    time.Duration
	RefreshTokenDuration   time.Duration
}

type SessionConfig struct {
	TimeoutMinutes int
}

type RateLimitConfig struct {
	RequestsPerMinute int
	RequestsPerHour   int
}

type LoggingConfig struct {
	Level  string
	Format string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type WorkerConfig struct {
	Concurrency  int
	PollInterval time.Duration
}

type UploadConfig struct {
	MaxSizeMB int
	UploadDir string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type SMSConfig struct {
	Provider            string
	TwilioAccountSID    string
	TwilioAuthToken     string
	TwilioPhoneNumber   string
}

type MetricsConfig struct {
	Enabled     bool
	MetricsPort string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists (development)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			APIPort:  getEnv("API_PORT", "8080"),
			GRPCPort: getEnv("GRPC_PORT", "9090"),
			Env:      getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "pos_saas"),
			SSLMode:         getEnv("DB_SSL_MODE", "require"),
			SSLRootCert:     getEnv("DB_SSL_ROOT_CERT", ""),
			SSLCert:         getEnv("DB_SSL_CERT", ""),
			SSLKey:          getEnv("DB_SSL_KEY", ""),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute),
			QueryTimeout:    getEnvAsDuration("DB_QUERY_TIMEOUT", 30*time.Second),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:                 getEnv("JWT_SECRET", "your-secret-key-change-me"),
			AccessTokenDuration:    getEnvAsDuration("JWT_ACCESS_TOKEN_DURATION", 15*time.Minute),
			RefreshTokenDuration:   getEnvAsDuration("JWT_REFRESH_TOKEN_DURATION", 7*24*time.Hour),
		},
		Session: SessionConfig{
			TimeoutMinutes: getEnvAsInt("SESSION_TIMEOUT_MINUTES", 480),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getEnvAsInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
			RequestsPerHour:   getEnvAsInt("RATE_LIMIT_REQUESTS_PER_HOUR", 1000),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"}),
		},
		Worker: WorkerConfig{
			Concurrency:  getEnvAsInt("WORKER_CONCURRENCY", 10),
			PollInterval: getEnvAsDuration("WORKER_POLL_INTERVAL", 5*time.Second),
		},
		Upload: UploadConfig{
			MaxSizeMB: getEnvAsInt("MAX_UPLOAD_SIZE_MB", 10),
			UploadDir: getEnv("UPLOAD_DIR", "./uploads"),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", ""),
			SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("SMTP_FROM_EMAIL", "noreply@example.com"),
			FromName:     getEnv("SMTP_FROM_NAME", "POS System"),
		},
		SMS: SMSConfig{
			Provider:          getEnv("SMS_PROVIDER", "twilio"),
			TwilioAccountSID:  getEnv("TWILIO_ACCOUNT_SID", ""),
			TwilioAuthToken:   getEnv("TWILIO_AUTH_TOKEN", ""),
			TwilioPhoneNumber: getEnv("TWILIO_PHONE_NUMBER", ""),
		},
		Metrics: MetricsConfig{
			Enabled:     getEnvAsBool("ENABLE_METRICS", true),
			MetricsPort: getEnv("METRICS_PORT", "9091"),
		},
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Validate checks if required configuration is present
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	// Validate SSL mode
	validSSLModes := map[string]bool{
		"disable": true, "allow": true, "prefer": true,
		"require": true, "verify-ca": true, "verify-full": true,
	}
	if !validSSLModes[c.Database.SSLMode] {
		return fmt.Errorf("invalid DB_SSL_MODE: %s (valid: disable, allow, prefer, require, verify-ca, verify-full)", c.Database.SSLMode)
	}

	// Production security requirements
	if c.Server.Env == "production" {
		if c.JWT.Secret == "your-secret-key-change-me" {
			return fmt.Errorf("JWT_SECRET must be changed in production")
		}
		if len(c.JWT.Secret) < 32 {
			return fmt.Errorf("JWT_SECRET must be at least 32 characters in production")
		}
		if c.Database.SSLMode == "disable" {
			return fmt.Errorf("DB_SSL_MODE cannot be 'disable' in production - use 'require' or higher")
		}
		if c.Database.SSLMode == "allow" || c.Database.SSLMode == "prefer" {
			return fmt.Errorf("DB_SSL_MODE must be 'require' or higher in production (current: %s)", c.Database.SSLMode)
		}
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	// Parse comma-separated values properly
	var result []string
	parts := strings.Split(valueStr, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}
	return result
}
