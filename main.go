package main

import (
	"log"
	"os"

	// Ensure these import paths are correct for your module setup.
	// Replace "github.com/Renlikesmoon/Proto-Go" with your actual module name if different.
	"github.com/Renlikesmoon/Proto-Go/commands" // For commands.HandleMessage
	"github.com/Renlikesmoon/Proto-Go/lib"      // For lib.StartClient and lib.Client
	// "github.com/Renlikesmoon/Proto-Go/config" // Config is used by commands, not directly in main, so this import is often not strictly needed here but doesn't hurt.

	"go.mau.fi/whatsmeow/types/events" // For whatsmeow event types
)

// main is the entry point of the Go application.
func main() {
	// Check if a phone number argument is provided.
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run . <phone_number>")
	}
	// Get the phone number from the command line arguments.
	phone := os.Args[1]

	// Start the Whatsmeow client. This function handles connecting,
	// pairing (if needed), and returning any errors.
	err := lib.StartClient(phone) // This calls the StartClient function from your lib package
	if err != nil {
		log.Fatalf("Error starting WhatsApp client: %v", err)
	}

	// Add an event handler to the Whatsmeow client to process incoming events.
	// This runs in a separate goroutine managed by whatsmeow itself.
	lib.Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			// Ignore messages sent by the bot itself to prevent infinite loops.
			if v.Info.MessageSource.IsFromMe {
				return
			}
			// Route incoming messages to the commands handler.
			commands.HandleMessage(v)
		case *events.Connected:
			// Log when the WhatsApp client successfully connects.
			log.Println("WhatsApp client connected!")
		case *events.Disconnected:
			// Log when the WhatsApp client disconnects.
			// Whatsmeow library typically handles automatic reconnection.
			log.Println("WhatsApp client disconnected.")
			// You could add custom reconnection logic here if needed,
			// but for most cases, whatsmeow's built-in handling is sufficient.
		case *events.QR:
			// Your lib.StartClient already handles displaying the QR code.
			// We'll just log that the QR event occurred here.
			log.Println("Received QR code event (QR displayed by lib.StartClient).")
		case *events.PairingCode:
			// Your lib.StartClient already handles displaying the pairing code.
			// We'll just log that the PairingCode event occurred here.
			log.Println("Received Pairing Code event (code displayed by lib.StartClient).")
		// Add other event types you might want to log or handle here.
		}
	})

	// This blocks the main goroutine indefinitely, keeping the application running.
	// Whatsmeow operates in its own goroutines for networking and event handling.
	select {}
}
