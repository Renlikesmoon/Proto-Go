package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	ctx := context.Background()
	logger := waLog.Noop
	client := whatsmeow.NewClient(nil, logger)

	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Println("📩 Pesan dari:", v.Info.Sender.String())
		case *events.Disconnected:
			fmt.Println("🔌 Terputus dari WhatsApp")
		}
	})

	qrChan, _ := client.GetQRChannel(ctx)
	err := client.Connect()
	if err != nil {
		fmt.Println("❌ Gagal konek:", err)
		return
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			fmt.Println("✅ Scan QR dengan WhatsApp kamu:", evt.Code)
		} else {
			fmt.Println("❌ QR error:", evt.Code)
		}
	}

	fmt.Println("✅ Terhubung")

	// Tunggu Ctrl+C untuk exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	fmt.Println("👋 Keluar.")
	client.Disconnect()
}
