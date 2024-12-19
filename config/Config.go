package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var ConfigFile Config

const fileName = "config.json"
const defaultModel = "gemini-1.5-flash-latest"
const defaultTimeout = 10000

type Config struct {
	ApiKey  string `json:"apiKey"`
	Model   string `json:"model"`
	Timeout int    `json:"requestTimeout"`
	Minimal bool   `json:"minimalNotifications"`
}

func LoadConfig() error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if err = saveDefaultConfig(); err != nil {
			return err
		}
		return fmt.Errorf("created default config - please restart the application")
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	var configFile Config
	err = json.Unmarshal(data, &configFile)
	if err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	if len(configFile.ApiKey) == 0 {
		return fmt.Errorf("config file contains no api key")
	}

	if len(configFile.Model) == 0 {
		configFile.Model = defaultModel
	}

	if configFile.Timeout <= 0 {
		configFile.Timeout = defaultTimeout
	}

	ConfigFile = configFile
	return nil
}

func saveDefaultConfig() error {
	defaultConfig := Config{
		ApiKey:  "your-api-key-here",
		Model:   defaultModel,
		Timeout: defaultTimeout,
		Minimal: false,
	}

	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	log.Println("Config saved.")
	return nil
}
