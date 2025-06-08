package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.mau.fi/whatsmeow"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	logger := waLog.Noop
	client := whatsmeow.NewClient(waLog.Stdout("Client"), logger)

	// Jangan langsung connect, tapi pair phone atau scan QR
	qrChan, err := client.GetQRChannel(context.Background())
	if err != nil {
		fmt.Println("Gagal mendapatkan QR channel:", err)
		return
	}

	go func() {
		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("Scan QR dengan WhatsApp kamu:", evt.Code)
			} else {
				fmt.Println("QR error:", evt.Code)
			}
		}
	}()

	err = client.Connect()
	if err != nil {
		fmt.Println("Gagal konek:", err)
		return
	}

	// Tunggu Ctrl+C untuk exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	fmt.Println("Keluar.")
	client.Disconnect()
}
