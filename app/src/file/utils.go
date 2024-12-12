package file

import (
	"github.com/codecrafters-io/grep-starter-go/src/logs"
	"os"
)

func CloseFile(f *os.File) {
	err := f.Close()
	if err != nil {
		logs.Error("error closing file: %v", err)
	}
}
