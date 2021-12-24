package config

import (
	"flag"
	"time"
)

type Config struct {
	Port                 int
	GracefulShutdownTime time.Duration
	DBConfig
}

type DBConfig struct {
	DBHost   string
	DBPort   int
	DBName   string
	User     string
	Password string
}

func AddFlags(c *Config) {
	flag.IntVar(&c.Port, "port", 8080, "server port")
	flag.DurationVar(&c.GracefulShutdownTime, "gracefulShutdownTime", 30*time.Second, "graceful shutdown time")
	flag.StringVar(&c.DBHost, "dbHost", "localhost", "db host")
	flag.IntVar(&c.DBPort, "dbPort", 5432, "db port")
	flag.StringVar(&c.DBName, "dbName", "weave", "db name")
	flag.StringVar(&c.User, "user", "postgres", "db user")
	flag.StringVar(&c.Password, "password", "123456", "db password")
}

func New() *Config {
	return &Config{}
}
