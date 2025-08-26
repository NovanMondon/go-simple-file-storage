package storage

import "os"

type OpenedStorage[T any] struct {
	filePath string

	marshal   func(T) ([]byte, error)
	unmarshal func([]byte) (T, error)
}

func (o OpenedStorage[T]) Read() (T, error) {
	var none T

	data, err := os.ReadFile(o.filePath)
	if err != nil {
		return none, err
	}

	result, err := o.unmarshal(data)
	if err != nil {
		return none, err
	}

	return result, nil
}

func (o OpenedStorage[T]) Write(content T) error {
	data, err := o.marshal(content)
	if err != nil {
		return err
	}

	err = os.WriteFile(o.filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
