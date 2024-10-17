package dsn_test

import (
	"github.com/hbttundar/diabuddy-api-config/config/db/dsn"
	"testing"
)

const schema = "postgres"

func TestDSN_GenerateConnectionString_WithSSLMode(t *testing.T) {
	databaseDsn := dsn.DSN{
		Host:     "localhost",
		Port:     "5432",
		Database: "testdb",
		Username: "user",
		Password: "password",
		SSLMode:  "require",
	}

	expected := "postgres://user:password@localhost:5432/testdb?sslmode=require"
	actual := databaseDsn.GenerateConnectionString(schema)

	if actual != expected {
		t.Errorf("expected database string '%s', but got '%s'", expected, actual)
	}
}

func TestDSN_GenerateConnectionString_WithoutSSLMode(t *testing.T) {
	databaseDsn := dsn.DSN{
		Host:     "localhost",
		Port:     "5432",
		Database: "testdb",
		Username: "user",
		Password: "password",
	}

	expected := "postgres://user:password@localhost:5432/testdb"
	actual := databaseDsn.GenerateConnectionString(schema)

	if actual != expected {
		t.Errorf("expected database string '%s', but got '%s'", expected, actual)
	}
}
