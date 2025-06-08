package main

import (
	"log"
	"os"
	"whatsbot/lib"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("❗ Masukkan nomor HP untuk pairing (contoh: go run main.go 6281234567890)")
	}
	phone := os.Args[1]
	err := lib.StartWhatsAppWithCode(phone)
	if err != nil {
		log.Fatal("❌ Error:", err)
	}
}
