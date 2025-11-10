package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"SERVER_ENV":                    "production",
				"DB_HOST":                       "localhost",
				"DB_PORT":                       "5432",
				"DB_USER":                       "testuser",
				"DB_PASSWORD":                   "testpass",
				"DB_NAME":                       "testdb",
				"DB_SSL_MODE":                   "require",
				"JWT_SECRET":                    "this-is-a-very-long-secret-key-for-testing-purposes-minimum-32-chars",
				"CORS_ALLOWED_ORIGINS":          "http://localhost:3000",
				"RATE_LIMIT_REQUESTS_PER_MINUTE": "60",
				"RATE_LIMIT_REQUESTS_PER_HOUR":  "1000",
			},
			wantErr: false,
		},
		{
			name: "production with SSL disabled - should fail",
			envVars: map[string]string{
				"SERVER_ENV":                    "production",
				"DB_HOST":                       "localhost",
				"DB_PORT":                       "5432",
				"DB_USER":                       "testuser",
				"DB_PASSWORD":                   "testpass",
				"DB_NAME":                       "testdb",
				"DB_SSL_MODE":                   "disable",
				"JWT_SECRET":                    "this-is-a-very-long-secret-key-for-testing-purposes-minimum-32-chars",
				"CORS_ALLOWED_ORIGINS":          "http://localhost:3000",
				"RATE_LIMIT_REQUESTS_PER_MINUTE": "60",
				"RATE_LIMIT_REQUESTS_PER_HOUR":  "1000",
			},
			wantErr: true,
			errMsg:  "DB_SSL_MODE cannot be 'disable' in production",
		},
		{
			name: "production with short JWT secret - should fail",
			envVars: map[string]string{
				"SERVER_ENV":                    "production",
				"DB_HOST":                       "localhost",
				"DB_PORT":                       "5432",
				"DB_USER":                       "testuser",
				"DB_PASSWORD":                   "testpass",
				"DB_NAME":                       "testdb",
				"DB_SSL_MODE":                   "require",
				"JWT_SECRET":                    "short",
				"CORS_ALLOWED_ORIGINS":          "http://localhost:3000",
				"RATE_LIMIT_REQUESTS_PER_MINUTE": "60",
				"RATE_LIMIT_REQUESTS_PER_HOUR":  "1000",
			},
			wantErr: true,
			errMsg:  "JWT_SECRET must be at least 32 characters",
		},
		{
			name: "development with SSL disabled - should pass",
			envVars: map[string]string{
				"SERVER_ENV":                    "development",
				"DB_HOST":                       "localhost",
				"DB_PORT":                       "5432",
				"DB_USER":                       "testuser",
				"DB_PASSWORD":                   "testpass",
				"DB_NAME":                       "testdb",
				"DB_SSL_MODE":                   "disable",
				"JWT_SECRET":                    "short-secret-ok-in-dev",
				"CORS_ALLOWED_ORIGINS":          "http://localhost:3000",
				"RATE_LIMIT_REQUESTS_PER_MINUTE": "60",
				"RATE_LIMIT_REQUESTS_PER_HOUR":  "1000",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Load configuration
			cfg, err := Load()

			// Check error expectation
			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() expected error, got nil")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Load() error = %v, want %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Load() unexpected error = %v", err)
				return
			}

			// Verify configuration
			if cfg == nil {
				t.Fatal("Load() returned nil config")
			}

			// Verify specific values
			if cfg.Server.Env != tt.envVars["SERVER_ENV"] {
				t.Errorf("Server.Env = %v, want %v", cfg.Server.Env, tt.envVars["SERVER_ENV"])
			}

			if cfg.Database.Host != tt.envVars["DB_HOST"] {
				t.Errorf("Database.Host = %v, want %v", cfg.Database.Host, tt.envVars["DB_HOST"])
			}
		})
	}
}

