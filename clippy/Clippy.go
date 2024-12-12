package clippy

import (
	"bytes"
	"clippy/config"
	"clippy/prompt"
	"clippy/response"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/kindlyfire/go-keylogger"
	"log"
	"net/http"
	"os"
	"time"
)

func Run() {
	systray.Run(onReady, onExit)
}

func Quit() {
	systray.Quit()
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
	icon, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read icon file: %w", err)
	}
	return icon, nil
}

func ListenForHotkey() {
	kl := keylogger.NewKeylogger()
	log.Println("Listening for F8...")

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

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=" + config.ApiKey
	resp, err := http.Post(url, "application/json", bytes.NewBufferString(promptStr))
	if err != nil {
		log.Printf("API error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		DisplayNotification("error", "api response does not indicate success: "+resp.Status)
		return
	}

	r, err := response.Deserialize(resp)
	if err != nil {
		log.Printf("API error: %v", err)
	}

	DisplayNotification("Clippy says:", r.Candidates[0].Content.Parts[0].Text)
}

func DisplayNotification(title, message string) {
	err := beeep.Notify(title, message, "")
	if err != nil {
		log.Printf("Error showing notification: %v", err)
	}
}
