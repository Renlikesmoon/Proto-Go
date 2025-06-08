package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	logger := waLog.Noop
	client := whatsmeow.NewClient(nil, logger)

	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Println("Pesan dari:", v.Info.Sender.String())
		case *events.Disconnected:
			fmt.Println("Terputus dari WhatsApp")
		}
	})

	// Tunggu Ctrl+C untuk exit
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("Keluar.")
		client.Disconnect()
		os.Exit(0)
	}()

	err := client.Connect()
	if err != nil {
		fmt.Println("Gagal konek:", err)
		return
	}
}
