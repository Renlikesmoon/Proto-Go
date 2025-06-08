package commands

import (
	"context"
	"fmt" // Required for fmt.Printf
	"strings"

	"github.com/Renlikesmoon/Proto-Go/config" // Import the config package for owner JID and command prefix
	"github.com/Renlikesmoon/Proto-Go/lib"    // Import lib to access the global whatsmeow.Client instance

	"go.mau.fi/whatsmeow"       // Required for whatsmeow.Client type
	"go.mau.fi/whatsmeow/types/events"
	waProto "go.mau.fi/whatsmeow/binary/proto" // Required for waProto.Message
	// If your AnimeCommand uses a separate 'anime' package, import it here:
	// "whatsbot/anime"
)

// --- IMPORTANT NOTE ---
// The 'Command' interface and 'GetMessageText' function MUST be defined and exported
// ONLY in 'commands/command.go' to avoid redeclaration errors.
// This file assumes they are correctly defined there.

// commandList holds all the registered commands.
// Each command must implement the 'Command' interface defined in commands/command.go.
var commandList = []Command{
	&PingCommand{},
	&AnimeCommand{},
	&TimeCommand{}, // Make sure you have this command file implemented
	&HelpCommand{},  // Make sure you have this command file implemented
	// Add any other new commands here
}

// HandleMessage processes incoming WhatsApp messages.
// It checks if the message is a command, if the sender is authorized,
// and then dispatches the command to the appropriate handler.
func HandleMessage(evt *events.Message) {
	// Get the message text using the shared utility function from commands/command.go.
	msg := GetMessageText(evt)

	// Check if the message starts with the configured command prefix.
	// Use config.CommandPrefix for consistency.
	if !strings.HasPrefix(msg, config.CommandPrefix) {
		return // Not a command, ignore
	}

	// Check if the sender is the owner as defined in config.
	// If not the owner, simply return without processing or replying.
	if string(evt.Info.Sender.String()) != config.OwnerJID {
		fmt.Printf("Pesan dari non-owner ditolak: %s (Konten: %s)\n", evt.Info.Sender.String(), msg)
		return
	}

	// Extract the command name from the message.
	// Example: If message is "!ping some_arg", commandName will be "ping".
	commandText := strings.TrimPrefix(msg, config.CommandPrefix)
	parts := strings.Fields(commandText)
	if len(parts) == 0 {
		return // Message was just the prefix, no actual command, ignore
	}
	commandName := parts[0] // The first word after the prefix is the command name

	// Iterate through the registered commands to find a match.
	for _, cmd := range commandList {
		// Compare the extracted command name with the command's registered prefix (without its own prefix).
		if commandName == strings.TrimPrefix(cmd.Prefix(), config.CommandPrefix) {
			// Ensure the Whatsmeow client is initialized before attempting to use it.
			if lib.Client != nil {
				// Execute the command's Run method, passing the event and the Whatsmeow client.
				cmd.Run(evt, lib.Client) // This is the crucial part for correct Command interface usage.
			} else {
				fmt.Println("Error: whatsmeow.Client belum diinisialisasi. Tidak bisa menjalankan perintah.")
				// Optionally, send an error message to the chat here if the client is not available.
				// sendText(evt, "Bot sedang tidak aktif. Coba lagi nanti.", nil) // Would require client argument for sendText
			}
			return // Command handled, exit function
		}
	}

	// If no command matched after checking all registered commands.
	fmt.Printf("Perintah tidak dikenal dari %s: %s\n", evt.Info.Sender.String(), msg)
	// You could optionally send a "command not found" message to the user:
	// if lib.Client != nil {
	//     sendText(evt, "Perintah tidak dikenal. Ketik !help untuk melihat daftar perintah.", lib.Client)
	// }
}

// sendText is a helper function to send a text message back to the chat.
// It now takes the whatsmeow.Client instance as an argument for consistency.
func sendText(evt *events.Message, text string, client *whatsmeow.Client) {
	if client == nil {
		fmt.Println("Error: Tidak bisa mengirim pesan, whatsmeow.Client nil di sendText.")
		return
	}
	_, err := client.SendMessage(context.Background(), evt.Info.Chat, &waProto.Message{
		Conversation: &text, // Correctly pass a pointer to the string
	})
	if err != nil {
		fmt.Printf("❌ Gagal mengirim pesan: %v\n", err)
	}
}
