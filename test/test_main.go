package testmain

import (
	"github.com/hbttundar/diabuddy-api-config/config/envmanager"
	"log"
	"os"
	"path/filepath"
)

var (
	EnvKey = []string{
		envmanager.AppNameKey,
		envmanager.AppEnvKey,
		envmanager.AppEncryptionKey,
		envmanager.AppDebugKey,
		envmanager.AppUrlKey,
		envmanager.AppTimezoneKey,
		envmanager.AppLocaleKey,
		envmanager.AppFallbackLocaleKey,
		envmanager.AppCipherKey,
		envmanager.DbUrlKey,
		envmanager.DbHostKey,
		envmanager.DbPortKey,
		envmanager.DbDatabaseKey,
		envmanager.DbUsernameKey,
		envmanager.DbPasswordKey,
		envmanager.DbSslModeKey,
		// other global environment variables
	}

	EnvVars = map[string]string{}
	RootDir string
)

func Setup() {
	log.Println("Setting up default environment variables...")
	SetEnvVars(EnvVars)

}

func TearDown() {
	log.Println("Tearing down environment variables...")
	ClearEnvVars(EnvKey)
	EnvVars = map[string]string{}

}

// SetEnvVars Helper function to set environment variables
func SetEnvVars(vars map[string]string) {
	for key, value := range vars {
		_ = os.Setenv(key, value)
	}
}

// ClearEnvVars Helper function to clear environment variables
func ClearEnvVars(keys []string) {
	for _, key := range keys {
		_ = os.Unsetenv(key)
	}
}

// SetupGoMod renames the go.mod file to go.mod.backup to simulate its absence
func SetupGoMod() {
	goModPath := filepath.Join(RootDir, "go.mod")
	backupPath := filepath.Join(RootDir, "go.mod.backup")

	if _, err := os.Stat(goModPath); err == nil {
		// Rename go.mod to go.mod.backup to simulate its absence
		if err := os.Rename(goModPath, backupPath); err != nil {
			log.Fatalf("Failed to rename go.mod to go.mod.backup: %v", err)
		}
	}
}

// TeardownGoMod renames the go.mod.backup file back to go.mod to restore it
func TeardownGoMod() {
	backupPath := filepath.Join(RootDir, "go.mod.backup")
	goModPath := filepath.Join(RootDir, "go.mod")

	if _, err := os.Stat(backupPath); err == nil {
		// Rename go.mod.backup back to go.mod to restore the original state
		if err := os.Rename(backupPath, goModPath); err != nil {
			log.Fatalf("Failed to rename go.mod.backup back to go.mod: %v", err)
		}
	}
}
