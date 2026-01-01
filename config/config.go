package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/xxxsen/common/logger"
)

type Config struct {
	Tailnet         string           `json:"tailnet"`
	APIKey          string           `json:"api_key"`
	Listen          string           `json:"listen"`
	RefreshInterval int64            `json:"refresh_interval"`
	Log             logger.LogConfig `json:"log"`
}

func Load(f string) (*Config, error) {
	raw, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := json.Unmarshal(raw, cfg); err != nil {
		return nil, err
	}
	if cfg.Tailnet == "" || cfg.APIKey == "" {
		return cfg, errors.New("tailnet and tskey are required")
	}
	if cfg.Listen == "" {
		cfg.Listen = ":8080"
	}
	if cfg.RefreshInterval <= 0 {
		cfg.RefreshInterval = 60 * 10 //second
	}
	return cfg, nil
}
