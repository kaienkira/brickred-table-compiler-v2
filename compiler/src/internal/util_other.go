//go:build !windows

package lib

import (
	"os"
)

func UtilOpenFileShared(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, os.O_RDONLY, 0)
}
