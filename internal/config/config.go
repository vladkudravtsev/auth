package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `env:"ENV" env-default:"local"`
	Database   Database
	HTTPServer HTTPServer
	GRPCServer GRPCServer
	Auth       Auth
	AppSecret  string `env:"APP_SECRET" env-required:"true"`
}

type Database struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	DBName   string `env:"DB_NAME" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `env:"HTTP_ADDRESS" env-default:"localhost:8080"`
	Timeout     time.Duration `env:"HTTP_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

type GRPCServer struct {
	Port int `env:"GRPC_PORT" env-default:"8081"`
}

type Auth struct {
	TokenTTL time.Duration `env:"AUTH_TOKEN_TTL" env-default:"1m"`
}

func MustLoad() *Config {
	var cfg Config

	path := fetchConfigPath()

	if path == "" {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("envs not loaded: %s", err)
		}
		return &cfg
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("envs not loaded: %s", err)
	}

	return &cfg
}

func LoadFromPath(path string) *Config {
	var cfg Config

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(err)
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("envs not loaded: %s", err)
		}
		return &cfg
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("file envs not loaded: %s", err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	return path
}
