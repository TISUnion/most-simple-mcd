//+build compress web file
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"fmt"
	"github.com/klauspost/compress/flate"
	"github.com/mholt/archiver/v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

/**
将前端编译好的文件写入到go文件中，达到把静态文件打包到exe中的 目地
*/
func main() {
	originWebDir := "web-admin/dist/*"
	compressFileName := "webfile.zip"
	codeFileName := "pack-webfile/pack-data.go"
	originWebDirFiles, _ := filepath.Glob(originWebDir)
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      true,
		ImplicitTopLevelFolder: false,
	}
	if err := z.Archive(originWebDirFiles, compressFileName); err != nil {
		fmt.Println(err)
		return
	}

	if compressFileData, err := ioutil.ReadFile(compressFileName); err != nil {
		fmt.Println(err)
		return
	} else {
		compressFileDataStr := ""
		for _, v := range compressFileData {
			compressFileDataStr += fmt.Sprintf("%d,", int(v))
		}
		compressFileDataStr = compressFileDataStr[:len(compressFileDataStr)-1]
		compressFileDataStr = fmt.Sprintf("{%s}", compressFileDataStr)
		//fmt.Println(compressFileDataStr)
		if codeData, err := ioutil.ReadFile(codeFileName); err != nil {
			fmt.Println(err)
			return
		} else {
			codeDataStr := string(codeData)
			if reg, err := regexp.Compile("\\{(.*)\\}"); err != nil {
				fmt.Println(err)
				return
			} else {
				str := reg.ReplaceAllString(codeDataStr, compressFileDataStr)
				if codefile, err := os.OpenFile(codeFileName, os.O_WRONLY, 0777); err != nil {
					fmt.Println(err)
					return
				} else {
					_, err := codefile.Write([]byte(str))
					if err != nil {
						fmt.Println(err)
					}
					codefile.Close()
				}
				_ = os.Remove(compressFileName)
			}
		}
	}
}
