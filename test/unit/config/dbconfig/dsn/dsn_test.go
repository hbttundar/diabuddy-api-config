package dsn_test

import (
	"github.com/hbttundar/diabuddy-api-config/config/dbconfig/dsn"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDSN_GenerateConnectionString(t *testing.T) {
	tests := []struct {
		name           string
		host           string
		port           string
		username       string
		password       string
		database       string
		params         map[string]string
		options        []dsn.ConnectionStringOption
		expectedString string
		expectError    bool
	}{
		// PostgreSQL with Params
		{
			name:           "Postgres with Params",
			host:           "localhost",
			port:           "5432",
			username:       "user",
			password:       "password",
			database:       "testdb",
			params:         map[string]string{"sslmode": "disable"},
			options:        []dsn.ConnectionStringOption{dsn.WithPostgresConnectionString()},
			expectedString: "postgres://user:password@localhost:5432/testdb?sslmode=disable",
			expectError:    false,
		},
		// PostgreSQL without Params
		{
			name:           "Postgres without Params",
			host:           "localhost",
			port:           "5432",
			username:       "user",
			password:       "password",
			database:       "testdb",
			params:         nil,
			options:        []dsn.ConnectionStringOption{dsn.WithPostgresConnectionString()},
			expectedString: "postgres://user:password@localhost:5432/testdb",
			expectError:    false,
		},
		// MySQL with Params
		{
			name:           "MySQL with Params",
			host:           "localhost",
			port:           "3306",
			username:       "user",
			password:       "password",
			database:       "testdb",
			params:         map[string]string{"charset": "utf8mb4"},
			options:        []dsn.ConnectionStringOption{dsn.WithMySqlConnectionString()},
			expectedString: "mysql://user:password@tcp(localhost:3306)/testdb?charset=utf8mb4",
			expectError:    false,
		},
		// MySQL without Params
		{
			name:           "MySQL without Params",
			host:           "localhost",
			port:           "3306",
			username:       "user",
			password:       "password",
			database:       "testdb",
			params:         nil,
			options:        []dsn.ConnectionStringOption{dsn.WithMySqlConnectionString()},
			expectedString: "mysql://user:password@tcp(localhost:3306)/testdb",
			expectError:    false,
		},
		// SQL Server with Params
		{
			name:           "SQL Server with Params",
			host:           "localhost",
			port:           "1433",
			username:       "user",
			password:       "password",
			database:       "testdb",
			params:         map[string]string{"encrypt": "disable"},
			options:        []dsn.ConnectionStringOption{dsn.WithSqlServerConnectionString()},
			expectedString: "sqlserver://user:password@localhost:1433?database=testdb&encrypt=disable",
			expectError:    false,
		},
		// SQL Server without Params
		{
			name:           "SQL Server without Params",
			host:           "localhost",
			port:           "1433",
			username:       "user",
			password:       "password",
			database:       "testdb",
			params:         nil,
			options:        []dsn.ConnectionStringOption{dsn.WithSqlServerConnectionString()},
			expectedString: "sqlserver://user:password@localhost:1433?database=testdb",
			expectError:    false,
		},
		// Oracle with Params
		{
			name:           "Oracle with Params",
			host:           "localhost",
			port:           "1521",
			username:       "user",
			password:       "password",
			database:       "oracledb",
			params:         map[string]string{"TRACE_LEVEL_CLIENT": "16"},
			options:        []dsn.ConnectionStringOption{dsn.WithOracleConnectionString()},
			expectedString: "oracle://user:password@localhost:1521/oracledb?TRACE_LEVEL_CLIENT=16",
			expectError:    false,
		},
		// Oracle without Params
		{
			name:           "Oracle without Params",
			host:           "localhost",
			port:           "1521",
			username:       "user",
			password:       "password",
			database:       "oracledb",
			params:         nil,
			options:        []dsn.ConnectionStringOption{dsn.WithOracleConnectionString()},
			expectedString: "oracle://user:password@localhost:1521/oracledb",
			expectError:    false,
		},
		// MongoDB with Params
		{
			name:           "MongoDB with Params",
			host:           "localhost",
			port:           "27017",
			username:       "user",
			password:       "password",
			database:       "mongodb",
			params:         map[string]string{"authSource": "admin"},
			options:        []dsn.ConnectionStringOption{dsn.WithMongoDBConnectionString()},
			expectedString: "mongodb://user:password@localhost:27017/mongodb?authSource=admin",
			expectError:    false,
		},
		// MongoDB without Params
		{
			name:           "MongoDB without Params",
			host:           "localhost",
			port:           "27017",
			username:       "user",
			password:       "password",
			database:       "mongodb",
			params:         nil,
			options:        []dsn.ConnectionStringOption{dsn.WithMongoDBConnectionString()},
			expectedString: "mongodb://user:password@localhost:27017/mongodb",
			expectError:    false,
		},
		// Cassandra with Params
		{
			name:           "Cassandra with Params",
			host:           "localhost",
			port:           "9042",
			username:       "user",
			password:       "password",
			database:       "cassandra",
			params:         map[string]string{"consistency": "QUORUM"},
			options:        []dsn.ConnectionStringOption{dsn.WithCassandraConnectionString()},
			expectedString: "cassandra://user:password@localhost:9042/cassandra?consistency=QUORUM",
			expectError:    false,
		},
		// Cassandra without Params
		{
			name:           "Cassandra without Params",
			host:           "localhost",
			port:           "9042",
			username:       "user",
			password:       "password",
			database:       "cassandra",
			params:         nil,
			options:        []dsn.ConnectionStringOption{dsn.WithCassandraConnectionString()},
			expectedString: "cassandra://user:password@localhost:9042/cassandra",
			expectError:    false,
		},
		// Redis with Params
		{
			name:           "Redis with Params",
			host:           "localhost",
			port:           "6379",
			username:       "user",
			password:       "password",
			database:       "",
			params:         map[string]string{"timeout": "10s"},
			options:        []dsn.ConnectionStringOption{dsn.WithRedisConnectionString()},
			expectedString: "redis://user:password@localhost:6379?timeout=10s",
			expectError:    false,
		},
		// Redis without Params
		{
			name:           "Redis without Params",
			host:           "localhost",
			port:           "6379",
			username:       "user",
			password:       "password",
			database:       "",
			params:         nil,
			options:        []dsn.ConnectionStringOption{dsn.WithRedisConnectionString()},
			expectedString: "redis://user:password@localhost:6379",
			expectError:    false,
		},
		// Default with Params (defaults to Postgres)
		{
			name:           "Default without Options (defaults to Postgres)",
			host:           "localhost",
			port:           "5432",
			username:       "user",
			password:       "password",
			database:       "testdb",
			params:         map[string]string{"sslmode": "disable"},
			options:        nil,
			expectedString: "postgres://user:password@localhost:5432/testdb?sslmode=disable",
			expectError:    false,
		},
		// SQL Server with multiple params
		{
			name:           "SQL Server with Multiple Params",
			host:           "localhost",
			port:           "1433",
			username:       "user",
			password:       "password",
			database:       "testdb",
			params:         map[string]string{"encrypt": "disable", "applicationIntent": "ReadOnly"},
			options:        []dsn.ConnectionStringOption{dsn.WithSqlServerConnectionString()},
			expectedString: "sqlserver://user:password@localhost:1433?database=testdb&encrypt=disable&applicationIntent=ReadOnly",
			expectError:    false,
		},
		// PostgreSQL with multiple params
		{
			name:           "Postgres with Multiple Params",
			host:           "localhost",
			port:           "5432",
			username:       "user",
			password:       "password",
			database:       "testdb",
			params:         map[string]string{"sslmode": "disable", "timezone": "UTC"},
			options:        []dsn.ConnectionStringOption{dsn.WithPostgresConnectionString()},
			expectedString: "postgres://user:password@localhost:5432/testdb?sslmode=disable&timezone=UTC",
			expectError:    false,
		},
		// Test: Passing multiple options should return an error
		{
			name:           "Error when Multiple Options Provided",
			host:           "localhost",
			port:           "5432",
			username:       "user",
			password:       "password",
			database:       "testdb",
			params:         nil,
			options:        []dsn.ConnectionStringOption{dsn.WithPostgresConnectionString(), dsn.WithMySqlConnectionString()},
			expectedString: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create DSN with provided options
			dsnInstance, err := dsn.NewDSN(tt.host, tt.port, tt.username, tt.password, tt.database, tt.params, tt.options...)

			if tt.expectError {
				assert.Error(t, err, "expected error but got nil")
				return
			}

			assert.NoError(t, err, "expected no error but got one")
			actualConnectionString := dsnInstance.GenerateConnectionString()
			assert.Equal(t, tt.expectedString, actualConnectionString, "connection string did not match expected")
		})
	}
}
