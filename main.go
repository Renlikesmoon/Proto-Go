package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/mem"
	waLog "go.mau.fi/whatsmeow/util/log"

	"go.mau.fi/whatsmeow/types/events"
	"github.com/mdp/qrterminal/v3"
)

var Client *whatsmeow.Client
var sessionFile = "session.json"

func main() {
	err := StartClient()
	if err != nil {
		fmt.Println("❌ Error:", err)
	}
}

func StartClient() error {
	ctx := context.Background()
	dbLog := waLog.Noop

	// In-memory store
	store := mem.NewStore()

	// Buat client kosong untuk awal
	Client = whatsmeow.NewClient(store.NewDevice(), dbLog)

	// Coba load session dari file JSON
	if _, err := os.Stat(sessionFile); err == nil {
		err := LoadSession()
		if err != nil {
			return fmt.Errorf("❌ Gagal load session: %w", err)
		}
	}

	Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Println("📩 Pesan masuk dari:", v.Info.Sender.String())
		case *events.Disconnected:
			fmt.Println("🔌 Terputus dari WhatsApp")
		}
	})

	if Client.Store.ID == nil {
		resp, err := Client.Pair(ctx)
		if err != nil {
			return fmt.Errorf("❌ Gagal pairing: %w", err)
		}
		qrterminal.GenerateHalfBlock(resp, qrterminal.L, os.Stdout)
		fmt.Println("✅ Scan QR di atas dengan WhatsApp kamu.")

		// Tunggu sampai paired sebelum simpan session
		Client.AddEventHandler(func(evt interface{}) {
			if _, ok := evt.(*events.PairSuccess); ok {
				err := SaveSession()
				if err != nil {
					fmt.Println("❌ Gagal simpan session:", err)
				} else {
					fmt.Println("💾 Session disimpan ke", sessionFile)
				}
			}
		})
	} else {
		err := Client.Connect()
		if err != nil {
			return fmt.Errorf("❌ Gagal konek ke WhatsApp: %w", err)
		}
		fmt.Println("✅ Terhubung ke WhatsApp sebagai", Client.Store.ID.User)
	}

	return nil
}

// Simpan session ke file JSON
func SaveSession() error {
	data, err := Client.Store.(*mem.Device).Marshal()
	if err != nil {
		return err
	}
	return os.WriteFile(sessionFile, data, 0600)
}

// Load session dari file JSON
func LoadSession() error {
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return err
	}
	return Client.Store.(*mem.Device).Unmarshal(data)
}
