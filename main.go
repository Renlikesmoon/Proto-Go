package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"github.com/Renlikesmoon/Proto-Go/jsonstore" // Pastikan sesuai dengan struktur proyek Anda
)

func main() {
	ctx := context.Background()

	// Konfigurasi pairing
	phoneNumber := "6285954540177"               // Ganti dengan nomor Anda
	clientType := whatsmeow.PairClientChrome     // Bisa juga Firefox, Edge, Safari
	clientName := "Go Bot (Desktop)"             // Nama klien WhatsApp

	// Load session dari file JSON
	store, err := jsonstore.NewJSONStore("session.json")
	if err != nil {
		fmt.Println("❌ Gagal load session:", err)
		return
	}

	logger := waLog.Noop
	client := whatsmeow.NewClient(store.GetDevice(), logger)

	// Event handler
	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Println("📩 Pesan dari:", v.Info.Sender.String())
		case *events.Disconnected:
			fmt.Println("🔌 Terputus dari WhatsApp")
		}
	})

	// Jika belum login, lakukan pairing otomatis
	if client.Store.ID == nil {
		resp, err := client.PairPhone(ctx, phoneNumber, false, clientType, clientName)
		if err != nil {
			fmt.Println("❌ Pairing gagal:", err)
			return
		}
		fmt.Println("✅ Pairing berhasil. Kode:", resp)

		// Simpan session
		if err := store.Save(); err != nil {
			fmt.Println("❌ Gagal simpan session:", err)
		}
	}

	// Connect
	if err := client.Connect(); err != nil {
		fmt.Println("❌ Gagal konek:", err)
		return
	}

	fmt.Println("✅ Terhubung ke WhatsApp sebagai", client.Store.ID.User)

	// Tunggu Ctrl+C untuk keluar
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("👋 Keluar dari aplikasi.")
	client.Disconnect()
}
