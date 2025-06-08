package commands

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// Command defines the interface that all commands must implement.
// It includes the Prefix to identify the command and the Run method
// that executes the command logic, with access to the Whatsmeow client.
type Command interface {
	Prefix() string // Example: ".ping", ".anime"
	// Run executes the command logic. It receives the incoming message event
	// and the whatsmeow.Client instance to perform actions like sending replies.
	Run(evt *events.Message, client *whatsmeow.Client)
}

// GetMessageText extracts the text content from different message types.
// This function is placed here as a shared utility to avoid redeclaration
// in multiple command files or the router.
func GetMessageText(evt *events.Message) string {
	msg := evt.Message.GetConversation()
	if msg == "" && evt.Message.ExtendedTextMessage != nil {
		msg = evt.Message.ExtendedTextMessage.GetText()
	}
	// Add other message types if necessary (e.g., image captions, document captions)
	// if msg == "" && evt.Message.ImageMessage != nil {
	//     msg = evt.Message.ImageMessage.GetCaption()
	// }
	return msg
}
