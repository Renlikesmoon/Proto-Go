package commands

import (
	"strings"
	"whatsbot/config"

	"go.mau.fi/whatsmeow/types/events"
)

var commandList = []Command{
	&PingCommand{},
	&AnimeCommand{},
}

type Command interface {
	Prefix() string
	Run(evt *events.Message)
}

func HandleMessage(evt *events.Message) {
	msg := getMessageText(evt)

	if !strings.HasPrefix(msg, config.CommandPrefix) {
		return
	}

	// Cek sender owner
	if string(evt.Info.Sender.String()) != config.OwnerJID {
		return // langsung stop, hanya owner boleh pakai
	}

	for _, cmd := range commandList {
		if strings.HasPrefix(msg, cmd.Prefix()) {
			cmd.Run(evt)
			break
		}
	}
}

func getMessageText(evt *events.Message) string {
	msg := evt.Message.GetConversation()
	if msg == "" && evt.Message.ExtendedTextMessage != nil {
		msg = evt.Message.ExtendedTextMessage.GetText()
	}
	return msg
}
