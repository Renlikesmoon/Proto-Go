package commands

import (
    "context"
    "fmt"
    "time"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
    waProto "go.mau.fi/whatsmeow/binary/proto"
    "google.golang.org/protobuf/proto"
)

type TimeCommand struct{}

func (c *TimeCommand) Prefix() string {
    return "!time"
}

func (c *TimeCommand) Run(evt *events.Message, client *whatsmeow.Client) {
    now := time.Now().Format("2006-01-02 15:04:05")
    msg := fmt.Sprintf("🕒 Current time: %s", now)

    client.SendMessage(context.Background(), evt.Info.Chat, &waProto.Message{
        Conversation: proto.String(msg),
    })
}
