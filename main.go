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
		clippy.HandleError(err)
		return
	}

	go clippy.ListenForHotkey()

	waitForExitSignal()
}

func waitForExitSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	<-sigCh
	clippy.Quit()
}
