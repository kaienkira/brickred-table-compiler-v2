//go:build windows

package lib

import (
	"os"

	"golang.org/x/sys/windows"
)

func UtilOpenFileShared(filePath string) (*os.File, error) {
	filePathUtf16Ptr, err := windows.UTF16PtrFromString(filePath)
	if err != nil {
		return nil, err
	}

	fileHandle, err := windows.CreateFile(
		filePathUtf16Ptr,
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0)
	if err != nil {
		return nil, err
	}

	return os.NewFile(uintptr(fileHandle), filePath), nil
}
