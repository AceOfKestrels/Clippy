package main

import (
	"bytes"
	"clippy/config"
	"clippy/prompt"
	"clippy/response"
	"encoding/json"
	"fmt"
	"github.com/gen2brain/beeep"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
	"github.com/kindlyfire/go-keylogger"
)

var apiKey string

func main() {
	go func() {
		systray.Run(onReady, onExit)
	}()

	err := loadConfig()
	if err != nil {
		displayNotification("Error", err.Error())
		return
	}

	displayNotification("Info", "Copy a question, then press F8 to ask Clippy!")

	go listenForHotkeys()

	waitForSignal()
}

func onReady() {
	icon, err := getIcon("icon.ico")
	if err == nil {
		systray.SetIcon(icon)
	}
	systray.SetTitle("Clippy")
	systray.SetTooltip("Clippy")

	mQuit := systray.AddMenuItem("Quit", "Quit the app")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		os.Exit(0)
	}()
}

func onExit() {
	log.Println("Exiting application")
}

func getIcon(filePath string) ([]byte, error) {
	// Read the icon file into memory
	icon, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read icon file: %w", err)
	}
	return icon, nil
}

func listenForHotkeys() {
	kl := keylogger.NewKeylogger()
	log.Println("Listening for F8...")

	//var ctrlPressed bool
	for {
		key := kl.GetKey()
		if key.Empty {
			time.Sleep(50 * time.Millisecond)
			continue
		}

		if key.Keycode == 119 {
			log.Printf("Asking gemini...")
			handleClipboard()
		}
	}
}

func handleClipboard() {
	content, err := clipboard.ReadAll()
	if err != nil {
		log.Printf("Clipboard error: %v", err)
		return
	}

	promptStr := prompt.New(content)

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=" + apiKey
	resp, err := http.Post(url, "application/json", bytes.NewBufferString(promptStr))
	if err != nil {
		log.Printf("API error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		displayNotification("error", "api response does not indicate success: "+resp.Status)
		return
	}

	r, err := response.Deserialize(resp)
	if err != nil {
		log.Printf("API error: %v", err)
	}

	displayNotification("Clippy says:", r.Candidates[0].Content.Parts[0].Text)
}

func loadConfig() error {
	// Check if the config file exists
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		// If it doesn't exist, create a default config
		defaultConfig := config.Config{
			APIKey: "your-api-key-here",
		}
		if err = saveConfig(defaultConfig); err != nil {
			return err
		}
		return fmt.Errorf("created default config - please restart the application")
	}

	// Read the config file
	data, err := os.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse the JSON into the Config struct
	var configFile config.Config
	err = json.Unmarshal(data, &configFile)
	if err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	apiKey = configFile.APIKey
	return nil
}

func saveConfig(config config.Config) error {
	// Convert the Config struct to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write the JSON data to the config file
	err = os.WriteFile("config.json", data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	log.Println("Config saved.")
	return nil
}

func waitForSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	<-sigCh
	systray.Quit()
}

func displayNotification(title, message string) {
	// Show a notification using the beeep package
	err := beeep.Notify(title, message, "")
	if err != nil {
		log.Printf("Error showing notification: %v", err)
	}
}
