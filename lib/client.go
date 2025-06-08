package lib

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

func StartClient() {
	ctx := context.Background()

	// Gunakan logger kosong agar tidak spam log
	dbLog := waLog.Noop

	// Buat koneksi database session SQLite
	container, err := sqlstore.New(ctx, "sqlite", "file:session.db?_foreign_keys=on", dbLog)
	if err != nil {
		fmt.Println("❌ Gagal konek database:", err)
		os.Exit(1)
	}

	// Ambil device pertama dari store
	device, err := container.GetFirstDevice(ctx)
	if err != nil {
		fmt.Println("❌ Gagal ambil device:", err)
		os.Exit(1)
	}

	// Inisialisasi client Whatsmeow dengan device dan logger
	Client = whatsmeow.NewClient(device, dbLog)

	// Event handler untuk menangani event masuk
	Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Println("📩 Pesan masuk dari:", v.Info.Sender.String())
		case *events.Disconnected:
			fmt.Println("🔌 Terputus dari WhatsApp")
		}
	})

	// Jika belum ada session ID, mulai proses pairing dengan QR
	if Client.Store.ID == nil {
		// Mengganti "client" dengan whatsmeow.PairClientChrome
		resp, err := Client.PairPhone(ctx, "", false, whatsmeow.PairClientChrome, "")
		if err != nil {
			fmt.Println("❌ Gagal pairing:", err)
			return
		}

		// Tampilkan QR code di terminal untuk discan
		qrterminal.GenerateHalfBlock(resp, qrterminal.L, os.Stdout)
		fmt.Println("✅ Scan QR di atas dengan WhatsApp kamu.")
	} else {
		// Jika sudah ada session, langsung konek
		err = Client.Connect()
		if err != nil {
			fmt.Println("❌ Gagal konek ke WhatsApp:", err)
			return
		}
		fmt.Println("✅ Terhubung ke WhatsApp sebagai", Client.Store.ID.User)
	}
}
