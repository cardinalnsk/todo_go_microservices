package config

import (
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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
	Salt              string        `yaml:"salt"`
	Expiration        time.Duration `yaml:"expiration"`
	JwtPrivateKeyPath string        `yaml:"jwt_private_key_path"`
	PrivateKey        *rsa.PrivateKey
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
		logrus.Warning("No .env file found or failed to load. Using default config.")
		logrus.Warning("ERROR: ", err.Error())
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

	cfg.Auth.PrivateKey = LoadPrivateKey(cfg.Auth.JwtPrivateKeyPath)
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

func LoadPrivateKey(path string) *rsa.PrivateKey {
	keyData, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		panic(err)
	}
	return privKey
}
