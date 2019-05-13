package Utils

import (
	"archive/zip"
	"bytes"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

// UnzipToCurrentPath 解压缩到压缩文件相同目录
func UnzipToCurrentPath(path string) error {
	return unzip(path, filepath.Dir(path))
}

// unzip 解压缩方法，需要zip文件地址及解压缩目标地址
func unzip(path, target string) error {
	if !isZip(path) {
		return errors.New("不是一个有效的zip压缩文件！")
	}
	// TODO: 如果有安装unzip，7zip等，直接解压缩
	// 否则用原生解压缩方式
	reader, err := zip.OpenReader(path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

// 判断这个文件是不是一个合法的zip文件
func isZip(zipPath string) bool {
	f, err := os.Open(zipPath)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 4)
	if n, err := f.Read(buf); err != nil || n < 4 {
		return false
	}

	return bytes.Equal(buf, []byte("PK\x03\x04"))
}
