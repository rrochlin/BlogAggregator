package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	var cfg Config
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	dat, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(dat, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil

}

func (cfg *Config) SetUser(name string) error {
	cfg.CurrentUserName = name
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	dat, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, dat, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(path, configFileName), nil

}
