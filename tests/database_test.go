package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"idm/inner/common"
	"idm/inner/database"
	"os"
	"path/filepath"
	"testing"
)

// TestConnectDb tests the ConnectDb and ConnectDbWithCfg functions for various scenarios.
func TestConnectDb(t *testing.T) {
	// Helper function to create a temporary .env file
	createEnvFile := func(content string) string {
		tmpDir := t.TempDir()
		envFile := filepath.Join(tmpDir, ".env")
		if content != "" {
			err := os.WriteFile(envFile, []byte(content), 0644)
			require.NoError(t, err, "Failed to create .env file")
		}

		return envFile
	}

	// Helper function to clear environment variables
	clearEnv := func() {
		_ = os.Unsetenv("DB_CONNECTION")
		_ = os.Unsetenv("DB_USER")
		_ = os.Unsetenv("DB_PASSWORD")
		_ = os.Unsetenv("DB_HOST")
		_ = os.Unsetenv("DB_PORT")
		_ = os.Unsetenv("DB_NAME")
	}

	// Helper function to set environment variables
	setEnv := func(conn, user, pass, host, port, name string) {
		_ = os.Setenv("DB_CONNECTION", conn)
		_ = os.Setenv("DB_USER", user)
		_ = os.Setenv("DB_PASSWORD", pass)
		_ = os.Setenv("DB_HOST", host)
		_ = os.Setenv("DB_PORT", port)
		_ = os.Setenv("DB_NAME", name)
	}

	t.Run("1. No .env file", func(t *testing.T) {
		clearEnv()
		cfg := common.GetConfig("non_existent.env")

		assert.Empty(t, cfg.DbDriverName, "DbDriverName should be empty")
		assert.Equal(t, "://:@:/?sslmode=disable", cfg.Dsn, "Dsn should be empty but formatted")
	})

	t.Run("2. .env file exists but no required variables, no env vars", func(t *testing.T) {
		clearEnv()
		envFile := createEnvFile("")
		cfg := common.GetConfig(envFile)

		assert.Empty(t, cfg.DbDriverName, "DbDriverName should be empty")
		assert.Equal(t, "://:@:/?sslmode=disable", cfg.Dsn, "Dsn should be empty but formatted")
	})

	t.Run("3. .env file exists but no variables, env vars present", func(t *testing.T) {
		clearEnv()
		setEnv("postgres", "testuser", "testpass", "localhost", "5432", "testdb")
		envFile := createEnvFile("")
		cfg := common.GetConfig(envFile)

		// Environment variables should be used
		expectedDsn := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
		assert.Equal(t, "postgres", cfg.DbDriverName, "DbDriverName should match env var")
		assert.Equal(t, expectedDsn, cfg.Dsn, "Dsn should match env vars")
	})

	t.Run("4. .env file with variables, conflicting env vars", func(t *testing.T) {
		clearEnv()
		setEnv("postgres", "envuser", "envpass", "envhost", "5432", "envdb")
		envContent := `DB_CONNECTION=postgres
			DB_USER=dotuser
			DB_PASSWORD=dotpass
			DB_HOST=dothost
			DB_PORT=5433
			DB_NAME=dotdb`
		envFile := createEnvFile(envContent)
		cfg := common.GetConfig(envFile)

		// .env file should NOT take precedence (godotenv.Overload overrides env vars)
		expectedDsn := "postgres://envuser:envpass@envhost:5432/envdb?sslmode=disable"
		assert.Equal(t, "postgres", cfg.DbDriverName, "DbDriverName should match .env")
		assert.Equal(t, expectedDsn, cfg.Dsn, "Dsn should match env vars")
	})

	t.Run("5. Correct .env file, no conflicting env vars", func(t *testing.T) {
		clearEnv()
		envContent := `DB_CONNECTION=postgres
			DB_USER=testuser
			DB_PASSWORD=testpass
			DB_HOST=localhost
			DB_PORT=5432
			DB_NAME=testdb`
		envFile := createEnvFile(envContent)
		cfg := common.GetConfig(envFile)

		// .env file should be used
		expectedDsn := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
		assert.Equal(t, "postgres", cfg.DbDriverName, "DbDriverName should match .env")
		assert.Equal(t, expectedDsn, cfg.Dsn, "Dsn should match .env")
	})

	t.Run("6. Cannot connect with incorrect config", func(t *testing.T) {
		// Simulate an incorrect config
		cfg := common.Config{
			DbDriverName: "postgres",
			Dsn:          "postgres://wronguser:wrongpass@invalidhost:9999/invaliddb?sslmode=disable",
		}
		// Since sqlx.MustConnect panics on failure, we catch the panic
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for invalid config, but none occurred")
			}
		}()
		database.ConnectDbWithCfg(cfg)
	})

	t.Run("7. Can connect with correct config", func(t *testing.T) {
		clearEnv()
		db := database.ConnectDbWithCfg(common.GetConfig(".env"))

		assert.NotNil(t, db, "DB should not be nil")
	})
}
