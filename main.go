package main

import (
	"log"
	"os"

	"github.com/username/Proto-Go/config"
	"github.com/username/Proto-Go/lib"
	"github.com/username/Proto-Go/commands"

	"go.mau.fi/whatsmeow/types/events"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go 6281234567890")
	}
	phone := os.Args[1]

	err := lib.StartWhatsAppWithCode(phone)
	if err != nil {
		log.Fatal("Error starting WhatsApp:", err)
	}

	lib.Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			if v.Info.MessageSource.IsFromMe {
				return
			}
			commands.HandleMessage(v)
		}
	})

	select {} // keep running
}
