package config

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"regexp"
	"time"
)

type ServiceConfig struct {
	Port string `yaml:"port"`
	Name string `yaml:"name"`
}

type DataSourceConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Pretty bool   `yaml:"pretty"`
}

type AuthConfig struct {
	Salt       string        `yaml:"salt"`
	SigningKey string        `yaml:"signing_key"`
	Expiration time.Duration `yaml:"expiration"`
}

type Config struct {
	DataSource DataSourceConfig `yaml:"datasource"`
	Main       ServiceConfig    `yaml:"main"`
	Logger     LoggerConfig     `yaml:"logger"`
	Auth       AuthConfig       `yaml:"auth"`
}

func LoadConfig(path string) (*Config, error) {
	gin.SetMode(gin.ReleaseMode)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	processed := replaceEnvVars(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(processed), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// replaceEnvVars Поддержка ${ENV_VAR:default} AKA SpringBoot
func replaceEnvVars(input string) string {
	re := regexp.MustCompile(`\$\{([^}:\s]+)(?::([^}]*))?\}`)
	return re.ReplaceAllStringFunc(input, func(s string) string {
		matches := re.FindStringSubmatch(s)
		envVar := matches[1]
		defVal := matches[2]
		val := os.Getenv(envVar)
		if val == "" {
			val = defVal
		}
		return val
	})
}
