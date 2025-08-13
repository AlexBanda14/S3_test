package config

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Endpoint     string `env:"ENDPOINT"`
	AccessKey    string `env:"ACCESSKEYID"`
	SecretKey    string `env:"SECRETACCESSKEY"`
	BucketName   string `env:"BUCKETNAME"`
	PathUpload   string `env:"UPLOADFOLDER"`
	PathDownload string `env:"DOWNLOADFOLDER"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load(".env")
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config error: %w", err)
	}
	return cfg, nil
}
