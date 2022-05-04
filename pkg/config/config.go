package config

import (
	"flag"
	"os"
	"weave/pkg/middleware/ratelimit"

	"gopkg.in/yaml.v2"
)

var appConfig = flag.String("config", "config/app.yaml", "application config path")

type Config struct {
	Server      ServerConfig           `yaml:"server"`
	DB          DBConfig               `yaml:"db"`
	Redis       RedisConfig            `yaml:"redis"`
	OAuthConfig map[string]OAuthConfig `yaml:"oauth"`
}

type ServerConfig struct {
	ENV                    string                  `yaml:"env"`
	Address                string                  `yaml:"address"`
	Port                   int                     `yaml:"port"`
	GracefulShutdownPeriod int                     `yaml:"gracefulShutdownPeriod"`
	LimitConfigs           []ratelimit.LimitConfig `yaml:"rateLimits"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type RedisConfig struct {
	Enable   bool   `yaml:"enable"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type OAuthConfig struct {
	AuthType     string `yaml:"authType"`
	ClientId     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
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
