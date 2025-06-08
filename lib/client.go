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

	dbLog := waLog.Noop // logger minimal agar tidak spam console
	container, err := sqlstore.New(ctx, "sqlite", "file:session.db?_foreign_keys=on", dbLog)
	if err != nil {
		fmt.Println("❌ Gagal konek database:", err)
		os.Exit(1)
	}

	device, err := container.GetFirstDevice(ctx)
	if err != nil {
		fmt.Println("❌ Gagal ambil device:", err)
		os.Exit(1)
	}

	Client = whatsmeow.NewClient(device, dbLog)

	// Event handler contoh: tampilkan pesan masuk dan disconnect
	Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Println("📩 Pesan masuk dari:", v.Info.Sender.String())
		case *events.Disconnected:
			fmt.Println("🔌 Terputus dari WhatsApp")
		}
	})

	if Client.Store.ID == nil {
		// Belum login, harus pairing
		resp, err := Client.PairPhone(ctx, "", false, "client", "")
		if err != nil {
			fmt.Println("❌ Gagal pairing:", err)
			return
		}

		qrterminal.GenerateHalfBlock(resp, qrterminal.L, os.Stdout)
		fmt.Println("✅ Scan QR di atas dengan WhatsApp kamu.")
	} else {
		// Sudah login, langsung connect
		err = Client.Connect()
		if err != nil {
			fmt.Println("❌ Gagal konek ke WhatsApp:", err)
			return
		}
		fmt.Println("✅ Terhubung ke WhatsApp sebagai", Client.Store.ID.User)
	}
}
