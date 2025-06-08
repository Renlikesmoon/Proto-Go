package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Renlikesmoon/Proto-Go/jsonstore"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	ctx := context.Background()

	logger := waLog.Noop

	// Inisialisasi store JSON
	js, err := jsonstore.NewJSONStore("session.json")
	if err != nil {
		fmt.Println("❌ Gagal load session:", err)
		return
	}

	device := js.GetDevice()
	if device == nil {
		// Jika tidak ada session, buat device baru
		memStore := mem.New()
		device = store.NewDevice(memStore, nil)
		js.SetDevice(device)
	}

	client := whatsmeow.NewClient(device, logger)

	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Println("📩 Pesan dari:", v.Info.Sender.String())
		case *events.Disconnected:
			fmt.Println("🔌 Terputus dari WhatsApp")
		}
	})

	// Pair jika belum login
	if client.Store.ID == nil {
		resp, err := client.PairPhone(ctx, "6285954540177", false, whatsmeow.PairClientChrome, "GoBot")
		if err != nil {
			fmt.Println("❌ Gagal pairing:", err)
			return
		}
		fmt.Println("✅ Silakan scan QR dengan WhatsApp kamu.")
		fmt.Println("Pairing code:", resp)

		// Simpan sesi
		if err := js.Save(); err != nil {
			fmt.Println("❌ Gagal simpan session:", err)
		}
	}

	// Connect
	err = client.Connect()
	if err != nil {
		fmt.Println("❌ Gagal konek:", err)
		return
	}
	fmt.Println("✅ Terhubung sebagai", client.Store.ID.User)

	// Tunggu Ctrl+C
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	fmt.Println("👋 Keluar dari bot.")
	client.Disconnect()
}
