package config

import "errors"

/*
DB_USER=admin
DB_PASSWORD=password
DB_NAME=db
*/
type Config struct {
	DbUser     string
	DbPassword string
	DbName     string
	DbHost     string
	TokenKey   string

	RedisPassword string
	RedisAddress  string
	RedisDb       string
}

func NewConfig(redisPassword string, redisAddress string, redisDb string, dbUser string, dbPassword string, dbName string, dbHost string, tokenKey string) (*Config, error) {
	//if we run locally, and do not include the service that runs our go backend into docker compose, we
	//need to get
	if dbUser == "" || dbPassword == "" || dbName == "" || dbHost == "" || tokenKey == "" || redisAddress == "" || redisDb == "" {
		return nil, errors.New("The environment is empty")
	}
	return &Config{RedisPassword: redisPassword, RedisAddress: redisAddress, RedisDb: redisDb, DbUser: dbUser, DbPassword: dbPassword, DbName: dbName, DbHost: dbHost, TokenKey: tokenKey}, nil
}
