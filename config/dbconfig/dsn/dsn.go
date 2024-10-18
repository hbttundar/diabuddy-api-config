package dsn

import (
	"errors"
	"fmt"
	"strings"
)

const (
	SchemaPostgres  = "postgres"
	SchemaMySQL     = "mysql"
	SchemaSQLServer = "sqlserver"
	SchemaOracle    = "oracle"
	SchemaMongoDB   = "mongodb"
	SchemaRedis     = "redis"
	SchemaCassandra = "cassandra"
)

type DSN struct {
	Host                    string
	Port                    string
	Database                string
	Username                string
	Password                string
	Params                  map[string]string
	connectionStringClosure ConnectionStringOption
}

type ConnectionStringOption func(*DSN) string

// NewDSN initializes a new DSN with provided parameters and connection options.
// If no connection option is provided, it defaults to Postgres.
func NewDSN(host, port, username, password, database string, params map[string]string, options ...ConnectionStringOption) (*DSN, error) {
	// Enforce that only one option should be passed, otherwise return an error
	if len(options) > 1 {
		return nil, errors.New("only one connection configuration option can be provided, please create separate DSNs for different databases")
	}

	dsn := &DSN{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
		Params:   params,
	}

	// Apply provided option or default to Postgres connection string.
	if len(options) == 0 {
		dsn.connectionStringClosure = WithPostgresConnectionString() // Default to Postgres if no option is given
	} else {
		dsn.connectionStringClosure = options[0]
	}

	return dsn, nil
}

// WithPostgresConnectionString builds a connection string for PostgreSQL.
func WithPostgresConnectionString() ConnectionStringOption {
	return func(dsn *DSN) string {
		baseString := fmt.Sprintf("%s://%s:%s@%s:%s/%s", SchemaPostgres, dsn.Username, dsn.Password, dsn.Host, dsn.Port, dsn.Database)
		return extendsConnectionStringWithParams(baseString, dsn.Params)
	}
}

// WithMySqlConnectionString builds a connection string for MySQL.
func WithMySqlConnectionString() ConnectionStringOption {
	return func(dsn *DSN) string {
		baseString := fmt.Sprintf("%s://%s:%s@tcp(%s:%s)/%s", SchemaMySQL, dsn.Username, dsn.Password, dsn.Host, dsn.Port, dsn.Database)
		return extendsConnectionStringWithParams(baseString, dsn.Params)
	}
}

// WithSqlServerConnectionString builds a connection string for SQL Server.
func WithSqlServerConnectionString() ConnectionStringOption {
	return func(dsn *DSN) string {
		baseString := fmt.Sprintf("%s://%s:%s@%s:%s?database=%s", SchemaSQLServer, dsn.Username, dsn.Password, dsn.Host, dsn.Port, dsn.Database)
		return extendsConnectionStringWithParams(baseString, dsn.Params)
	}
}

// WithOracleConnectionString builds a connection string for Oracle.
func WithOracleConnectionString() ConnectionStringOption {
	return func(dsn *DSN) string {
		baseString := fmt.Sprintf("%s://%s:%s@%s:%s/%s", SchemaOracle, dsn.Username, dsn.Password, dsn.Host, dsn.Port, dsn.Database)
		return extendsConnectionStringWithParams(baseString, dsn.Params)
	}
}

// WithMongoDBConnectionString builds a connection string for MongoDB.
func WithMongoDBConnectionString() ConnectionStringOption {
	return func(dsn *DSN) string {
		baseString := fmt.Sprintf("%s://%s:%s@%s:%s/%s", SchemaMongoDB, dsn.Username, dsn.Password, dsn.Host, dsn.Port, dsn.Database)
		return extendsConnectionStringWithParams(baseString, dsn.Params)
	}
}

// WithRedisConnectionString builds a connection string for Redis.
func WithRedisConnectionString() ConnectionStringOption {
	return func(dsn *DSN) string {
		baseString := fmt.Sprintf("%s://%s:%s@%s:%s", SchemaRedis, dsn.Username, dsn.Password, dsn.Host, dsn.Port)
		return extendsConnectionStringWithParams(baseString, dsn.Params)
	}
}

// WithCassandraConnectionString builds a connection string for Cassandra.
func WithCassandraConnectionString() ConnectionStringOption {
	return func(dsn *DSN) string {
		baseString := fmt.Sprintf("%s://%s:%s@%s:%s/%s", SchemaCassandra, dsn.Username, dsn.Password, dsn.Host, dsn.Port, dsn.Database)
		return extendsConnectionStringWithParams(baseString, dsn.Params)
	}
}

func extendsConnectionStringWithParams(base string, params map[string]string) string {
	if params == nil || len(params) == 0 {
		return base
	}

	// Determine whether to use '?' or '&' to start appending parameters
	paramStr := ""
	if strings.Contains(base, "?") {
		paramStr = "&"
	} else {
		paramStr = "?"
	}

	// Construct the parameters in query format
	first := true
	for key, value := range params {
		if !first {
			paramStr += "&"
		}
		paramStr += fmt.Sprintf("%s=%s", key, value)
		first = false
	}

	return base + paramStr
}

func (dsn *DSN) GenerateConnectionString() string {
	return dsn.connectionStringClosure(dsn)
}
