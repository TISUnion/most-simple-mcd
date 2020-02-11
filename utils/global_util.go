package utils

import (
	"fmt"
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net"
	"strconv"
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
func ParseCharacter(data []byte) ([]byte, error) {
	if !IsUTF8(data) {
		if result, err := simplifiedchinese.GBK.NewDecoder().Bytes(data); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
	return data, nil
}

// GetFreePort
// 获取系统空闲端口
// 如果port为0，则表示随机获取一个空闲端口，不为0则为指定端口
func GetFreePort(port int) (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// IntToint32
func IntToint32(i int) int32 {
	iStr := strconv.Itoa(i)
	i64, _ := strconv.ParseInt(iStr, 10, 32)
	return int32(i64)
}

// Uint64Tofloat64
func Uint64Tofloat64(ui uint64) float64 {
	uiStr := strconv.FormatUint(ui, 10)
	f64, _ := strconv.ParseFloat(uiStr, 64)
	return f64
}

// 是否是UTF8编码
func IsUTF8(data []byte) bool {
	i := 0
	for i < len(data) {
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
			for j := 0; j < num-1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else {
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
	for i := 0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}
