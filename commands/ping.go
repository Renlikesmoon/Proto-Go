package commands

import (
    "context"
    "log"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
    waProto "go.mau.fi/whatsmeow/binary/proto"

    "google.golang.org/protobuf/proto"
)

type PingCommand struct{}

func (c *PingCommand) Prefix() string {
    return "!ping"
}

func (c *PingCommand) Run(evt *events.Message, client *whatsmeow.Client) {
    reply := &waProto.Message{Conversation: proto.String("🏓 Pong!")}
    _, err := client.SendMessage(context.Background(), evt.Info.Chat, reply)
    if err != nil {
        log.Println("❌ Failed to reply:", err)
    }
}
