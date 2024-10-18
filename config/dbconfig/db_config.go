package dbconfig

import (
	"fmt"
	"github.com/hbttundar/diabuddy-api-config/config/dbconfig/dsn"
	"github.com/hbttundar/diabuddy-api-config/config/envmanager"
	diabuddyErrors "github.com/hbttundar/diabuddy-errors"
	"net/url"
	"strings"
)

const (
	Mysql                = "mysql"
	SqlServer            = "sqlserver"
	Oracle               = "oracle"
	MongoDb              = "mongodb"
	Redis                = "redis"
	Cassandra            = "cassandra"
	Postgres             = "postgres"
	MysqlDefaultPort     = "3306"
	SqlServerDefaultPort = "1433"
	OracleDefaultPort    = "1521"
	MongoDbDefaultPort   = "27017"
	RedisDefaultPort     = "6379"
	CassandraDefaultPort = "9042"
	PostgresDefaultPort  = "5432"
)

type Config interface {
	ConnectionString() (string, diabuddyErrors.ApiErrors)
	Validate(requiredKeys []string) diabuddyErrors.ApiErrors
	Get(key string, defaultValue ...string) string
}

type DBConfig struct {
	envManager  *envmanager.EnvManager
	dsn         *dsn.DSN
	dbType      string
	params      map[string]string
	isTypeSet   bool
	isParamsSet bool
}

// ConfigOption Option function type for configuring DbConfig.
type ConfigOption func(*DBConfig) diabuddyErrors.ApiErrors

// NewDBConfig creates a new DbConfig with provided options.
func NewDBConfig(envManager *envmanager.EnvManager, options ...ConfigOption) (*DBConfig, diabuddyErrors.ApiErrors) {
	databaseConfig := &DBConfig{
		envManager: envManager,
		dbType:     Postgres, // Default to Postgres if not specified
		params:     make(map[string]string),
	}

	// Apply provided options
	for _, opt := range options {
		if err := opt(databaseConfig); err != nil {
			return nil, err
		}
	}

	// Ensure environment variables are loaded at the point of initialization
	if err := envManager.LoadEnvironmentVariables(); err != nil {
		return nil, diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "failed to load environment variables", diabuddyErrors.WithInternalError(err))
	}

	return databaseConfig, nil
}

// WithType sets the type of database (e.g., postgres, mysql, etc.).
func WithType(dbType string) ConfigOption {
	return func(dbConfig *DBConfig) diabuddyErrors.ApiErrors {
		if dbConfig.isTypeSet {
			return diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "multiple database types provided")
		}
		dbConfig.dbType = dbType
		dbConfig.isTypeSet = true
		return nil
	}
}

// WithDsnParameters sets DSN parameters for DBConfig.
func WithDsnParameters(params map[string]string) ConfigOption {
	return func(dbConfig *DBConfig) diabuddyErrors.ApiErrors {
		if dbConfig.isParamsSet {
			return diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "multiple DSN parameter sets provided")
		}
		dbConfig.params = params
		dbConfig.isParamsSet = true
		return nil
	}
}

// Get retrieves the value of an environment variable, or the default if not set.
func (c *DBConfig) Get(key string, defaultValue ...string) string {
	return c.envManager.Get(key, defaultValue...)
}

// Validate checks that all required environment variables are present.
func (c *DBConfig) Validate() diabuddyErrors.ApiErrors {
	requiredKeys := getRequiredKeysForDBType(c.dbType)

	var missingKeys []string
	for _, key := range requiredKeys {
		if strings.TrimSpace(c.envManager.Get(key)) == "" {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		return diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, fmt.Sprintf("missing required key(s): %s", strings.Join(missingKeys, ", ")))
	}
	return nil
}

