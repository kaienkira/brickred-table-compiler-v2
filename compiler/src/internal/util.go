package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func UtilGetFullPath(filePath string) string {
	fullPath, err := filepath.Abs(filePath)
	if err != nil {
		return ""
	} else {
		return fullPath
	}
}

func UtilCheckFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}

	return true
}

func UtilCheckDirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	if info.IsDir() == false {
		return false
	}

	return true
}

func UtilWriteAllText(filePath string, fileContent string) bool {
	err := os.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"error: write file %s failed: %s",
			filePath, err.Error())
		return false
	}

	return true
}

func UtilCamelToUnderscore(camelName string) string {
	ret := g_camelToUnderscoreCase1Regexp.ReplaceAllString(camelName, "${1}_${2}")
	ret = g_camelToUnderscoreCase2Regexp.ReplaceAllString(ret, "${1}_${2}")
	ret = strings.ToLower(ret)

	return ret
}
