package config

import (
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"regexp"
)

type ServiceConfig struct {
	Port string `yaml:"port"`
	Name string `yaml:"name"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Pretty bool   `yaml:"pretty"`
}

type GatewayConfig struct {
	JwtKeyPath string `yaml:"jwt_public_key_path"`
	PublicKey  *rsa.PublicKey
}

type UpstreamServiceConfig struct {
	Url string `yaml:"url"`
}

type Config struct {
	Main    ServiceConfig         `yaml:"main"`
	Logger  LoggerConfig          `yaml:"logger"`
	Gateway GatewayConfig         `yaml:"gateway"`
	Auth    UpstreamServiceConfig `yaml:"auth"`
	Todo    UpstreamServiceConfig `yaml:"todo"`
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
	cfg.Gateway.PublicKey = LoadPublicKey(cfg.Gateway.JwtKeyPath)
	return &cfg, nil
}

// replaceEnvVars поддержка ${ENV_VAR:default} — аналогично вашему стилю
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

func LoadPublicKey(path string) *rsa.PublicKey {
	pubKeyData, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		panic(err)
	}
	return pubKey
}
