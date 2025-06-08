package commands

import (
    "strings"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

var cmds = []Command{
    &PingCommand{},
    &TimeCommand{},
    &HelpCommand{},
}

func HandleCommand(evt *events.Message, client *whatsmeow.Client) {
    msg := getMessageText(evt)
    if msg == "" || !strings.HasPrefix(msg, "!") {
        return
    }

    for _, cmd := range cmds {
        if strings.HasPrefix(msg, cmd.Prefix()) {
            cmd.Run(evt, client)
            break
        }
    }
}

func getMessageText(evt *events.Message) string {
    if evt.Message.Conversation != "" {
        return evt.Message.GetConversation()
    }
    if ext := evt.Message.GetExtendedTextMessage(); ext != nil {
        return ext.GetText()
    }
    return ""
}
