package dbconfig_test

import (
	dbconfig "github.com/hbttundar/diabuddy-api-config/config/dbconfig"
	"github.com/hbttundar/diabuddy-api-config/config/envmanager"
	testmain "github.com/hbttundar/diabuddy-api-config/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDBConfig_ConnectionString(t *testing.T) {
	tests := []struct {
		name            string
		dbType          string
		setupEnv        map[string]string
		params          map[string]string
		expectedConnStr string
		expectError     bool
	}{
		// PostgresSQL Tests
		{
			name:            "Postgres with default port without Database URL",
			dbType:          "postgres",
			setupEnv:        map[string]string{"DATABASE_URL": "", "DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_PORT": "", "DB_HOST": "localhost", "DB_DATABASE": "testdb"},
			expectedConnStr: "postgres://user:password@localhost:5432/testdb",
			expectError:     false,
		},
		{
			name:            "Postgres with Database url with no port",
			dbType:          "postgres",
			setupEnv:        map[string]string{"DATABASE_URL": "postgres://user:password@localhost/testdb"},
			expectedConnStr: "postgres://user:password@localhost:5432/testdb",
			expectError:     false,
		},
		{
			name:            "Postgres with custom port",
			dbType:          "postgres",
			setupEnv:        map[string]string{"DATABASE_URL": "", "DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_PORT": "6543", "DB_HOST": "localhost", "DB_DATABASE": "testdb"},
			expectedConnStr: "postgres://user:password@localhost:6543/testdb",
			expectError:     false,
		},
		{
			name:            "Postgres with multiple params",
			dbType:          "postgres",
			setupEnv:        map[string]string{"DATABASE_URL": "", "DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_PORT": "5432", "DB_HOST": "localhost", "DB_DATABASE": "testdb"},
			params:          map[string]string{"sslmode": "disable", "timezone": "UTC"},
			expectedConnStr: "postgres://user:password@localhost:5432/testdb?sslmode=disable&timezone=UTC",
			expectError:     false,
		},

		// MySQL Tests
		{
			name:            "MySQL with default port without database url",
			dbType:          "mysql",
			setupEnv:        map[string]string{"DATABASE_URL": "", "DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_PORT": "", "DB_HOST": "localhost", "DB_DATABASE": "testdb"},
			expectedConnStr: "mysql://user:password@tcp(localhost:3306)/testdb",
			expectError:     false,
		},
		{
			name:            "MySQL with database url",
			dbType:          "mysql",
			setupEnv:        map[string]string{"DATABASE_URL": "mysql://user:password@localhost:3306/testdb"},
			expectedConnStr: "mysql://user:password@tcp(localhost:3306)/testdb",
			expectError:     false,
		},
		{
			name:            "MySQL with database url with no port",
			dbType:          "mysql",
			setupEnv:        map[string]string{"DATABASE_URL": "mysql://user:password@localhost/testdb"},
			expectedConnStr: "mysql://user:password@tcp(localhost:3306)/testdb",
			expectError:     false,
		},
		{
			name:            "MySQL with custom port",
			dbType:          "mysql",
			setupEnv:        map[string]string{"DATABASE_URL": "", "DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_PORT": "3307", "DB_HOST": "localhost", "DB_DATABASE": "testdb"},
			expectedConnStr: "mysql://user:password@tcp(localhost:3307)/testdb",
			expectError:     false,
		},

		// Oracle Tests
		{
			name:            "Oracle with default port without database url",
			dbType:          "oracle",
			setupEnv:        map[string]string{"DATABASE_URL": "", "DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_PORT": "1521", "DB_HOST": "localhost", "DB_DATABASE": "testdb"},
			expectedConnStr: "oracle://user:password@localhost:1521/testdb",
			expectError:     false,
		},
		{
			name:            "Oracle with database url",
			dbType:          "oracle",
			setupEnv:        map[string]string{"DATABASE_URL": "mysql://user:password@localhost:1521/testdb"},
			expectedConnStr: "oracle://user:password@localhost:1521/testdb",
			expectError:     false,
		},
		{
			name:            "Oracle with database url with no port",
			dbType:          "oracle",
			setupEnv:        map[string]string{"DATABASE_URL": "mysql://user:password@localhost/testdb"},
			expectedConnStr: "oracle://user:password@localhost:1521/testdb",
			expectError:     false,
		},
		{
			name:            "Oracle with custom port",
			dbType:          "oracle",
			setupEnv:        map[string]string{"DATABASE_URL": "", "DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_PORT": "1526", "DB_HOST": "localhost", "DB_DATABASE": "testdb"},
			expectedConnStr: "oracle://user:password@localhost:1526/testdb",
			expectError:     false,
		},

		// SQL Server Tests
		{
			name:            "SQL Server with default port",
			dbType:          "sqlserver",
			setupEnv:        map[string]string{"DATABASE_URL": "", "DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_PORT": "1433", "DB_HOST": "localhost", "DB_DATABASE": "testdb"},
			expectedConnStr: "sqlserver://user:password@localhost:1433?database=testdb",
			expectError:     false,
		},
		{
			name:            "SQL Server with database url ",
			dbType:          "sqlserver",
			setupEnv:        map[string]string{"DATABASE_URL": "sqlserver://user:password@localhost:1433/testdb"},
			expectedConnStr: "sqlserver://user:password@localhost:1433?database=testdb",
			expectError:     false,
		},
		{
			name:            "SQL Server with database url  with no port ",
			dbType:          "sqlserver",
			setupEnv:        map[string]string{"DATABASE_URL": "sqlserver://user:password@localhost/testdb"},
			expectedConnStr: "sqlserver://user:password@localhost:1433?database=testdb",
			expectError:     false,
		},
		{
			name:            "SQL Server with multiple params",
			dbType:          "sqlserver",
			setupEnv:        map[string]string{"DATABASE_URL": "", "DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_PORT": "1433", "DB_HOST": "localhost", "DB_DATABASE": "testdb"},
			params:          map[string]string{"encrypt": "disable", "applicationIntent": "ReadOnly"},
			expectedConnStr: "sqlserver://user:password@localhost:1433?database=testdb&encrypt=disable&applicationIntent=ReadOnly",
			expectError:     false,
		},
		// MongoDB Tests
		{
			name:            "MongoDB with default port",
			dbType:          "mongodb",
			setupEnv:        map[string]string{"DATABASE_URL": "", "DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_PORT": "", "DB_HOST": "localhost", "DB_DATABASE": "testdb"},
			expectedConnStr: "mongodb://user:password@localhost:27017/testdb",
			expectError:     false,
		},
		{
			name:            "MongoDB with database url",
			dbType:          "mongodb",
			setupEnv:        map[string]string{"DATABASE_URL": "mongodb://user:password@localhost:27017/testdb"},
			expectedConnStr: "mongodb://user:password@localhost:27017/testdb",
			expectError:     false,
		},
		{
			name:            "MongoDB with database url with no port",
			dbType:          "mongodb",
			setupEnv:        map[string]string{"DATABASE_URL": "mongodb://user:password@localhost/testdb"},
			expectedConnStr: "mongodb://user:password@localhost:27017/testdb",
			expectError:     false,
		},
		{
			name:            "MongoDB with custom port",
			dbType:          "mongodb",
			setupEnv:        map[string]string{"DATABASE_URL": "", "DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_PORT": "28017", "DB_HOST": "localhost", "DB_DATABASE": "testdb"},
			expectedConnStr: "mongodb://user:password@localhost:28017/testdb",
			expectError:     false,
		},
		// Redis Tests
		{
			name:            "Redis with default port",
			dbType:          "redis",
			setupEnv:        map[string]string{"DATABASE_URL": "", "DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_PORT": "", "DB_HOST": "localhost", "DB_DATABASE": ""},
			expectedConnStr: "redis://user:password@localhost:6379",
			expectError:     false,
		},
		// Cassandra Tests
		{
			name:            "Cassandra with default port",
			dbType:          "cassandra",
			setupEnv:        map[string]string{"DATABASE_URL": "", "DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_PORT": "", "DB_HOST": "localhost", "DB_DATABASE": "testdb"},
			expectedConnStr: "cassandra://user:password@localhost:9042/testdb",
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment for each test
			setupTestEnvironment(tt.setupEnv)
			var dbConfig = &dbconfig.DBConfig{}
			envManager, err := envmanager.NewEnvManager(envmanager.WithUseDefault(false))
			assert.NoError(t, err, "expected no error while creating EnvManager")
			err = envManager.LoadEnvironmentVariables()
			assert.NoError(t, err, "expected no error while loading env variables EnvManager")
			if tt.params != nil {
				dbConfig, err = dbconfig.NewDBConfig(envManager, dbconfig.WithType(tt.dbType), dbconfig.WithDsnParameters(tt.params))
				if tt.expectError {
					assert.Error(t, err, "expected error but got nil")
					return
				}
			} else {
				dbConfig, err = dbconfig.NewDBConfig(envManager, dbconfig.WithType(tt.dbType))
				if tt.expectError {
					assert.Error(t, err, "expected error but got nil")
					return
				}
			}
			assert.NoError(t, err, "expected no error while creating DBConfig")

			connStr, err := dbConfig.ConnectionString()
			assert.NoError(t, err, "expected no error while generating connection string")
			assert.Equal(t, tt.expectedConnStr, connStr, "expected connection string to match")
		})
	}
}