func TestConfig_IsDevelopment(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want bool
	}{
		{"development", "development", true},
		{"dev", "dev", true},
		{"production", "production", false},
		{"test", "test", false},
		{"staging", "staging", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{
					Env: tt.env,
				},
			}

			if got := cfg.IsDevelopment(); got != tt.want {
				t.Errorf("IsDevelopment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want bool
	}{
		{"production", "production", true},
		{"prod", "prod", true},
		{"development", "development", false},
		{"test", "test", false},
		{"staging", "staging", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{
					Env: tt.env,
				},
			}

			if got := cfg.IsProduction(); got != tt.want {
				t.Errorf("IsProduction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		setValue     string
		want         string
	}{
		{
			name:         "existing env var",
			key:          "TEST_VAR",
			defaultValue: "default",
			setValue:     "actual",
			want:         "actual",
		},
		{
			name:         "non-existing env var",
			key:          "NON_EXISTING_VAR",
			defaultValue: "default",
			setValue:     "",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear and set test env var
			os.Unsetenv(tt.key)
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
			}

			got := getEnv(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}

			// Cleanup
			os.Unsetenv(tt.key)
		})
	}
}

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		setValue     string
		want         int
	}{
		{
			name:         "valid integer",
			key:          "TEST_INT",
			defaultValue: 10,
			setValue:     "42",
			want:         42,
		},
		{
			name:         "invalid integer",
			key:          "TEST_INT",
			defaultValue: 10,
			setValue:     "not_a_number",
			want:         10,
		},
		{
			name:         "empty string",
			key:          "TEST_INT",
			defaultValue: 10,
			setValue:     "",
			want:         10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.key)
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
			}

			got := getEnvAsInt(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnvAsInt() = %v, want %v", got, tt.want)
			}

			os.Unsetenv(tt.key)
		})
	}
}

func TestGetEnvAsDuration(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue time.Duration
		setValue     string
		want         time.Duration
	}{
		{
			name:         "valid duration",
			key:          "TEST_DURATION",
			defaultValue: 10 * time.Second,
			setValue:     "30s",
			want:         30 * time.Second,
		},
		{
			name:         "invalid duration",
			key:          "TEST_DURATION",
			defaultValue: 10 * time.Second,
			setValue:     "invalid",
			want:         10 * time.Second,
		},
		{
			name:         "empty string",
			key:          "TEST_DURATION",
			defaultValue: 10 * time.Second,
			setValue:     "",
			want:         10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.key)
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
			}

			got := getEnvAsDuration(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnvAsDuration() = %v, want %v", got, tt.want)
			}

			os.Unsetenv(tt.key)
		})
	}
}

func TestGetEnvAsSlice(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue []string
		setValue     string
		want         []string
	}{
		{
			name:         "comma separated values",
			key:          "TEST_SLICE",
			defaultValue: []string{"default"},
			setValue:     "value1,value2,value3",
			want:         []string{"value1", "value2", "value3"},
		},
		{
			name:         "with spaces",
			key:          "TEST_SLICE",
			defaultValue: []string{"default"},
			setValue:     "value1, value2 , value3",
			want:         []string{"value1", "value2", "value3"},
		},
		{
			name:         "empty string",
			key:          "TEST_SLICE",
			defaultValue: []string{"default"},
			setValue:     "",
			want:         []string{"default"},
		},
		{
			name:         "single value",
			key:          "TEST_SLICE",
			defaultValue: []string{"default"},
			setValue:     "single",
			want:         []string{"single"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.key)
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
			}

			got := getEnvAsSlice(tt.key, tt.defaultValue)

			if len(got) != len(tt.want) {
				t.Errorf("getEnvAsSlice() length = %v, want %v", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("getEnvAsSlice()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}

			os.Unsetenv(tt.key)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid production config",
			cfg: &Config{
				Server: ServerConfig{
					Env: "production",
				},
				Database: DatabaseConfig{
					SSLMode: "require",
				},
				JWT: JWTConfig{
					Secret: "this-is-a-very-long-secret-key-for-testing-purposes-minimum-32-chars",
				},
			},
			wantErr: false,
		},
		{
			name: "production with disabled SSL",
			cfg: &Config{
				Server: ServerConfig{
					Env: "production",
				},
				Database: DatabaseConfig{
					SSLMode: "disable",
				},
				JWT: JWTConfig{
					Secret: "this-is-a-very-long-secret-key-for-testing-purposes-minimum-32-chars",
				},
			},
			wantErr: true,
			errMsg:  "DB_SSL_MODE cannot be 'disable' in production",
		},
		{
			name: "production with short JWT secret",
			cfg: &Config{
				Server: ServerConfig{
					Env: "production",
				},
				Database: DatabaseConfig{
					SSLMode: "require",
				},
				JWT: JWTConfig{
					Secret: "short",
				},
			},
			wantErr: true,
			errMsg:  "JWT_SECRET must be at least 32 characters in production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()

			if tt.wantErr {
				if err == nil {
					t.Error("Validate() expected error, got nil")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Validate() error = %v, want %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error = %v", err)
			}
		})
	}
}
