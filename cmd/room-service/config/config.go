package config

import (
	"fmt"
	"os"
	"time"

	"github.com/tajogii/goWatch/pkg/storage"
	"gopkg.in/yaml.v3"
)

type cache struct {
	Ttl time.Duration `yaml:"ttl"`
}

type General struct {
	PublicPort int    `yaml:"public"`
	Level      string `yaml:"level"`
}

type DB struct {
	RoomDb *storage.PgConf `yaml:"roomDb"`
}

type Cache struct {
	RoomCache cache `yaml:"room"`
}

type Config struct {
	General General `yaml:"general"`
	DB      DB      `yaml:"db"`
	Cache   Cache   `yaml:"cache"`
}

func LoadConfig() (*Config, error) {
	buf, err := os.ReadFile(os.Getenv("CONFIG"))
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	cfg.DB.RoomDb.Host = os.Getenv("ROOM_DB_HOST")

	return &cfg, nil
}
