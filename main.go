package main

import (
    "log"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/store/sqlstore"
    _ "modernc.org/sqlite"

    "github.com/mdp/qrterminal/v3"
    "context"
)

func main() {
    dbLog := log.Default()
    container, err := sqlstore.New("sqlite", "file:session.db?_foreign_keys=on", dbLog)
    if err != nil {
        log.Fatal(err)
    }

    deviceStore, err := container.GetFirstDevice()
    if err != nil {
        log.Fatal(err)
    }

    client := whatsmeow.NewClient(deviceStore, nil)

    if client.Store.ID == nil {
        qrChan, _ := client.GetQRChannel(context.Background())
        err = client.Connect()
        if err != nil {
            log.Fatal(err)
        }

        for evt := range qrChan {
            if evt.Event == "code" {
                qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, log.Writer())
            } else {
                log.Println("QR Event:", evt.Event)
            }
        }
    } else {
        err = client.Connect()
        if err != nil {
            log.Fatal(err)
        }
    }

    client.AddEventHandler(func(evt interface{}) {
        messageHandler(evt, client)
    })

    select {} // block
}
