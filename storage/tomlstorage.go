package storage

import "github.com/BurntSushi/toml"

type TomlStorage[T any] struct {
	storage StringStorage
}

func NewTomlStorage[T any](filePath string, opts ...Option) TomlStorage[T] {
	return TomlStorage[T]{
		storage: NewStringStorage(filePath, opts...),
	}
}

func (j TomlStorage[T]) Save(data T) error {
	tomlData, err := toml.Marshal(data)
	if err != nil {
		return err
	}
	return j.storage.Save(string(tomlData))
}

func (j TomlStorage[T]) Load() (T, error) {
	var data T
	tomlData, err := j.storage.Load()
	if err != nil {
		return data, err
	}
	err = toml.Unmarshal([]byte(tomlData), &data)
	return data, err
}

func (j TomlStorage[T]) TrySave(data T) error {
	tomlData, err := toml.Marshal(data)
	if err != nil {
		return err
	}
	return j.storage.TrySave(string(tomlData))
}

func (j TomlStorage[T]) TryLoad() (T, error) {
	var data T
	tomlData, err := j.storage.TryLoad()
	if err != nil {
		return data, err
	}
	err = toml.Unmarshal([]byte(tomlData), &data)
	return data, err
}
