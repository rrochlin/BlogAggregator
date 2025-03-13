package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	var cfg Config
	path, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	dat, err := os.ReadFile(path + "/" + configFileName)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(dat, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil

}

func (cfg *Config) SetUser(name string) {
	cfg.CurrentUserName = name
}
