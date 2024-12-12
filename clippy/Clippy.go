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
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const f8Key = 119

var url string

func Run() {
	url = fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%v:generateContent?key=%v", config.ConfigFile.Model, config.ConfigFile.ApiKey)

	displayNotification("Welcome", "Copy a question, then press F8 to ask Clippy!\n\nUsing model: "+config.ConfigFile.Model)
	systray.Run(onReady, onExit)
}

func Quit() {
	displayNotification("Goodbye", "Thank you for using Clippy")
	systray.Quit()
}

func onReady() {
	systray.SetIcon(icon)
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

func ListenForHotkey() {
	kl := keylogger.NewKeylogger()
	log.Println("Listening for F8...")

	for {
		key := kl.GetKey()
		if key.Empty {
			time.Sleep(50 * time.Millisecond)
			continue
		}

		if key.Keycode == f8Key {
			handleClipboard()
		}
	}
}

func handleClipboard() {
	content, err := clipboard.ReadAll()
	if err != nil {
		HandleError(err)
		return
	}

	log.Printf("Asking gemini...")

	promptStr := prompt.New(content)

	displayNotification("Thinking...", "Please give me a moment.")

	resp, err := http.Post(url, "application/json", bytes.NewBufferString(promptStr))
	if err != nil {
		HandleError(err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		HandleError(fmt.Errorf("api response does not indicate success: %v", resp.Status))
		return
	}

	r, err := response.Deserialize(resp)
	if err != nil {
		HandleError(err)
	}

	responseText := r.Candidates[0].Content.Parts[0].Text
	displayNotification("Clippy says:", responseText)
	err = clipboard.WriteAll(responseText)
	if err != nil {
		HandleError(err)
	}
}

func HandleError(err error) {
	log.Println(err)
	displayNotification("Error", err.Error())
}

func displayNotification(title, message string) {
	err := beeep.Notify(title, message, "")
	if err != nil {
		log.Printf("Error showing notification: %v", err)
	}
}
