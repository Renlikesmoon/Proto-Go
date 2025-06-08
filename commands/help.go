package commands

import (
	"context"
	"fmt" // Added for error printing

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	// "google.golang.org/protobuf/proto" // This import is not needed for simple string pointers
)

// HelpCommand implements the Command interface to list available commands.
type HelpCommand struct{}

// Prefix returns the command's prefix.
func (c *HelpCommand) Prefix() string {
	return "!help"
}

// Run executes the Help command, sending a list of commands back to the sender.
func (c *HelpCommand) Run(evt *events.Message, client *whatsmeow.Client) {
	// Define the help message.
	text := "📖 Commands available:\n!ping - Cek koneksi\n!time - Waktu server\n!help - Lihat bantuan\n!anime <query> - Cari informasi anime"

	// Send the message back to the chat where the message originated.
	// The Conversation field expects a *string, so we pass the address of 'text'.
	_, err := client.SendMessage(context.Background(), evt.Info.Chat, &waProto.Message{
		Conversation: &text, // Corrected: Use &text to get a pointer to the string
	})

	if err != nil {
		fmt.Printf("❌ Gagal mengirim pesan bantuan: %v\n", err)
	}
}
