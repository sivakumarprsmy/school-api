package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/periasamy/school-api/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "local", cfg.Environment)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, 5432, cfg.DBPort)
	assert.Equal(t, "school", cfg.DBName)
	assert.Empty(t, cfg.DBUser)
	assert.Empty(t, cfg.DBPassword)
}

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("PORT", "9090")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_NAME", "mydb")
	t.Setenv("DB_USER", "user")
	t.Setenv("DB_PASSWORD", "secret")

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "production", cfg.Environment)
	assert.Equal(t, 9090, cfg.Port)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "db.example.com", cfg.DBHost)
	assert.Equal(t, 5433, cfg.DBPort)
	assert.Equal(t, "mydb", cfg.DBName)
	assert.Equal(t, "user", cfg.DBUser)
	assert.Equal(t, "secret", cfg.DBPassword)
}

func TestLoad_InvalidPort_FallsBackToDefault(t *testing.T) {
	t.Setenv("PORT", "not-a-number")

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, 8080, cfg.Port)
}

func TestPostgresDSN(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "admin",
		DBPassword: "password",
		DBName:     "school",
	}

	dsn := cfg.PostgresDSN()

	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "port=5432")
	assert.Contains(t, dsn, "user=admin")
	assert.Contains(t, dsn, "password=password")
	assert.Contains(t, dsn, "dbname=school")
	assert.Contains(t, dsn, "sslmode=disable")
}
