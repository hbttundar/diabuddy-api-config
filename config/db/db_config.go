package dbconfig

import (
	"fmt"
	"github.com/hbttundar/diabuddy-api-config/config/db/dsn"
	"github.com/hbttundar/diabuddy-api-config/config/envmanager"
	diabuddyErrors "github.com/hbttundar/diabuddy-errors"
	"net/url"
	"strings"
)

const (
	DefaultPort     = "5432"
	SslModeQueryKey = "sslmode"
	DefaultSchema   = "postgres"
)

type DBConfig interface {
	ConnectionString() (string, diabuddyErrors.ApiErrors)
}

type PostgresConfig struct {
	dsn        *dsn.DSN
	envManager *envmanager.EnvManager
}

func NewPostgresConfig(envManager *envmanager.EnvManager) (*PostgresConfig, diabuddyErrors.ApiErrors) {
	pc := &PostgresConfig{
		dsn:        &dsn.DSN{},
		envManager: envManager,
	}

	// Ensure environment variables are loaded at the point of initialization
	err := envManager.LoadEnvironmentVariables()
	if err != nil {
		return nil, diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "failed to load environment variables", diabuddyErrors.WithInternalError(err))
	}

	return pc, nil
}

func (pc *PostgresConfig) Get(key string, defaultValue ...string) string {
	return pc.envManager.Get(key, defaultValue...)
}

func (pc *PostgresConfig) ConnectionString() (string, diabuddyErrors.ApiErrors) {
	if err := pc.loadFromEnv(); err != nil {
		return "", err
	}
	return pc.dsn.GenerateConnectionString(DefaultSchema), nil
}

func (pc *PostgresConfig) Validate() diabuddyErrors.ApiErrors {
	if strings.TrimSpace(pc.envManager.Get(envmanager.DbUrlKey)) != "" {
		return nil
	}
	requiredKeys := []string{envmanager.DbHostKey, envmanager.DbPortKey, envmanager.DbDatabaseKey}
	var missingKeys []string

	for _, key := range requiredKeys {
		if strings.TrimSpace(pc.envManager.Get(key)) == "" {
			missingKeys = append(missingKeys, key)
		}
	}
	if len(missingKeys) > 0 {
		return diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, fmt.Sprintf("db configuration is invalid required key(s):%s is/are missed or empty.", strings.Join(missingKeys, ",")))
	}
	return nil
}

func (pc *PostgresConfig) loadFromEnv() diabuddyErrors.ApiErrors {
	if err := pc.Validate(); err != nil {
		return err
	}
	dbUrl := strings.TrimSpace(pc.envManager.Get(envmanager.DbUrlKey))
	if dbUrl != "" {
		return pc.parseDatabaseURL(dbUrl)
	}

	dbHost := pc.envManager.Get(envmanager.DbHostKey)
	dbPort := pc.envManager.Get(envmanager.DbPortKey, DefaultPort)
	dbUsername := pc.envManager.Get(envmanager.DbUsernameKey)
	dbPassword := pc.envManager.Get(envmanager.DbPasswordKey)
	dbDatabase := pc.envManager.Get(envmanager.DbDatabaseKey)
	sslMode := pc.envManager.Get(envmanager.DbSslModeKey)

	if dbHost == "" || dbDatabase == "" || dbPort == "" {
		return diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "missing required db configuration")
	}

	pc.dsn = dsn.NewDsn(dbHost, dbPort, dbUsername, dbPassword, dbDatabase, sslMode)
	return nil
}

func (pc *PostgresConfig) parseDatabaseURL(databaseUrl string) diabuddyErrors.ApiErrors {
	parsedURL, err := url.Parse(databaseUrl)
	if err != nil {
		return diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "db URL is invalid", diabuddyErrors.WithInternalError(err))
	}

	// Extract user info
	dbUsername := parsedURL.User.Username()
	dbPassword, _ := parsedURL.User.Password()

	// Extract host and port
	dbHost := parsedURL.Hostname()
	dbPort := parsedURL.Port()
	if dbPort == "" {
		dbPort = DefaultPort
	}

	// Extract the database name
	dbDatabase := strings.TrimPrefix(parsedURL.Path, "/")

	// Extract sslMode from Query if exist
	var sslMode string
	for key, value := range parsedURL.Query() {
		if strings.ToLower(key) == SslModeQueryKey {
			sslMode = value[0]
		}
	}

	pc.dsn = dsn.NewDsn(dbHost, dbPort, dbUsername, dbPassword, dbDatabase, sslMode)
	return nil
}
