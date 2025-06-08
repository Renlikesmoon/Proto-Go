package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"github.com/Renlikesmoon/Proto-Go/jsonstore"
)

var Client *whatsmeow.Client

func main() {
	err := StartClient()
	if err != nil {
		fmt.Println("❌ Error:", err)
	}
}

func StartClient() error {
	ctx := context.Background()
	dbLog := waLog.Noop

	js, err := jsonstore.NewJSONStore("session.json")
	if err != nil {
		return err
	}

	Client = whatsmeow.NewClient(js.GetDevice(), dbLog)

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

		Client.AddEventHandler(func(evt interface{}) {
			if _, ok := evt.(*events.PairSuccess); ok {
				js.Save()
				fmt.Println("💾 Session disimpan ke session.json")
			}
		})
	} else {
		err = Client.Connect()
		if err != nil {
			return fmt.Errorf("❌ Gagal konek ke WhatsApp: %w", err)
		}
		fmt.Println("✅ Terhubung ke WhatsApp sebagai", Client.Store.ID.User)
	}

	return nil
}
