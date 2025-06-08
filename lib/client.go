package lib

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/mdp/qrterminal/v3"
	"github.com/tulir/whatsmeow"
	"github.com/tulir/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var Client *whatsmeow.Client

func InitClient() {
	dbLog := waLog.Noop // kosongkan logger (tidak spam)
	ctx := context.Background()

	container, err := sqlstore.New("sqlite", "file:session.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Fatalf("Device store error: %v", err)
	}

	Client = whatsmeow.NewClient(deviceStore, waLog.Noop)

	Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *whatsmeow.events.ConnectionOpened:
			log.Println("[OK] Bot connected as:", Client.Store.ID.User)
		case *whatsmeow.events.Disconnected:
			log.Println("[!] Disconnected")
		}
	})

	if Client.Store.ID == nil {
		pairResp, err := Client.PairPhone(ctx, "628xxx", false, whatsmeow.PairClientTypeClient, "")
		if err != nil {
			log.Fatalf("Pair error: %v", err)
		}

		qrterminal.Generate(pairResp, qrterminal.L, os.Stdout)
		log.Println("[!] Scan QR Code above to login")
	} else {
		err = Client.Connect()
		if err != nil {
			log.Fatalf("Connect error: %v", err)
		}
	}
}
