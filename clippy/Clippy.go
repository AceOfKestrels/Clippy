package clippy

import (
	"bytes"
	"clippy/clippy/errors/contextCanceledError"
	"clippy/clippy/errors/deadlineExceededError"
	"clippy/clippy/icon"
	"clippy/config"
	"clippy/prompt"
	"clippy/response"
	"context"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
	"github.com/go-toast/toast"
	"github.com/kindlyfire/go-keylogger"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const AppName = "Clippy"

const f8Key = 119

var url string

var cancelLastRequest context.CancelFunc

func Run() {
	url = fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%v:generateContent?key=%v", config.Config.Model, config.Config.ApiKey)

	displayNotification("Welcome", "Copy a question, then press F8 to ask Clippy!\n\nUsing model: "+config.Config.Model, false)
	systray.Run(onReady, onExit)
}

func Quit() {
	displayNotification(
		"Goodbye",
		"Thank you for using Clippy",
		false,
		toast.Action{Type: "protocol", Label: "Github Source", Arguments: "https://github.com/AceOfKestrels/Clippy"},
		toast.Action{Type: "protocol", Label: "Report a bug", Arguments: "https://github.com/AceOfKestrels/Clippy/issues"},
	)
	systray.Quit()
}

func onReady() {
	systray.SetIcon(icon.Icon)
	systray.SetTitle(AppName)
	systray.SetTooltip(AppName)

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
			if cancelLastRequest != nil {
				cancelLastRequest()
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Config.Timeout)*time.Millisecond)
			cancelLastRequest = cancel
			go handleClipboard(ctx)
		}
	}
}

func handleClipboard(ctx context.Context) {
	content, err := clipboard.ReadAll()
	if err != nil {
		HandleError(err)
		return
	}

	log.Printf("Asking gemini...")

	promptStr := prompt.New(content)

	displayNotification("Thinking...", "Please give me a moment.", false)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(promptStr))
	if err != nil {
		HandleError(err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		if deadlineExceededError.Is(err) {
			HandleError(fmt.Errorf("no response within timeout"))
		} else if !contextCanceledError.Is(err) {
			HandleError(err)
		}
	}
	if resp == nil {
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		HandleError(fmt.Errorf("api response does not indicate success: %v", resp.Status))
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
		return
	}

	r, err := response.Deserialize(resp)
	if err != nil {
		HandleError(err)
	}

	responseText := r.Candidates[0].Content.Parts[0].Text
	displayNotification("Clippy says:", responseText, true)
	err = clipboard.WriteAll(responseText)
	if err != nil {
		HandleError(err)
	}
}

func HandleError(err error) {
	log.Println(err)
	displayNotification("Error", err.Error(), true)
}

func displayNotification(title, message string, ignoreMinimal bool, actions ...toast.Action) {
	if config.Config.Minimal && !ignoreMinimal {
		return
	}

	notification := toast.Notification{
		AppID:   AppName,
		Title:   title,
		Message: message,
		Actions: actions,
	}

	err := notification.Push()
	if err != nil {
		log.Printf("Error showing notification: %v", err)
	}
}
