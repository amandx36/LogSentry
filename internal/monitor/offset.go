package monitor 

import (
	"sync"
)
type OffsetManager struct {
	mu sync.RWMutex
	offset  map[string]int64
}
// making a constructor for  offset to prevent nil map error 
func NewOffsetManager() *OffsetManager  {
	// adding mutex 
	
	return &OffsetManager{
		offset: make(map[string]int64),
	}
}

// read lock 
func (of *OffsetManager) GetOffset(filename string) int64 {
    of.mu.RLock()     
    defer of.mu.RUnlock()
    return of.offset[filename]
}
// write lock 
func (of *OffsetManager) UpdateOffset(fileName string, offset int64) {
	// write lock 
	of.mu.Lock()
	defer of.mu.Unlock()	
	of.offset[fileName] = offset
}