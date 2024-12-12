package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var ApiKey string

type Config struct {
	APIKey string `json:"apiKey"`
}

func LoadConfig() error {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		defaultConfig := Config{
			APIKey: "your-api-key-here",
		}
		if err = saveConfig(defaultConfig); err != nil {
			return err
		}
		return fmt.Errorf("created default config - please restart the application")
	}

	data, err := os.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	var configFile Config
	err = json.Unmarshal(data, &configFile)
	if err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	ApiKey = configFile.APIKey
	return nil
}

func saveConfig(config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	err = os.WriteFile("config.json", data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	log.Println("Config saved.")
	return nil
}
