package commands

import (
	"context"
	"fmt" // Added for error printing

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	// The 'lib' import is no longer strictly necessary within this command
	// as the whatsmeow.Client will be passed directly to the Run method.
)

// PingCommand implements the Command interface to reply with "Pong!".
type PingCommand struct{}

// Prefix returns the command's prefix.
func (p *PingCommand) Prefix() string {
	return ".ping"
}

// Run executes the Ping command.
// It now correctly accepts the whatsmeow.Client instance.
func (p *PingCommand) Run(evt *events.Message, client *whatsmeow.Client) {
	text := "Pong!"
	fmt.Printf("Received .ping command from %s\n", evt.Info.Sender.String())

	// Use the 'client' parameter received in the Run method to send the message.
	_, err := client.SendMessage(context.Background(), evt.Info.Chat, &waProto.Message{
		Conversation: &text,
	})

	if err != nil {
		fmt.Printf("❌ Gagal mengirim pesan 'Pong!': %v\n", err)
	}
}
