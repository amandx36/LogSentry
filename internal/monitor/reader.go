package monitor

import (
	"os"
	"io"
)

// ReadNewLogs(file,lastOffset) -> os.open -> seek(lastOffset) -> read remaining bytes -> current file size -> return data,newOffset

// offset = how many bytes to move
// whence = from where to start counting   [ o.SeekStart (0) , o.SeekCurrent (1), o.SeekEnd (2) ]
func ReadNewLogs(fileName string, lastOffset int64) ([]byte, int64, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	_, err = file.Seek(lastOffset, 0)
	if err != nil {
		return nil, 0, err
	}
	// io.ReadAll() reads from the file's current cursor position to the end of the file
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, 0, err
	}

	newOffset, err := file.Seek(0, 2)
	if err != nil {
		return nil, 0, err
	}

	return data, newOffset, nil
}