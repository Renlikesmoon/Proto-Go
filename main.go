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

// StartClient initializes the Whatsmeow client and handles authentication.
// It now accepts a phoneNumber string to be used during the pairing process.
func StartClient(phoneNumber string) error { // Added phoneNumber parameter and error return
	ctx := context.Background()

	// Use no-op logger to avoid spamming logs
	dbLog := waLog.Noop

	// Create SQLite session database connection
	container, err := sqlstore.New(ctx, "sqlite", "file:session.db?_foreign_keys=on", dbLog)
	if err != nil {
		return fmt.Errorf("❌ Gagal konek database: %w", err) // Return error instead of os.Exit
	}

	// Get the first device from the store
	device, err := container.GetFirstDevice(ctx)
	if err != nil {
		return fmt.Errorf("❌ Gagal ambil device: %w", err) // Return error
	}

	// Initialize Whatsmeow client with device and logger
	Client = whatsmeow.NewClient(device, dbLog)

	// Event handler for incoming events
	Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Println("📩 Pesan masuk dari:", v.Info.Sender.String())
			// You'll likely want to route this to your commands.HandleMessage here
			// For example: commands.HandleMessage(v) (make sure to import "wa_bot/commands")
		case *events.Disconnected:
			fmt.Println("🔌 Terputus dari WhatsApp")
		}
	})

	// If no session ID exists, start the pairing process with QR code
	if Client.Store.ID == nil {
		// THIS IS WHERE YOUR PROVIDED SNIPPET FITS
		resp, err := Client.PairPhone(
			ctx,                     // The context for the operation
			phoneNumber,             // The phone number you want to pair with
			false,                   // showPushNotification (usually false for terminal bots)
			whatsmeow.PairClientChrome, // The client type (e.g., Chrome, Firefox, etc.)
			"Go Bot (Desktop)",      // A display name for the client (e.g., "Chrome (Linux)")
		)
		if err != nil {
			return fmt.Errorf("❌ Gagal pairing: %w", err) // Return error
		}

		// Display QR code in the terminal
		qrterminal.GenerateHalfBlock(resp, qrterminal.L, os.Stdout)
		fmt.Println("✅ Scan QR di atas dengan WhatsApp kamu.")
		fmt.Printf("Kode Pairing: %s\n", resp) // Also show pairing code for convenience
	} else {
		// If session exists, just connect
		err = Client.Connect()
		if err != nil {
			return fmt.Errorf("❌ Gagal konek ke WhatsApp: %w", err) // Return error
		}
		fmt.Println("✅ Terhubung ke WhatsApp sebagai", Client.Store.ID.User)
	}
	return nil // No error
}
