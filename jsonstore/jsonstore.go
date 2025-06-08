package jsonstore

import (
	"encoding/json"
	"os"
	"sync"

	"go.mau.fi/whatsmeow/store"
	waTypes "go.mau.fi/whatsmeow/types"
)

type JSONStore struct {
	mu     sync.Mutex
	Store  *store.Device
	Path   string
	Client *waTypes.DeviceID
}

func NewJSONStore(path string) (*JSONStore, error) {
	js := &JSONStore{Path: path}
	data, err := os.ReadFile(path)
	if err == nil {
		var dev store.Device
		err := json.Unmarshal(data, &dev)
		if err != nil {
			return nil, err
		}
		js.Store = &dev
	}
	if js.Store == nil {
		js.Store = &store.Device{}
	}
	return js, nil
}

func (j *JSONStore) Save() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	data, err := json.MarshalIndent(j.Store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(j.Path, data, 0600)
}

func (j *JSONStore) GetDevice() *store.Device {
	return j.Store
}
