package jsonstore

import (
	"os"
)

type JSONStore struct {
	filePath string
}

func NewJSONStore(filePath string) *JSONStore {
	return &JSONStore{filePath}
}

// Save menerima data session []byte hasil client.Store.Serialize()
func (js *JSONStore) Save(data []byte) error {
	return os.WriteFile(js.filePath, data, 0644)
}

// Load membaca file dan mengembalikan []byte session untuk Restore()
func (js *JSONStore) Load() ([]byte, error) {
	data, err := os.ReadFile(js.filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}
