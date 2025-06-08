package lib

import (
	"context"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	_ "modernc.org/sqlite"
)

var Client *whatsmeow.Client

func StartWhatsAppWithCode(phone string) error {
	dbLog := store.NoopLogger{}

	container, err := sqlstore.New("sqlite", "file:session.db?_foreign_keys=on", dbLog)
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return fmt.Errorf("device error: %w", err)
	}

	Client = whatsmeow.NewClient(deviceStore, dbLog)

	Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Ready:
			fmt.Println("✅ WhatsApp connected as", Client.Store.PushName)
		case *events.Disconnected:
			fmt.Println("⚠️ Disconnected:", v.Reason)
		}
	})

	if Client.Store.ID == nil {
		err = Client.Connect()
		if err != nil {
			return fmt.Errorf("connect error: %w", err)
		}

		// Pair with code
		resp, err := Client.PairPhone(context.Background(), phone, whatsmeow.PairClientChrome)
		if err != nil {
			return fmt.Errorf("pairing error: %w", err)
		}

		fmt.Println("🔑 Kode pairing:", resp.Code)
		fmt.Println("➡️ Masukkan kode ini di WhatsApp Web (Link with code).")
		return nil
	} else {
		err = Client.Connect()
		if err != nil {
			return fmt.Errorf("reconnect error: %w", err)
		}
		fmt.Println("✅ WhatsApp reconnected as", Client.Store.PushName)
	}

	return nil
}
