package commands

import (
    "context"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
    waProto "go.mau.fi/whatsmeow/binary/proto"
    "google.golang.org/protobuf/proto"
)

type HelpCommand struct{}

func (c *HelpCommand) Prefix() string {
    return "!help"
}

func (c *HelpCommand) Run(evt *events.Message, client *whatsmeow.Client) {
    text := "📖 Commands available:\n!ping - Cek koneksi\n!time - Waktu server\n!help - Lihat bantuan"
    client.SendMessage(context.Background(), evt.Info.Chat, &waProto.Message{
        Conversation: proto.String(text),
    })
}
