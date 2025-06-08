package main

import (
	"context"
	"fmt"
	"os"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.mau.fi/whatsmeow/types/events"

	"github.com/mdp/qrterminal/v3"
)

var Client *whatsmeow.Client

func main() {
	err := StartClient()
	if err != nil {
		fmt.Println("❌ Error:", err)
	}
}

// StartClient menginisialisasi klien Whatsmeow dan menangani autentikasi.
func StartClient() error {
	// --- Detail Pairing Ditetapkan Langsung di sini ---
	phoneNumber := "6285954540177"
	clientType := whatsmeow.PairClientChrome
	clientDisplayName := "Go Bot (Desktop)"
	// --- Akhir Hardcode Detail ---

	ctx := context.Background()
	dbLog := waLog.Noop

	container, err := sqlstore.New(ctx, "sqlite", "file:session.db?_foreign_keys=on", dbLog)
	if err != nil {
		return fmt.Errorf("❌ Gagal konek database: %w", err)
	}

	device, err := container.GetFirstDevice(ctx)
	if err != nil {
		return fmt.Errorf("❌ Gagal ambil device: %w", err)
	}

	Client = whatsmeow.NewClient(device, dbLog)

	Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Println("📩 Pesan masuk dari:", v.Info.Sender.String())
		case *events.Disconnected:
			fmt.Println("🔌 Terputus dari WhatsApp")
		}
	})

	if Client.Store.ID == nil {
		resp, err := Client.PairPhone(ctx, phoneNumber, false, clientType, clientDisplayName)
		if err != nil {
			return fmt.Errorf("❌ Gagal pairing: %w", err)
		}
		qrterminal.GenerateHalfBlock(resp, qrterminal.L, os.Stdout)
		fmt.Println("✅ Scan QR di atas dengan WhatsApp kamu.")
		fmt.Printf("Kode Pairing: %s\n", resp)
	} else {
		err = Client.Connect()
		if err != nil {
			return fmt.Errorf("❌ Gagal konek ke WhatsApp: %w", err)
		}
		fmt.Println("✅ Terhubung ke WhatsApp sebagai", Client.Store.ID.User)
	}

	return nil
}
