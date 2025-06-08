package jsonstore

import (
	"encoding/json"
	"fmt"
	"os"

	"go.mau.fi/whatsmeow/store"
)

type JSONStore struct {
	filePath string
	store    *store.Store
}

func NewJSONStore(filePath string) (*JSONStore, error) {
	js := &JSONStore{
		filePath: filePath,
	}

	err := js.Load()
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return js, nil
}

func (js *JSONStore) Load() error {
	file, err := os.Open(js.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	s := store.NewStore()
	err = json.NewDecoder(file).Decode(&s.Data)
	if err != nil {
		return fmt.Errorf("decode session json error: %w", err)
	}

	js.store = s
	return nil
}

func (js *JSONStore) Save() error {
	if js.store == nil {
		return fmt.Errorf("store is nil")
	}
	file, err := os.Create(js.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(js.store.Data)
}

func (js *JSONStore) GetStore() *store.Store {
	if js.store == nil {
		js.store = store.NewStore()
	}
	return js.store
}
