package pack_webfile

import (
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/utils"
	"os"
	"path/filepath"
)

// 解压PackCompressData里的web文件
// 主要用于将web静态文件打包进二进制文件中
func UnCompress() error {
	currentPath, _ := utils.GetCurrentPath()
	compressFilepath := filepath.Join(currentPath, constant.COMPRESS_FILE_NAME)
	fl, err := os.OpenFile(compressFilepath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		utils.PanicError("无法创建管理后台文件", err)
	}
	_, err = fl.Write([]byte(PackCompressData))
	if err != nil {
		utils.PanicError("无法写入管理后台文件内容", err)
	}
	fl.Close()

	uncompressFilepath := filepath.Join(currentPath, constant.Web_FILE_DIR_NAME)
	err = utils.UnCompressDir(compressFilepath, uncompressFilepath)
	if err != nil {
		utils.PanicError("无法解压管理后台文件", err)
	}
	// 移除生成的压缩包
	_ = os.Remove(compressFilepath)
	return nil
}
