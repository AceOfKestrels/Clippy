package main

import (
	"clippy/clippy"
	"clippy/config"
	"os"
	"os/signal"
)

func main() {
	err := config.Load()
	if err != nil {
		clippy.HandleError(err)
		return
	}

	go clippy.Run()
	go clippy.ListenForHotkey()
	waitForExitSignal()
}

func waitForExitSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	<-sigCh
	clippy.Quit()
}
