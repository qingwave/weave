package config

import (
	"flag"
	"os"
	"weave/pkg/middleware/ratelimit"

	"gopkg.in/yaml.v2"
)

var appConfig = flag.String("config", "config/app.yaml", "application config path")

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
}

type ServerConfig struct {
	ENV                    string                  `yaml:"env"`
	Address                string                  `yaml:"address"`
	Port                   int                     `yaml:"port"`
	GracefulShutdownPeriod int                     `yaml:"gracefulShutdownPeriod"`
	LimitConfigs           []ratelimit.LimitConfig `yaml:"rateLimits"`
}

type DBConfig struct {
	DBHost   string `yaml:"dbHost"`
	DBPort   int    `yaml:"dbPort"`
	DBName   string `yaml:"dbName"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func Parse() (*Config, error) {
	config := &Config{}

	file, err := os.Open(*appConfig)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
