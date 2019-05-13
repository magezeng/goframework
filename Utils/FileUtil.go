package Utils

import (
	"path"
	"strings"
)

func GetFileNameWithoutExt(pathStr string) string {
	filenameWithSuffix := path.Base(pathStr)
	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix)
	fileName := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	return fileName
}
