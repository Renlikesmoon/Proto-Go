package commands

import (
	"fmt"
	"strings"

	"github.com/Renlikesmoon/Proto-Go/config" // Import the config package for owner JID and command prefix
	"github.com/Renlikesmoon/Proto-Go/lib"    // Import lib to access the global whatsmeow.Client instance

	"go.mau.fi/whatsmeow/types/events"
)

// commandList holds all the registered commands.
// Each command must implement the 'Command' interface defined in commands/command.go.
var commandList = []Command{
	&PingCommand{},  // Make sure ping.go exists and implements Command
	&AnimeCommand{}, // Make sure anime.go exists and implements Command
	&TimeCommand{},  // Make sure time.go exists and implements Command
	&HelpCommand{},  // Make sure help.go exists and implements Command
	// Add any other new commands here after creating their respective files
}

// HandleMessage processes incoming WhatsApp messages.
// It checks if the message is a command, if the sender is authorized,
// and then dispatches the command to the appropriate handler.
func HandleMessage(evt *events.Message) {
	// Get the message text using the shared utility function from commands/command.go.
	// This assumes GetMessageText is defined and exported in command.go (e.g., func GetMessageText(...) string).
	msg := GetMessageText(evt)

	// Check if the message starts with the configured command prefix.
	if !strings.HasPrefix(msg, config.CommandPrefix) {
		return // Not a command, ignore
	}

	// Check if the sender is the owner as defined in config.
	// If not the owner, simply return without processing or replying.
	// Ensure config.OwnerJID is a string representation of the JID (e.g., "6281234567890@s.whatsapp.net").
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
		// We trim the prefix from cmd.Prefix() because commandName already has the main prefix removed.
		if commandName == strings.TrimPrefix(cmd.Prefix(), config.CommandPrefix) {
			// Ensure the Whatsmeow client is initialized before attempting to use it.
			if lib.Client != nil {
				// Execute the command's Run method, passing the event and the Whatsmeow client.
				cmd.Run(evt, lib.Client)
			} else {
				fmt.Println("Error: whatsmeow.Client belum diinisialisasi. Tidak bisa menjalankan perintah.")
				// Optionally, send an error message to the chat if the client is not available.
				// You would need a separate helper function here if you want to reply.
			}
			return // Command handled, exit function
		}
	}

	// If no command matched after checking all registered commands.
	fmt.Printf("Perintah tidak dikenal dari %s: %s\n", evt.Info.Sender.String(), msg)
	// Optionally, send a "command not found" message to the user:
	// if lib.Client != nil {
	//     // You might need a helper function here, similar to the one we defined in anime.go
	//     // func reply(evt *events.Message, text string, client *whatsmeow.Client) { ... }
	//     // reply(evt, "Perintah tidak dikenal. Ketik !help untuk melihat daftar perintah.", lib.Client)
	// }
}

// NOTE: The 'Command' interface and 'GetMessageText' function
// must be defined and exported ONLY in 'commands/command.go'.
// They should NOT be present in this file (router.go) to avoid redeclaration errors.
//
// Example of how they should look in 'commands/command.go':
// type Command interface {
//     Prefix() string
//     Run(evt *events.Message, client *whatsmeow.Client)
// }
//
// func GetMessageText(evt *events.Message) string {
//     msg := evt.Message.GetConversation()
//     if msg == "" && evt.Message.ExtendedTextMessage != nil {
//         msg = evt.Message.ExtendedTextMessage.GetText()
//     }
//     return msg
// }
