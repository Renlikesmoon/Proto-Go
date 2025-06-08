package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Renlikesmoon/Proto-Go/lib"
	"go.mau.fi/whatsmeow/types/events"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

type AnimeCommand struct{}

func (a *AnimeCommand) Prefix() string {
	return ".anime"
}

func (a *AnimeCommand) Run(evt *events.Message) {
	msg := getMessageText(evt)
	args := strings.TrimSpace(strings.TrimPrefix(msg, a.Prefix()))
	if args == "" {
		reply(evt, "Usage: .anime <query>")
		return
	}

	url := fmt.Sprintf("https://api.jikan.moe/v4/anime?q=%s&limit=1", args)
	resp, err := http.Get(url)
	if err != nil {
		reply(evt, "❌ Error fetching anime.")
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
		reply(evt, "❌ Anime not found.")
		return
	}

	anime := result.Data[0]
	text := fmt.Sprintf("📺 *%s*\n\n📝 %s\n🔗 %s", anime.Title, anime.Synopsis, anime.URL)
	reply(evt, text)
}

func getMessageText(evt *events.Message) string {
	msg := evt.Message.GetConversation()
	if msg == "" && evt.Message.ExtendedTextMessage != nil {
		msg = evt.Message.ExtendedTextMessage.GetText()
	}
	return msg
}

func reply(evt *events.Message, text string) {
	lib.Client.SendMessage(context.Background(), evt.Info.Chat, &waProto.Message{
		Conversation: &text,
	})
}
