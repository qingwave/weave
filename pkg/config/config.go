package config

import (
	"os"
	"path"
	"path/filepath"
	"weave/pkg/utils/ratelimit"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server      ServerConfig           `yaml:"server"`
	DB          DBConfig               `yaml:"db"`
	Redis       RedisConfig            `yaml:"redis"`
	OAuthConfig map[string]OAuthConfig `yaml:"oauth"`
	AuthConfig  AuthenticationConfig   `yaml:"authentication"`
}

type ServerConfig struct {
	ENV                    string                  `yaml:"env"`
	Address                string                  `yaml:"address"`
	Port                   int                     `yaml:"port"`
	GracefulShutdownPeriod int                     `yaml:"gracefulShutdownPeriod"`
	LimitConfigs           []ratelimit.LimitConfig `yaml:"rateLimits"`
	JWTSecret              string                  `yaml:"jwtSecret"`
	DockerHost             string                  `yaml:"dockerHost"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Migrate  bool   `yaml:"migrate"`
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

type AuthenticationConfig struct {
	AuthModelConfigName             string `yaml:"authModelConfigName"`
	AuthModelConfigFullName         string
	AuthDefaultPolicyConfig         string `yaml:"authDefaultPolicyConfig"`
	LoadDefaultPolicyAlways         bool   `yaml:"loadDefaultPolicyAlways"`
	AuthDefaultPolicyConfigFullName string
	AuthTablePrefix                 string `yaml:"authTablePrefix"`
	AuthTableName                   string `yaml:"authTableName"`
}

func Parse(appConfig string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(appConfig)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	baseDir := filepath.Dir(appConfig)
	config.AuthConfig.AuthModelConfigFullName = path.Join(baseDir, config.AuthConfig.AuthModelConfigName)
	config.AuthConfig.AuthDefaultPolicyConfigFullName = path.Join(baseDir, config.AuthConfig.AuthDefaultPolicyConfig)

	return config, nil
}
