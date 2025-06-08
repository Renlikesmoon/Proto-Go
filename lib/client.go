package lib

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	_ "modernc.org/sqlite"
)

var Client *whatsmeow.Client

func StartWhatsAppWithCode(phone string) error {
	// Gunakan silent logger
	dbLog := store.NoopLogger{}

	// Setup DB session
	container, err := sqlstore.New("sqlite", "file:session.db?_foreign_keys=on", dbLog)
	if err != nil {
		return fmt.Errorf("DB error: %w", err)
	}

	// Coba ambil device session pertama
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return fmt.Errorf("device error: %w", err)
	}

	Client = whatsmeow.NewClient(deviceStore, dbLog)
	Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Ready:
			fmt.Println("✅ WhatsApp connected as", Client.Store.PushName)
		}
	})

	if Client.Store.ID == nil {
		// Connect untuk pairing
		err = Client.Connect()
		if err != nil {
			return fmt.Errorf("connect error: %w", err)
		}

		// Pairing code (via nomor)
		pairResp, err := Client.PairPhone(context.Background(), phone, whatsmeow.PairClientChrome)
		if err != nil {
			return fmt.Errorf("pairing error: %w", err)
		}

		fmt.Println("🔑 Kode pairing:", pairResp.Code)
		fmt.Println("➡️  Masukkan kode ini di WhatsApp Web pada perangkat tersebut.")
		return nil
	} else {
		// Sudah login sebelumnya, reconnect
		err := Client.Connect()
		if err != nil {
			return fmt.Errorf("reconnect error: %w", err)
		}
		fmt.Println("✅ WhatsApp reconnected as", Client.Store.PushName)
	}

	return nil
}
  