func TestDBConfig_MultipleOptionsShouldReturnError(t *testing.T) {
	// Setup for test
	testmain.Setup()
	defer testmain.TearDown()

	envManager, err := envmanager.NewEnvManager()
	assert.NoError(t, err, "expected no error while creating EnvManager")

	t.Run("Multiple connection options should return error", func(t *testing.T) {
		_, err := dbconfig.NewDBConfig(envManager, dbconfig.WithType(dbconfig.Postgres), dbconfig.WithType(dbconfig.Mysql))
		assert.Error(t, err, "expected error when providing multiple connection options")
	})

	t.Run("Multiple DSNParams options should return error", func(t *testing.T) {
		_, err := dbconfig.NewDBConfig(envManager, dbconfig.WithDsnParameters(map[string]string{"timezone": "UTC"}), dbconfig.WithDsnParameters(map[string]string{"sslmode": "disable"}))
		assert.Error(t, err, "expected error when providing multiple connection options")
	})
}

func TestDBConfig_Validate(t *testing.T) {
	tests := []struct {
		name            string
		dbType          string
		setupEnv        map[string]string
		expectedToError bool
	}{
		// Postgres
		{
			name:            "Postgres validation passes with all required keys set",
			dbType:          dbconfig.Postgres,
			setupEnv:        map[string]string{"DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_DATABASE": "testdb", "DB_HOST": "localhost"},
			expectedToError: false,
		},
		{
			name:            "Postgres validation fails with missing keys",
			dbType:          dbconfig.Postgres,
			setupEnv:        map[string]string{"DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_DATABASE": ""},
			expectedToError: true,
		},
		// Mysql
		{
			name:            "Mysql validation passes with all required keys set",
			dbType:          dbconfig.Mysql,
			setupEnv:        map[string]string{"DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_DATABASE": "testdb", "DB_HOST": "localhost"},
			expectedToError: false,
		},
		{
			name:            "Mysql validation fails with missing keys",
			dbType:          dbconfig.Mysql,
			setupEnv:        map[string]string{"DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_DATABASE": ""},
			expectedToError: true,
		},
		// Sqlserver
		{
			name:            "Sqlserver validation passes with all required keys set",
			dbType:          dbconfig.SqlServer,
			setupEnv:        map[string]string{"DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_DATABASE": "testdb", "DB_HOST": "localhost"},
			expectedToError: false,
		},
		{
			name:            "Sqlserver validation fails with missing keys",
			dbType:          dbconfig.SqlServer,
			setupEnv:        map[string]string{"DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_DATABASE": ""},
			expectedToError: true,
		},
		// Mongodb
		{
			name:            "MongoDb validation passes with all required keys set",
			dbType:          dbconfig.MongoDb,
			setupEnv:        map[string]string{"DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_DATABASE": "testdb", "DB_HOST": "localhost"},
			expectedToError: false,
		},
		{
			name:            "MongoDb validation fails with missing keys",
			dbType:          dbconfig.MongoDb,
			setupEnv:        map[string]string{"DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_DATABASE": ""},
			expectedToError: true,
		},
		// Oracle
		{
			name:            "Oracle validation passes with all required keys set",
			dbType:          dbconfig.Oracle,
			setupEnv:        map[string]string{"DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_DATABASE": "testdb", "DB_HOST": "localhost"},
			expectedToError: false,
		},
		{
			name:            "Oracle validation fails with missing keys",
			dbType:          dbconfig.Oracle,
			setupEnv:        map[string]string{"DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_DATABASE": ""},
			expectedToError: true,
		},
		// Cassandra
		{
			name:            "Cassandra validation passes with all required keys set",
			dbType:          dbconfig.Cassandra,
			setupEnv:        map[string]string{"DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_DATABASE": "testdb", "DB_HOST": "localhost"},
			expectedToError: false,
		},
		{
			name:            "Cassandra validation fails with missing keys",
			dbType:          dbconfig.Cassandra,
			setupEnv:        map[string]string{"DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_DATABASE": ""},
			expectedToError: true,
		},

		//redis
		{
			name:            "Redis validation passes without database key",
			dbType:          "redis",
			setupEnv:        map[string]string{"DB_USERNAME": "user", "DB_PASSWORD": "password", "DB_HOST": "localhost"},
			expectedToError: false,
		},
		{
			name:            "Redis validation fails with missing username",
			dbType:          "redis",
			setupEnv:        map[string]string{"DB_PASSWORD": "password", "DB_HOST": "localhost", "DB_USERNAME": ""},
			expectedToError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestEnvironment(tt.setupEnv)
			envManager, _ := envmanager.NewEnvManager(envmanager.WithUseDefault(false))

			dbConfig, _ := dbconfig.NewDBConfig(envManager, dbconfig.WithType(tt.dbType))
			err := dbConfig.Validate()

			if tt.expectedToError {
				assert.Error(t, err, "expected validation to fail")
			} else {
				assert.NoError(t, err, "expected validation to pass")
			}
		})
	}
}

