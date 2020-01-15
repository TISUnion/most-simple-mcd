package utils

import (
	"github.com/mholt/archiver/v3"
	"os"
	"path/filepath"
)

// 判断所给路径文件/文件夹是否存在
func ExistsResource(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// 压缩多个文件
func CompressFiles(sources []string, destination string) error {
	return archiver.Archive(sources, destination)
}

// 压缩单个文件
func CompressFile(source, destination string) error {
	return archiver.Archive([]string{source}, destination)
}

// 创建文件
func CreateFile(path string) error {
	dirPath, _ := filepath.Split(path)
	if IsDir(dirPath) {
		if err := CreatDir(dirPath) ; err != nil {
			return err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}