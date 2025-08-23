package storage

import "encoding/json"

type JsonStorage[T any] struct {
	storage StringStorage
}

// New is alias for NewJsonStorage
func New[T any](filePath string, opts ...Option) JsonStorage[T] {
	return NewJsonStorage[T](filePath, opts...)
}

func NewJsonStorage[T any](filePath string, opts ...Option) JsonStorage[T] {
	return JsonStorage[T]{
		storage: NewStringStorage(filePath, opts...),
	}
}

func (j JsonStorage[T]) Save(data T) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return j.storage.Save(string(jsonData))
}

func (j JsonStorage[T]) Load() (T, error) {
	var data T
	jsonData, err := j.storage.Load()
	if err != nil {
		return data, err
	}
	err = json.Unmarshal([]byte(jsonData), &data)
	return data, err
}

func (j JsonStorage[T]) TrySave(data T) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return j.storage.TrySave(string(jsonData))
}

func (j JsonStorage[T]) TryLoad() (T, error) {
	var data T
	jsonData, err := j.storage.TryLoad()
	if err != nil {
		return data, err
	}
	err = json.Unmarshal([]byte(jsonData), &data)
	return data, err
}