func TestDBConfig_Get(t *testing.T) {
	tests := []struct {
		name                 string
		setupEnv             map[string]string
		useDefault           bool
		key                  string
		defaultValue         string
		expectedValue        string
		expectedToUseDefault bool
	}{
		{
			name:          "Value exists in environment variables",
			setupEnv:      map[string]string{"DB_USERNAME": "env_user"},
			useDefault:    false,
			key:           "DB_USERNAME",
			defaultValue:  "default_user",
			expectedValue: "env_user",
		},
		{
			name:                 "Value does not exist, use default value when configured to use default",
			setupEnv:             map[string]string{"DB_USERNAME": ""},
			useDefault:           true,
			key:                  "DB_USERNAME",
			defaultValue:         "default_user",
			expectedValue:        "default_user",
			expectedToUseDefault: true,
		},
		{
			name:                 "Value does not exist, no default value available, configured not to use default",
			setupEnv:             map[string]string{"DB_USERNAME": ""},
			useDefault:           false,
			key:                  "DB_USERNAME",
			defaultValue:         "",
			expectedValue:        "",
			expectedToUseDefault: false,
		},
		{
			name:                 "Value does not exist, default value should be returned",
			setupEnv:             map[string]string{},
			useDefault:           true,
			key:                  "DB_PASSWORD",
			defaultValue:         "default_pass",
			expectedValue:        "root",
			expectedToUseDefault: true,
		},
		{
			name:                 "Empty environment value, configured not to use default value, return empty",
			setupEnv:             map[string]string{"DB_DATABASE": ""},
			useDefault:           false,
			key:                  "DB_DATABASE",
			defaultValue:         "",
			expectedValue:        "",
			expectedToUseDefault: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestEnvironment(tt.setupEnv)
			envManager, _ := envmanager.NewEnvManager(envmanager.WithUseDefault(tt.useDefault))
			dbConfig, _ := dbconfig.NewDBConfig(envManager)
			var actualValue string
			if tt.defaultValue != "" {
				actualValue = dbConfig.Get(tt.key, tt.defaultValue)
			}
			actualValue = dbConfig.Get(tt.key)

			assert.Equal(t, tt.expectedValue, actualValue, "expected value to match")
		})
	}
}

func setupTestEnvironment(setupEnv map[string]string) {
	// Clear and set up the environment based on the provided map
	testmain.TearDown()
	for key, value := range setupEnv {
		testmain.EnvVars[key] = value
	}
	testmain.Setup()
}
