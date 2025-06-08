package commands

import (
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/types/events"
	"github.com/Renlikesmoon/Proto-Go/lib" // Pastikan import path ini benar sesuai struktur proyek Anda
)

// commandList holds all the registered commands.
// It uses the Command interface defined in commands/command.go.
var commandList = []Command{
	&PingCommand{},
	&AnimeCommand{},
}

// HandleMessage processes incoming messages, checks for command prefixes,
// and dispatches to the appropriate command.
func HandleMessage(evt *events.Message) {
	// Use the GetMessageText helper function from the 'commands' package
	// (which should be defined and exported in commands/command.go).
	msg := GetMessageText(evt)

	// Check if the message starts with the configured command prefix.
	if !strings.HasPrefix(msg, config.CommandPrefix) {
		return // Not a command, ignore
	}

	// Check if the sender is the owner (if configured).
	// Ensure config.OwnerJID is a string representation of the JID.
	if string(evt.Info.Sender.String()) != config.OwnerJID {
		fmt.Printf("Pesan dari non-owner ditolak: %s\n", evt.Info.Sender.String())
		return // Block messages from non-owners
	}

	// Extract the command name from the message.
	// For example, if msg is "!ping", commandName will be "ping".
	commandText := strings.TrimPrefix(msg, config.CommandPrefix)
	parts := strings.Fields(commandText)
	if len(parts) == 0 {
		return // Only prefix, no command, ignore
	}
	commandName := parts[0]

	// Find and run the command.
	for _, cmd := range commandList {
		// Use cmd.Name() (assuming your Command interface has a Name() method now)
		// or cmd.Prefix() depending on how you want to match commands.
		// For simplicity, let's match based on Prefix() as in your original code.
		if commandName == strings.TrimPrefix(cmd.Prefix(), config.CommandPrefix) {
			// Crucially, pass the whatsmeow.Client instance (from lib.Client)
			// to the command's Run method.
			if lib.Client != nil {
				cmd.Run(evt, lib.Client) // Correctly passing whatsmeow.Client
			} else {
				fmt.Println("Error: whatsmeow.Client belum diinisialisasi.")
				// Optionally, send an error message to the chat if client is nil
			}
			break // Command found and executed, stop processing
		}
	}
}

// NOTE: The 'Command' interface and 'getMessageText' function
// are assumed to be defined and exported in 'commands/command.go'.
// They should NOT be present in this file (router.go) to avoid redeclaration errors.
//
// The 'Command' interface should look like this in 'commands/command.go':
// type Command interface {
// 	Prefix() string // Or Name() depending on your preference
// 	Run(evt *events.Message, client *whatsmeow.Client)
// }
//
// And GetMessageText should be:
// func GetMessageText(evt *events.Message) string { ... }
