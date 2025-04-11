package numeral

import (
	"os"
	"strconv"
)

// TODO something
type Config struct {
	Port               int64
	DbLocation         string
	BankFolderLocation string
}

type ConfigOption func(*Config)

func NewConfigFromEnv() *Config {
	port := getEnvInt("LISTEN_PORT", 8080)
	bankFolderLocation := getEnvString("BANK_FOLDER", "data/bank/")
	dbLocation := getEnvString("SQLITE_DB_FILE_LOCATION", "data/data.sqlite")
	return &Config{
		Port:               port,
		DbLocation:         dbLocation,
		BankFolderLocation: bankFolderLocation,
	}
}

func getEnvString(name, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		return defaultValue
	}
	return v
}

func getEnvInt(name string, defaultValue int64) int64 {
	port := os.Getenv(name)
	v, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		return defaultValue
	}
	return v
}

func NewDefaultConfig(opts ...ConfigOption) *Config {
	c := &Config{
		Port:               8080,
		DbLocation:         "data.sqlite",
		BankFolderLocation: "results",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
