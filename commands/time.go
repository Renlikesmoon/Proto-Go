package commands

import (
	"context"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	// "google.golang.org/protobuf/proto" // This import is not needed for simple string pointers
)

// TimeCommand implements the Command interface to display the current time.
type TimeCommand struct{}

// Prefix returns the command's prefix.
func (c *TimeCommand) Prefix() string {
	return "!time"
}

// Run executes the Time command, sending the current time back to the sender.
func (c *TimeCommand) Run(evt *events.Message, client *whatsmeow.Client) {
	// Format the current time into a readable string.
	now := time.Now().Format("2006-01-02 15:04:05 MST") // Added MST for timezone clarity
	msg := fmt.Sprintf("🕒 Current time: %s", now)

	// Send the message back to the chat where the message originated.
	// The Conversation field expects a *string, so we pass the address of 'msg'.
	_, err := client.SendMessage(context.Background(), evt.Info.Chat, &waProto.Message{
		Conversation: &msg, // Corrected: Use &msg to get a pointer to the string
	})

	if err != nil {
		fmt.Printf("❌ Gagal mengirim pesan waktu: %v\n", err)
	}
}
