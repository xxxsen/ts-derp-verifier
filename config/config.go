package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/xxxsen/common/logger"
)

type Config struct {
	Tailnet         string           `json:"tailnet"`
	ClientID        string           `json:"client_id"`
	ClientSecret    string           `json:"client_secret"`
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
	if cfg.Tailnet == "" || cfg.ClientID == "" || cfg.ClientSecret == "" {
		return cfg, errors.New("tailnet, client_id and client_secret are required")
	}
	if cfg.Listen == "" {
		cfg.Listen = ":8080"
	}
	if cfg.RefreshInterval <= 0 {
		cfg.RefreshInterval = 60 * 10 //second
	}
	return cfg, nil
}
