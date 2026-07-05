package monitor 

type OffsetManager struct {
	offset  map[string]int64
}
// making a constructor for  offset to prevent nil map error 
func NewOffsetManager() *OffsetManager  {
	return &OffsetManager{
		offset: make(map[string]int64),
	}
}

func (of *OffsetManager) GetOffset(filename string) int64 {
	return of.offset[filename]
}
func (of *OffsetManager) UpdateOffset(fileName string, offset int64) {
	of.offset[fileName] = offset
}