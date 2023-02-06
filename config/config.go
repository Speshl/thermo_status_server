package config

import (
	"fmt"
	"os"
)

type Config struct {
	ConnString string
	HTTPPort   string
}

func NewConfig() *Config {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("TS_USER_NAME"),
		os.Getenv("TS_USER_PASS"),
		os.Getenv("TS_HOST"),
		os.Getenv("TS_PORT"),
		os.Getenv("TS_DB"),
	)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	return &Config{
		ConnString: connStr,
		HTTPPort:   httpPort,
	}
}
