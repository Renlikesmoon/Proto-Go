package main

import (
	"log"
	// os library tidak diperlukan jika tidak menggunakan os.Args lagi
	// "os" 

	// Pastikan jalur impor ini benar untuk pengaturan modul Anda.
	// Ganti "github.com/Renlikesmoon/Proto-Go" dengan nama modul Anda jika berbeda.
	"github.com/Renlikesmoon/Proto-Go/commands" // Untuk commands.HandleMessage
	"github.com/Renlikesmoon/Proto-Go/lib"      // Untuk lib.StartClient dan lib.Client

	"go.mau.fi/whatsmeow/types/events" // Untuk event whatsmeow
)

// main adalah titik masuk aplikasi Go.
func main() {
	// --- Nomor telepon langsung ditetapkan di sini ---
	// Ganti dengan nomor yang Anda inginkan.
	phone := "6285954540177" // Nomor telepon yang Anda berikan
	// --- Akhir hardcode ---

	// Mulai klien Whatsmeow. Fungsi ini menangani koneksi,
	// pairing (jika diperlukan), dan mengembalikan kesalahan apa pun.
	// Panggilan ini akan memanggil fungsi StartClient yang ada di paket 'lib'.
	err := lib.StartClient(phone)
	if err != nil {
		log.Fatalf("Error saat memulai klien WhatsApp: %v", err)
	}

	// Tambahkan event handler ke klien Whatsmeow untuk memproses event yang masuk.
	// Ini berjalan dalam goroutine terpisah yang dikelola oleh whatsmeow itu sendiri.
	lib.Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			// Abaikan pesan yang dikirim oleh bot itu sendiri untuk mencegah loop tak terbatas.
			if v.Info.MessageSource.IsFromMe {
				return
			}
			// Rute pesan masuk ke handler perintah.
			commands.HandleMessage(v)
		case *events.Connected:
			// Catat ketika klien WhatsApp berhasil terhubung.
			log.Println("Klien WhatsApp terhubung!")
		case *events.Disconnected:
			// Catat ketika klien WhatsApp terputus.
			// Pustaka Whatsmeow biasanya menangani penyambungan ulang otomatis.
			log.Println("Klien WhatsApp terputus.")
		case *events.QR:
			// lib.StartClient Anda sudah menangani tampilan kode QR.
			// Kita hanya akan mencatat bahwa event QR terjadi di sini tanpa mencoba mengakses v.QRCode.
			log.Println("Menerima event kode QR (QR ditampilkan oleh lib.StartClient).")
		case *events.PairingCode:
			// Serupa dengan QR, event PairingCode ditangani oleh PairPhone di StartClient.
			// Kita hanya akan mencatat bahwa event Pairing Code terjadi di sini tanpa mencoba mengakses v.Code.
			log.Println("Menerima event kode Pairing (kode ditampilkan oleh lib.StartClient).")
		}
	})

	// Ini akan memblokir goroutine utama tanpa batas waktu, menjaga aplikasi tetap berjalan.
	// Whatsmeow beroperasi dalam goroutine sendiri untuk jaringan dan penanganan event.
	select {}
}
