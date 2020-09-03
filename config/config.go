package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Storage struct {
		DBType   string `json:"db_type"`
		DBDsn    string `json:"db_dsn"`
		RedisUrl string `json:"redis_url"`
	} `json:"storage"`

	ServeOn struct {
		Addr string `json:"addr,omitempty"`
		Port string `json:"port"`
	} `json:"serve_on"`
}

func FromFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func MustFromFile(path string) *Config {
	c, err := FromFile(path)
	if err != nil {
		log.Panicln(err)
	}
	return c
}
