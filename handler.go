package commands

import (
	"strings"
	"whatsbot/config"
	"whatsbot/lib"

	"go.mau.fi/whatsmeow/types/events"
)

func HandleMessage(evt *events.Message) {
	msg := getMessageText(evt)

	if !strings.HasPrefix(msg, ".") {
		return // bukan perintah
	}

	// Cek sender, jika bukan owner, langsung return tanpa balas
	if string(evt.Info.Sender.String()) != config.OwnerJID {
		return // langsung stop di sini, tidak reply
	}

	// proses perintah owner
	switch {
	case strings.HasPrefix(msg, ".ping"):
		sendText(evt, "Pong!")
	case strings.HasPrefix(msg, ".anime"):
		// panggil fungsi anime search misal
		// ...
	default:
		// perintah lain
	}
}

func getMessageText(evt *events.Message) string {
	msg := evt.Message.GetConversation()
	if msg == "" && evt.Message.ExtendedTextMessage != nil {
		msg = evt.Message.ExtendedTextMessage.GetText()
	}
	return msg
}

func sendText(evt *events.Message, text string) {
	lib.Client.SendMessage(evt.Info.Context, evt.Info.Chat, &waProto.Message{
		Conversation: &text,
	})
}
