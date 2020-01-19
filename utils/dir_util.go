package utils

import (
	"errors"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/mholt/archiver/v3"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 递归创建文件夹
func CreatDir(dirPath string) error {
	if !ExistsResource(dirPath) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
		if err := os.Chmod(dirPath, 0777); err != nil {
			return err
		}
	}
	return nil
}

// 复制文件夹对外暴露方法
func CopyDir(src, dest string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return copyDir(src, dest, info)
}

// 解压到文件夹
func UnCompressDir(source, destination string) error {
	return archiver.Unarchive(source, destination)
}

// 获取当前执行文件目录
func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, "\\", "/", -1)
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		return "", errors.New(constant.GET_CURRENT_PATH_FAILED)
	}
	return string(path[0 : i+1]), nil
}

/* --------------内部方法---------------*/
// 复制文件夹主方法
func copyDir(src, dest string, info os.FileInfo) error {
	if info.Mode()&os.ModeSymlink != 0 {
		return linkCopy(src, dest, info)
	}
	if info.IsDir() {
		return dirCopy(src, dest, info)
	}
	return fileCopy(src, dest, info)
}

// 复制文件
func fileCopy(src, dest string, info os.FileInfo) error {

	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	return err
}

// 复制目录
func dirCopy(srcdir, destdir string, info os.FileInfo) error {

	if err := os.MkdirAll(destdir, info.Mode()); err != nil {
		return err
	}

	contents, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return err
	}

	for _, content := range contents {
		cs, cd := filepath.Join(srcdir, content.Name()), filepath.Join(destdir, content.Name())
		if err := copyDir(cs, cd, content); err != nil {
			return err
		}
	}
	return nil
}

// 复制链接
func linkCopy(src, dest string, info os.FileInfo) error {
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(src, dest)
}
