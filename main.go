package main

import (
	"clippy/clippy"
	"clippy/config"
	"os"
	"os/signal"
)

func main() {
	go clippy.Run()

	err := config.LoadConfig()
	if err != nil {
		clippy.DisplayNotification("Error", err.Error())
		return
	}

	clippy.DisplayNotification("Info", "Copy a question, then press F8 to ask Clippy!")

	go clippy.ListenForHotkey()

	waitForExitSignal()
}

func waitForExitSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	<-sigCh
	clippy.Quit()
}
