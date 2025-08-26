package storage

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
)

func NewStringStorage(filePath string, opts ...Option) Storage[string] {
	return New(filePath,
		func(content string) ([]byte, error) { return []byte(content), nil },
		func(data []byte) (string, error) { return string(data), nil },
		opts...,
	)
}

func NewJSONStorage[U any](filePath string, opts ...Option) Storage[U] {
	return New(filePath,
		func(content U) ([]byte, error) {
			return json.Marshal(content)
		},
		func(data []byte) (U, error) {
			var content U
			err := json.Unmarshal(data, &content)
			return content, err
		},
		opts...,
	)
}

func NewTOMLStorage[U any](filePath string, opts ...Option) Storage[U] {
	return New(filePath,
		func(content U) ([]byte, error) {
			return toml.Marshal(content)
		},
		func(data []byte) (U, error) {
			var content U
			err := toml.Unmarshal(data, &content)
			return content, err
		},
		opts...,
	)
}
