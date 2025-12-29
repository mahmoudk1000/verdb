package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	baseVerDBConfigDir = "verdb"
)

type ConfigBuilder[T any] struct {
	path           string
	configFileName string
	model          T
}

func NewConfigBuilder[T any](configFileName string, model T) *ConfigBuilder[T] {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil
	}

	path := filepath.Join(configDir, baseVerDBConfigDir, configFileName)

	cb := &ConfigBuilder[T]{
		path:           path,
		configFileName: configFileName,
		model:          model,
	}

	cb.initializeModel()

	return cb
}

func (cb *ConfigBuilder[T]) initializeModel() {
	if v, ok := any(&cb.model).(interface{ InitializeMap() }); ok {
		v.InitializeMap()
	}
}

func (cb *ConfigBuilder[T]) BuildConfigDir() error {
	if err := os.MkdirAll(filepath.Dir(cb.path), 0755); err != nil {
		return err
	}

	if f, err := os.Open(cb.path); err == nil {
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Printf("failed to close file: %v\n", err)
			}
		}()

		decoder := json.NewDecoder(f)
		if err := decoder.Decode(&cb.model); err != nil {
			fmt.Printf("warning: could not decode existing config: %v\n", err)
		}
	}

	cb.initializeModel()

	return nil
}

func (cb *ConfigBuilder[T]) Model() T {
	return cb.model
}


func (cb *ConfigBuilder[T]) Save() error {
	f, err := os.OpenFile(
		cb.path,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("failed to close file: %v\n", err)
		}
	}()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(cb.model); err != nil {
		return err
	}

	return nil
}
