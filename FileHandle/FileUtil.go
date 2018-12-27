package FileHandle

import (
	Errors "go-framework/errors"
	"os"
	"path/filepath"
)

// IsFileExsist 返回文件是否存在
func IsFileExsist(filePathStr string) (bool, error) {
	pathStr := GetFileAbsPath(filePathStr)
	_, err := os.Stat(pathStr)
	if err != nil && (os.IsNotExist(err) || os.IsPermission(err)) {
		return false, err
	}
	return true, nil
}

// GetFileAbsPath 获取文件的绝对路径
func GetFileAbsPath(filePathStr string) string {
	pathStr, err := filepath.Abs(filePathStr)
	if err != nil {
		Errors.Raise(err, Errors.IOErrCode)
	}
	return pathStr
}

// CreateLocalFile 创建文件，返回创建文件的指针
func CreateLocalFile(filePathStr string) *os.File {
	pathStr := GetFileAbsPath(filePathStr)
	handle, err := os.OpenFile(pathStr, os.O_CREATE, 0766)
	if err != nil {
		Errors.Raise(err, Errors.IOErrCode)
	}
	return handle
}

// OpenLocalFile 返回指定路径的文件的指针（名称和os.OpenFile区分开）
// 找不到文件时报系统错退出
func OpenLocalFile(filePathStr string) *os.File {
	if ok, err := IsFileExsist(filePathStr); !ok {
		Errors.Raise(err, Errors.IOErrCode)
	}
	pathStr := GetFileAbsPath(filePathStr)
	handle, err := os.OpenFile(pathStr, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		Errors.Raise(err, Errors.IOErrCode)
	}

	return handle
}

// CloseFile 关闭文件
func CloseFile(file *os.File) {
	file.Close()
}
