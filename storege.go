package main

import (
	"encoding/json"
	"os"
)

type Storege[T any] struct {
	FileName string
}

func NewStorage[T any](fileName string) *Storege[T] {
	return &Storege[T]{FileName: fileName}
}

func (s *Storege[T]) Save(data T) error {
	fileData, err := json.MarshalIndent(data, "", "		")

	if err != nil{
		return err
	}

	return os.WriteFile(s.FileName, fileData, 0644)
}

func (s *Storege[T]) Load (data *T) error{
	fileData, err := os.ReadFile((s.FileName))

	if err != nil{
		return err
	}

	return json.Unmarshal(fileData, data)
}