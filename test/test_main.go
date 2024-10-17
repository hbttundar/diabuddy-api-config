package testmain

import (
	"github.com/hbttundar/diabuddy-api-config/config/envmanager"
	"log"
	"os"
	"testing"
)

var (
	EnvVars = map[string]string{
		envmanager.AppNameKey:           "",
		envmanager.AppEnvKey:            "",
		envmanager.AppEncryptionKey:     "",
		envmanager.AppDebugKey:          "",
		envmanager.AppUrlKey:            "",
		envmanager.AppTimezoneKey:       "",
		envmanager.AppLocaleKey:         "",
		envmanager.AppFallbackLocaleKey: "",
		envmanager.AppCipherKey:         "",
		envmanager.DbUrlKey:             "",
		envmanager.DbHostKey:            "",
		envmanager.DbPortKey:            "",
		envmanager.DbDatabaseKey:        "",
		envmanager.DbUsernameKey:        "",
		envmanager.DbPasswordKey:        "",
		envmanager.DbSslModeKey:         "",
		// other global environment variables
	}
)

// TestMain sets up and tears down global test setup
func TestMain(m *testing.M) {
	// Global Setup
	Setup()
	err := ensureRequiredFiles()
	if err != nil {
		log.Fatalf("Failed to ensure required files: %v", err)
	}
	// Run migrations to ensure a test database is ready
	err = runMigrations()
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	code := m.Run()
	TearDown()
	// Exit with the appropriate code
	os.Exit(code)
}

func Setup() {
	log.Println("Ensuring empty value for all required env variables in App...")
	// Set all environment variables to empty by default
	for k, v := range EnvVars {
		_ = os.Setenv(k, v)
	}
}

func TearDown() {
	log.Println("unset all env variables...")
	// Unset all environment variables to clean up after tests
	for k := range EnvVars {
		_ = os.Unsetenv(k)
	}
	// Ensure `go.mod` exists
	if _, err := os.Stat("./go.mod"); os.IsNotExist(err) {
		log.Println("go.mod not found, attempting to create or restore...")
		_ = restoreGoMod()
	}

}

// ensureRequiredFiles checks if required files like `go.mod` exist and handles them appropriately
func ensureRequiredFiles() error {
	log.Println("Ensuring required files exist...")

	// Ensure `go.mod` exists
	if _, err := os.Stat("./go.mod"); os.IsNotExist(err) {
		log.Println("go.mod not found, attempting to create or restore...")
		return restoreGoMod()
	}
	return nil
}

// restoreGoMod restores the `go.mod` file if it is missing
func restoreGoMod() error {
	// Assume you have a backup or an example file to restore from
	log.Println("Restoring go.mod file...")
	backupFile := "./go.mod.example"
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return err // Backup also does not exist
	}

	return os.Rename(backupFile, "./go.mod")
}

// runMigrations runs database migrations to prepare for testing
func runMigrations() error {
	log.Println("Running migrations for test database...")

	//TODO as soon as possible create a util service which allows to run migrations using Migrate tool and if migrate too
	//     not found on system show the appropriate message for install it.e.g url to github
	//return util.database.migrations.Run()
	return nil
}
