package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	// Ensure this import path is correct for your project
	// It's assumed that lib.Client is available via the router passing it.
	// We'll use the 'client' argument directly in the Run method now.
	"go.mau.fi/whatsmeow/types/events"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow" // Added this import for whatsmeow.Client type
)

// AnimeCommand implements the Command interface.
type AnimeCommand struct{}

// Prefix returns the command's prefix.
func (a *AnimeCommand) Prefix() string {
	return ".anime"
}

// Run executes the anime search command.
// The signature now includes the *whatsmeow.Client parameter.
func (a *AnimeCommand) Run(evt *events.Message, client *whatsmeow.Client) {
	// Use the shared GetMessageText function from the commands package.
	// Assuming GetMessageText is now defined in commands/command.go and is exported.
	msg := GetMessageText(evt)

	args := strings.TrimSpace(strings.TrimPrefix(msg, a.Prefix()))
	if args == "" {
		reply(evt, "Usage: .anime <query>", client) // Pass client to reply
		return
	}

	url := fmt.Sprintf("https://api.jikan.moe/v4/anime?q=%s&limit=1", args)
	resp, err := http.Get(url)
	if err != nil {
		reply(evt, "❌ Error fetching anime.", client) // Pass client to reply
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			Title    string `json:"title"`
			Synopsis string `json:"synopsis"`
			URL      string `json:"url"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil || len(result.Data) == 0 {
		reply(evt, "❌ Anime not found.", client) // Pass client to reply
		return
	}

	anime := result.Data[0]
	text := fmt.Sprintf("📺 *%s*\n\n📝 %s\n🔗 %s", anime.Title, anime.Synopsis, anime.URL)
	reply(evt, text, client) // Pass client to reply
}

// NOTE: getMessageText should be moved to commands/command.go or a shared utility file
// and removed from here to resolve the redeclaration error.
// The version in commands/command.go should be exported (e.g., GetMessageText).
// I am assuming it has been moved and is imported/accessible.


// reply sends a message back to the chat where the event originated.
// It now explicitly takes the *whatsmeow.Client as an argument.
func reply(evt *events.Message, text string, client *whatsmeow.Client) {
	if client == nil {
		fmt.Println("Error: whatsmeow.Client is nil in reply function.")
		return
	}
	_, err := client.SendMessage(context.Background(), evt.Info.Chat, &waProto.Message{
		Conversation: &text,
	})
	if err != nil {
		fmt.Printf("❌ Gagal mengirim pesan: %v\n", err)
	}
}
