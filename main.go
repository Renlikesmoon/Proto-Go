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

	js := jsonstore.NewJSONStore("session.json")

	client := whatsmeow.NewClient(nil, logger)

	// coba load session
	data, err := js.Load()
	if err == nil {
		err = client.Store.Restore(data)
		if err != nil {
			fmt.Println("❌ Gagal restore session:", err)
		}
	} else {
		fmt.Println("ℹ️ Session file belum ada atau gagal load:", err)
	}

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

		// Simpan session baru
		data, err := client.Store.Serialize()
		if err == nil {
			err = js.Save(data)
			if err != nil {
				fmt.Println("❌ Gagal simpan session:", err)
			}
		}
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
