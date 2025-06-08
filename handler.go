package main

import (
    "mybot/commands" // Ganti dengan path modulmu

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

func messageHandler(evt interface{}, client *whatsmeow.Client) {
    switch v := evt.(type) {
    case *events.Message:
        commands.HandleCommand(v, client)
    }
}
