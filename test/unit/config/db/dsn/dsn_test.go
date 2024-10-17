package dsn_test

import (
	"github.com/hbttundar/diabuddy-api-config/config/db/dsn"
	"testing"
)

func TestDSN_GenerateConnectionString_WithSSLMode(t *testing.T) {
	databaseDsn := dsn.DSN{
		Host:     "localhost",
		Port:     "5432",
		Database: "testdb",
		Username: "user",
		Password: "password",
		SSLMode:  "require",
	}

	expected := "host=localhost port=5432 user=user password=password dbname=testdb sslmode=require"
	actual := databaseDsn.GenerateConnectionString()

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

	expected := "host=localhost port=5432 user=user password=password dbname=testdb"
	actual := databaseDsn.GenerateConnectionString()

	if actual != expected {
		t.Errorf("expected database string '%s', but got '%s'", expected, actual)
	}
}
