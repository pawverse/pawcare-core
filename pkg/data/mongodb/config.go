package mongodb

import "fmt"

func ConnectionString(username string, password string, host string, port uint32) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
}

type DBConfig interface {
	GetUsername() string
	GetPassword() string
	GetHost() string
	GetPort() uint32
}

func ConnectionStringConfig(config DBConfig) string {
	return ConnectionString(config.GetUsername(), config.GetPassword(), config.GetHost(), config.GetPort())
}
