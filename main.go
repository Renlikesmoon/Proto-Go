package main

import (
	"context"
	"fmt"
	"os"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.mau.fi/whatsmeow/types/events"

	"github.com/mdp/qrterminal/v3"
)

var Client *whatsmeow.Client

// StartClient menginisialisasi klien Whatsmeow dan menangani autentikasi.
// Nomor telepon untuk pairing sekarang langsung ditetapkan di dalam fungsi ini.
func StartClient() error { // Parameter 'phoneNumber' telah dihapus
	// --- Nomor telepon langsung ditetapkan di sini ---
	phoneNumber := "6285954540177" // Nomor hardcode langsung di dalam fungsi
	// --- Akhir hardcode ---

	ctx := context.Background()

	// Gunakan logger kosong agar tidak spam log
	dbLog := waLog.Noop

	// Buat koneksi database session SQLite
	container, err := sqlstore.New(ctx, "sqlite", "file:session.db?_foreign_keys=on", dbLog)
	if err != nil {
		return fmt.Errorf("❌ Gagal konek database: %w", err)
	}

	// Ambil device pertama dari store
	device, err := container.GetFirstDevice(ctx)
	if err != nil {
		return fmt.Errorf("❌ Gagal ambil device: %w", err)
	}

	// Inisialisasi client Whatsmeow dengan device dan logger
	Client = whatsmeow.NewClient(device, dbLog)

	// Event handler untuk menangani event masuk
	Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Println("📩 Pesan masuk dari:", v.Info.Sender.String())
			// Anda mungkin ingin merutekan ini ke commands.HandleMessage di sini
			// Misalnya: commands.HandleMessage(v) (pastikan untuk mengimpor "wa_bot/commands")
		case *events.Disconnected:
			fmt.Println("🔌 Terputus dari WhatsApp")
		}
	})

	// Jika belum ada session ID, mulai proses pairing dengan QR
	if Client.Store.ID == nil {
		resp, err := Client.PairPhone(
			ctx,
			phoneNumber,             // Sekarang menggunakan variabel lokal 'phoneNumber'
			false,                   // showPushNotification
			whatsmeow.PairClientChrome, // Tipe klien
			"Go Bot (Desktop)",      // Nama tampilan untuk klien
		)
		if err != nil {
			return fmt.Errorf("❌ Gagal pairing: %w", err)
		}

		// Tampilkan QR code di terminal untuk discan
		qrterminal.GenerateHalfBlock(resp, qrterminal.L, os.Stdout)
		fmt.Println("✅ Scan QR di atas dengan WhatsApp kamu.")
		fmt.Printf("Kode Pairing: %s\n", resp)
	} else {
		// Jika sudah ada session, langsung konek
		err = Client.Connect()
		if err != nil {
			return fmt.Errorf("❌ Gagal konek ke WhatsApp: %w", err)
		}
		fmt.Println("✅ Terhubung ke WhatsApp sebagai", Client.Store.ID.User)
	}
	return nil
}
