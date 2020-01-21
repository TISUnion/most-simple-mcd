package utils

import (
	"fmt"
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// 致命错误，退出程序
func PanicError(msg string, err error) {
	panic(fmt.Sprintf("%s, reason: %v", msg, err))
}

// 创建confParam实例
func NewConfParam(Name, ConfVal, description string, level int) *_interface.ConfParam {
	return &_interface.ConfParam{
		ConfVal:        ConfVal,
		DefaultConfVal: ConfVal,
		Name:           Name,
		Level:          level,
		Description:    description,
	}
}

// 去除数组内相同的元素（set化）
func RemoveRepeatedElement(arr []string) []string {
	noRepeatArr := make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			noRepeatArr = append(noRepeatArr, arr[i])
		}
	}
	return noRepeatArr
}

// 判断字符编码是否不是UTF8， 如果不是，尝试转成UTF8
func ParseCharacter(data []byte) ([]byte, error){
	if !IsUTF8(data) {
		if result, err := simplifiedchinese.GBK.NewDecoder().Bytes(data);err !=nil {
			return nil, err
		} else {
			return result, nil
		}
	}
	return data, nil
}

// 是否是UTF8编码
func IsUTF8(data []byte) bool {
	i := 0
	for i < len(data)  {
		if (data[i] & 0x80) == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num - 1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else  {
			//其他情况说明不是utf-8
			return false
		}
	}
	return true
}

// 判断 UTF8首位
func preNUm(data byte) int {
	var mask byte = 0x80
	var num int = 0
	//8bit中首个0bit前有多少个1bits
	for i:=0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}