package dsn

import "fmt"

type DSN struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	SSLMode  string
}

func NewDsn(host, port, username, password, database, sslMode string) *DSN {
	return &DSN{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
		SSLMode:  sslMode,
	}
}

func (dsn *DSN) GenerateConnectionString(schema string) string {
	if dsn.SSLMode != "" {
		return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
			schema, dsn.Username, dsn.Password, dsn.Host, dsn.Port, dsn.Database, dsn.SSLMode)
	}
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		schema, dsn.Username, dsn.Password, dsn.Host, dsn.Port, dsn.Database)
}
