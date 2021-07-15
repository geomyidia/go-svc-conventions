package util

import (
	"os"
	"os/signal"
)

func HandleSignal(handler func(int, os.Signal), signals ...os.Signal) {
	signalHandler := make(chan os.Signal, 1)
	signal.Notify(signalHandler, signals...)
	s := <-signalHandler
	handler(os.Getpid(), s)
}
