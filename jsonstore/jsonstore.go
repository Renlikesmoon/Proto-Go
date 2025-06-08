package jsonstore

import (
	"encoding/json"
	"fmt"
	"os"

	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/legacy"
	"go.mau.fi/whatsmeow/store/mem"
	waTypes "go.mau.fi/whatsmeow/types"
)

type JSONSessionStore struct {
	filePath string
	data     *store.Device
	memStore *mem.Store
}

func NewJSONStore(filePath string) (*JSONSessionStore, error) {
	memStore := mem.New()
	js := &JSONSessionStore{
		filePath: filePath,
		memStore: memStore,
	}

	err := js.load()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load session error: %w", err)
	}
	return js, nil
}

// Memuat device dari file JSON
func (s *JSONSessionStore) load() error {
	file, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var session legacy.Device
	err = json.NewDecoder(file).Decode(&session)
	if err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	device := store.NewDevice(s.memStore, waTypes.DeviceID{})
	err = device.Restore(&session)
	if err != nil {
		return fmt.Errorf("restore error: %w", err)
	}

	s.data = device
	return nil
}

// Simpan session ke file
func (s *JSONSessionStore) Save() error {
	if s.data == nil {
		return fmt.Errorf("device is nil")
	}

	session := s.data.Serialize()
	file, err := os.Create(s.filePath)
	if err != nil {
		return fmt.Errorf("file create error: %w", err)
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(session)
}

func (s *JSONSessionStore) GetDevice() *store.Device {
	return s.data
}

func (s *JSONSessionStore) SetDevice(device *store.Device) {
	s.data = device
}
