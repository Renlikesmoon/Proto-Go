package commands

import (
    "strings"
    "context"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
    waProto "go.mau.fi/whatsmeow/binary/proto"
    "google.golang.org/protobuf/proto"
)

var cmds = []Command{
    &PingCommand{},
    &TimeCommand{},
    &HelpCommand{},
}

func HandleCommand(evt *events.Message, client *whatsmeow.Client) {
    msg := getMessageText(evt)
    if msg == "" {
        return
    }

    isReply := false
    var prefixUsed string

    if strings.HasPrefix(msg, "!") {
        prefixUsed = "!"
    } else if strings.HasPrefix(msg, ".") {
        prefixUsed = "."
        isReply = true
    } else {
        return
    }

    for _, cmd := range cmds {
        if strings.HasPrefix(msg, cmd.Prefix()) {
            if isReply {
                // Balas sebagai reply (bukan hanya kirim pesan)
                cmdWithReply := wrapWithReply(cmd)
                cmdWithReply.Run(evt, client)
            } else {
                cmd.Run(evt, client)
            }
            break
        }
    }
}

func getMessageText(evt *events.Message) string {
    if evt.Message.GetConversation() != "" {
        return evt.Message.GetConversation()
    }
    if ext := evt.Message.GetExtendedTextMessage(); ext != nil {
        return ext.GetText()
    }
    return ""
}

// Membungkus command agar kirim balasan sebagai reply
func wrapWithReply(cmd Command) Command {
    return &replyCommand{cmd}
}

type replyCommand struct {
    inner Command
}

func (r *replyCommand) Prefix() string {
    return r.inner.Prefix()
}

func (r *replyCommand) Run(evt *events.Message, client *whatsmeow.Client) {
    // Buat inner command mengirim balasan dengan quoted
    quotedMsg := &waProto.ContextInfo{
        StanzaId:   evt.Info.ID,
        Participant: proto.String(string(evt.Info.Sender.User) + "@s.whatsapp.net"),
        QuotedMessage: evt.Message,
    }

    // Gunakan contextInfo pada inner command
    switch c := r.inner.(type) {
    case *PingCommand:
        client.SendMessage(context.Background(), evt.Info.Chat, &waProto.Message{
            Conversation: proto.String("🏓 Pong (reply)!"),
            ContextInfo:  quotedMsg,
        })
    case *TimeCommand:
        client.SendMessage(context.Background(), evt.Info.Chat, &waProto.Message{
            Conversation: proto.String("⏰ Waktu server (reply)!"),
            ContextInfo:  quotedMsg,
        })
    case *HelpCommand:
        client.SendMessage(context.Background(), evt.Info.Chat, &waProto.Message{
            Conversation: proto.String("📚 Command list (reply)!"),
            ContextInfo:  quotedMsg,
        })
    default:
        // Fallback: jalankan normal
        r.inner.Run(evt, client)
    }
}
