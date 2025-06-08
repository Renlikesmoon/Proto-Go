package commands

import (
	"context"
	"whatsbot/lib"

	"go.mau.fi/whatsmeow/types/events"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

type PingCommand struct{}

func (p *PingCommand) Prefix() string {
	return ".ping"
}

func (p *PingCommand) Run(evt *events.Message) {
	text := "Pong!"
	lib.Client.SendMessage(context.Background(), evt.Info.Chat, &waProto.Message{
		Conversation: &text,
	})
}
