//+build compress web file
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"bytes"
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

const TMPL = `package pack_webfile

var PackCompressData = []byte{}
`

func main() {
	originWebDir := "web-admin/dist/*"
	compressFileName := "webfile.zip"
	codeFileName := "pack-webfile/pack-data.go"
	initCodeFile(codeFileName)
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
		panic(err)
	}

	if compressFileData, err := ioutil.ReadFile(compressFileName); err != nil {
		panic(err)
	} else {
		compressFileDataStr := ""
		buff := bytes.NewBuffer([]byte{})
		for _, v := range compressFileData {
			buff.WriteString(fmt.Sprintf("%d,", int(v)))
		}
		compressFileDataStr = buff.String()
		compressFileDataStr = compressFileDataStr[:len(compressFileDataStr)-1]
		compressFileDataStr = fmt.Sprintf("{%s}", compressFileDataStr)
		if codeData, err := ioutil.ReadFile(codeFileName); err != nil {
			panic(err)
		} else {
			codeDataStr := string(codeData)
			if reg, err := regexp.Compile("\\{(.*)\\}"); err != nil {
				panic(err)
			} else {
				str := reg.ReplaceAllString(codeDataStr, compressFileDataStr)
				if codefile, err := os.OpenFile(codeFileName, os.O_WRONLY, 0777); err != nil {
					panic(err)
				} else {
					_, err := codefile.Write([]byte(str))
					if err != nil {
						panic(err)
					}
					codefile.Close()
				}
				_ = os.Remove(compressFileName)
			}
		}
	}
}

func initCodeFile(path string) {
	os.Remove(path)
	file, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666);
	file.Write([]byte(TMPL))
}
