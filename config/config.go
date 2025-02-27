package config

import (
	"errors"
	"os"
	"strconv"
)

var errNoEnvFound = errors.New("no enviroment variable found")

type Config struct {
	ServicePort        string
	Database           DB
	VerseString        string
	EnrichmentAddr     string
	EnrichmentEndpoint string
	LogLevel           string
	PageSize           int
}

type DB struct {
	DBName string
	DBPort string
	DBHost string
	DBPwd  string
	DBUser string
}

func NewConfig() (*Config, error) {
	ServicePort, exists := os.LookupEnv("SERVICE_PORT")
	if !exists {
		return nil, errNoEnvFound
	}
	DBName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		return nil, errNoEnvFound
	}
	DBPort, exists := os.LookupEnv("DB_PORT")
	if !exists {
		return nil, errNoEnvFound
	}
	DBUser, exists := os.LookupEnv("DB_USER")
	if !exists {
		return nil, errNoEnvFound
	}
	DBPwd, exists := os.LookupEnv("DB_PWD")
	if !exists {
		return nil, errNoEnvFound
	}
	DBHost, exists := os.LookupEnv("DB_HOST")
	if !exists {
		return nil, errNoEnvFound
	}
	DB := DB{DBName: DBName, DBPort: DBPort, DBHost: DBHost, DBUser: DBUser, DBPwd: DBPwd}
	EnrichmentAddr, exists := os.LookupEnv("ENRICHMENT_ADDR")
	if !exists {
		return nil, errNoEnvFound
	}
	EnrichmentEndpoint, exists := os.LookupEnv("ENRICHMENT_ENDPOINT")
	if !exists {
		return nil, errNoEnvFound
	}
	VerseString, exists := os.LookupEnv("VERSE")
	if !exists {
		return nil, errNoEnvFound
	}
	LogLevel, exists := os.LookupEnv("LOG_LEVEL")
	if !exists {
		return nil, errNoEnvFound
	}
	PageSizeStr, exists := os.LookupEnv("PAGE_SIZE")
	if !exists {
		return nil, errNoEnvFound
	}
	PageSize, err := strconv.Atoi(PageSizeStr)
	if err != nil {
		return nil, err
	}
	return &Config{ServicePort: ServicePort, Database: DB, EnrichmentAddr: EnrichmentAddr, EnrichmentEndpoint: EnrichmentEndpoint, VerseString: VerseString, LogLevel: LogLevel, PageSize: PageSize}, nil
}
