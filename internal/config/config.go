package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("getting config file path: %v", err)
	}
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("opening config file: %v", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("closing JSON file: %w", cerr)
		}
	}()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("reading config file: %v", err)
	}

	return cfg, err
}

func getConfigFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("getting user config directory: %v", err)
	}

	gatorDir := filepath.Join(configDir, "gator")
	if err := os.MkdirAll(gatorDir, 0755); err != nil {
		return "", fmt.Errorf("creating config directory: %w", err)
	}

	path := filepath.Join(gatorDir, configFileName)
	return path, err
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	if err := cfg.write(); err != nil {
		return fmt.Errorf("writing config file: %v", err)
	}

	return nil
}

func (cfg *Config) write() (err error) {
	path, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("getting config file path: %v", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("closing JSON file: %w", cerr)
		}
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return err
}
