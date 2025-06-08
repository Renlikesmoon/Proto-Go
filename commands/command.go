package commands

import (
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

type Command interface {
    Prefix() string
    Run(evt *events.Message, client *whatsmeow.Client)
}
