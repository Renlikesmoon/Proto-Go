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
	dbLog := waLog.Noop // tidak spam console
	container, err := sqlstore.New("sqlite", "file:session.db?_foreign_keys=on", dbLog)
	if err != nil {
		fmt.Println("❌ Gagal konek DB:", err)
		os.Exit(1)
	}

	device, err := container.GetFirstDevice(context.Background())
	if err != nil {
		fmt.Println("❌ Gagal ambil device:", err)
		os.Exit(1)
	}

	Client = whatsmeow.NewClient(device, dbLog)

	// Event handler
	Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Println("📩 Pesan masuk dari:", v.Info.Sender.String())
		case *events.Disconnected:
			fmt.Println("🔌 Terputus dari WhatsApp")
		}
	})

	if Client.Store.ID == nil {
		// Pairing baru
		resp, err := Client.PairPhone(context.Background(), "", false, whatsmeow.PairClientTypeClient, "")
		if err != nil {
			fmt.Println("❌ Gagal pairing:", err)
			return
		}

		qrterminal.GenerateHalfBlock(resp.URI)
		fmt.Println("✅ Scan QR di atas dengan WhatsApp kamu.")
	} else {
		// Langsung konek
		err = Client.Connect()
		if err != nil {
			fmt.Println("❌ Gagal konek ke WhatsApp:", err)
			return
		}
		fmt.Println("✅ Terhubung ke WhatsApp sebagai", Client.Store.ID.User)
	}
}
