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

func (dsn *DSN) GenerateConnectionString() string {
	if dsn.SSLMode != "" {
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dsn.Host, dsn.Port, dsn.Username, dsn.Password, dsn.Database, dsn.SSLMode)
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		dsn.Host, dsn.Port, dsn.Username, dsn.Password, dsn.Database)
}
