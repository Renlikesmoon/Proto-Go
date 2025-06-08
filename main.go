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

	if client.Store.ID == nil {
		resp, err := client.PairPhone(ctx, "6285954540177", false, whatsmeow.PairClientChrome, "GoBot")
		if err != nil {
			fmt.Println("❌ Gagal pairing:", err)
			return
		}

		fmt.Println("✅ Scan QR dengan WhatsApp kamu:")
		fmt.Println(resp)

	}

	err = client.Connect()
	if err != nil {
		fmt.Println("❌ Gagal konek:", err)
		return
	}
	fmt.Println("✅ Terhubung sebagai", client.Store.ID.User)

	// Tunggu Ctrl+C untuk exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	fmt.Println("👋 Keluar.")
	client.Disconnect()
}