// loadFromEnv loads environment variables to initialize the DSN.
func (c *DBConfig) loadFromEnv() diabuddyErrors.ApiErrors {
	dbUrl := c.envManager.Get(envmanager.DbUrlKey)
	if dbUrl != "" {
		return c.parseDatabaseURL(dbUrl)
	}

	dbHost := c.envManager.Get(envmanager.DbHostKey)
	dbPort := c.envManager.Get(envmanager.DbPortKey, getDefaultPort(c.dbType))
	dbUsername := c.envManager.Get(envmanager.DbUsernameKey)
	dbPassword := c.envManager.Get(envmanager.DbPasswordKey)
	dbDatabase := c.envManager.Get(envmanager.DbDatabaseKey)

	dsnInstance, err := dsn.NewDSN(dbHost, dbPort, dbUsername, dbPassword, dbDatabase, c.params, dsnConnectionStringClosureWithDbType(c.dbType))
	if err != nil {
		return diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "cannot initialize DSN", diabuddyErrors.WithInternalError(err))
	}

	c.dsn = dsnInstance
	return nil
}

// parseDatabaseURL parses a database URL and initializes the DSN.
func (c *DBConfig) parseDatabaseURL(databaseUrl string) diabuddyErrors.ApiErrors {
	parsedURL, err := url.Parse(databaseUrl)
	if err != nil {
		return diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "invalid database URL", diabuddyErrors.WithInternalError(err))
	}

	dbUsername := parsedURL.User.Username()
	dbPassword, _ := parsedURL.User.Password()

	dbHost := parsedURL.Hostname()
	dbPort := parsedURL.Port()
	if dbPort == "" {
		dbPort = getDefaultPort(c.dbType)
	}

	dbDatabase := strings.TrimPrefix(parsedURL.Path, "/")

	// Merge existing parameters with those in the URL
	for key, value := range parsedURL.Query() {
		c.params[key] = value[0]
	}

	dsnInstance, dsnErr := dsn.NewDSN(dbHost, dbPort, dbUsername, dbPassword, dbDatabase, c.params, dsnConnectionStringClosureWithDbType(c.dbType))
	if dsnErr != nil {
		return diabuddyErrors.NewApiError(diabuddyErrors.InternalServerErrorType, "cannot create DSN", diabuddyErrors.WithInternalError(dsnErr))
	}

	c.dsn = dsnInstance
	return nil
}

// ConnectionString generates the DSN string.
func (c *DBConfig) ConnectionString() (string, diabuddyErrors.ApiErrors) {
	err := c.generateDSN()
	if err != nil {
		return "", err
	}
	return c.dsn.GenerateConnectionString(), nil
}

func (c *DBConfig) generateDSN() diabuddyErrors.ApiErrors {
	if c.dsn == nil {
		return c.loadFromEnv()
	}
	return nil
}

// dsnConnectionStringClosureWithDbType returns the appropriate DSN connection string function based on the dbType.
func dsnConnectionStringClosureWithDbType(dbType string) dsn.ConnectionStringOption {
	switch dbType {
	case Mysql:
		return dsn.WithMySqlConnectionString()
	case SqlServer:
		return dsn.WithSqlServerConnectionString()
	case Oracle:
		return dsn.WithOracleConnectionString()
	case MongoDb:
		return dsn.WithMongoDBConnectionString()
	case Redis:
		return dsn.WithRedisConnectionString()
	case Cassandra:
		return dsn.WithCassandraConnectionString()
	case Postgres:
		fallthrough
	default:
		return dsn.WithPostgresConnectionString()
	}
}

// getDefaultPort returns the default port for the specified dbType.
func getDefaultPort(dbType string) string {
	switch dbType {
	case Mysql:
		return MysqlDefaultPort
	case SqlServer:
		return SqlServerDefaultPort
	case Oracle:
		return OracleDefaultPort
	case MongoDb:
		return MongoDbDefaultPort
	case Redis:
		return RedisDefaultPort
	case Cassandra:
		return CassandraDefaultPort
	case Postgres:
		fallthrough
	default:
		return PostgresDefaultPort
	}
}

func getRequiredKeysForDBType(dbType string) []string {
	switch dbType {
	case Mysql, Postgres, SqlServer, Oracle, MongoDb, Cassandra:
		return []string{envmanager.DbHostKey, envmanager.DbUsernameKey, envmanager.DbPasswordKey, envmanager.DbDatabaseKey}
	case Redis:
		// Redis doesn't require a database name
		return []string{envmanager.DbHostKey, envmanager.DbUsernameKey, envmanager.DbPasswordKey}
	default:
		// Default keys if the dbType is not recognized
		return []string{envmanager.DbHostKey, envmanager.DbUsernameKey, envmanager.DbPasswordKey}
	}
}
