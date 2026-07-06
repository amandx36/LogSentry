package monitor

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type OffsetManager struct {
	mu     sync.RWMutex
	offset map[string]int64
}

// Constructor
func NewOffsetManager() *OffsetManager {
	return &OffsetManager{
		offset: make(map[string]int64),
	}
}

// Read Offset
func (of *OffsetManager) GetOffset(filename string) int64 {
	of.mu.RLock()
	defer of.mu.RUnlock()

	return of.offset[filename]
}

// Update Offset
func (of *OffsetManager) UpdateOffset(filename string, offset int64) {
	of.mu.Lock()
	defer of.mu.Unlock()

	of.offset[filename] = offset
}

// Load offsets from JSON file
func LoadOffsets(path string) (*OffsetManager, error) {

	manager := NewOffsetManager()

	data, err := os.ReadFile(path)
	if err != nil {

		// File doesn't exist 	
		if errors.Is(err, os.ErrNotExist) {
			return manager, nil
		}

		return nil, err
	}

	if len(data) == 0 {
		return manager, nil
	}

	err = json.Unmarshal(data, &manager.offset)
	if err != nil {
		return nil, err
	}

	return manager, nil
}

// Save offsets to JSON file
func (of *OffsetManager) Save(path string) error {

	of.mu.RLock()
	defer of.mu.RUnlock()

	data, err := json.MarshalIndent(of.offset, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}