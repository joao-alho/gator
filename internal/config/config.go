package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	config_file_name = ".gatorconfig.json"
)

type Config struct {
	config_path string
	DBURL       string `json:"db_url"`
	CurrentUser string `json:"current_user"`
}

func getConfigPath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("Missing HOME environment variable!")
	}
	home_path := strings.TrimRight(home_dir, "/")
	config_path := fmt.Sprintf("%s/%s", home_path, config_file_name)
	return config_path, nil
}

func writeConfig(cfg *Config) error {
	json_data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	config_path, err := getConfigPath()
	if err != nil {
		return err
	}
	if err := os.WriteFile(config_path, json_data, os.FileMode(0666)); err != nil {
		msg := fmt.Sprintf("Something went wrong writting to file %s", config_path)
		return errors.New(msg)
	}
	return nil
}
func Read() (Config, error) {
	config_path, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}
	config_data, err := os.ReadFile(config_path)
	if err != nil {
		return Config{}, errors.New("Could not find file")
	}
	var config Config
	if err := json.Unmarshal(config_data, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUser = username
	if err := writeConfig(c); err != nil {
		return err
	}
	return nil
}
