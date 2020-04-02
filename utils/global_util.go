package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net"
	"regexp"
	"strconv"
)

const parseMessageRGX = `\[(\d+:\d+:\d+)]\s+\[Server thread/INFO\]:\s+[<|\[]{1}(.+)[>|\]]{1}\s+(.+)`

// 致命错误，退出程序
func PanicError(msg string, err error) {
	panic(fmt.Sprintf("%s, reason: %v", msg, err))
}

// 创建confParam实例
func NewConfParam(Name, ConfVal, description string, level int, IsAlterable bool) *json_struct.ConfParam {
	return &json_struct.ConfParam{
		ConfVal:        ConfVal,
		DefaultConfVal: ConfVal,
		Name:           Name,
		Level:          level,
		Description:    description,
		IsAlterable:    IsAlterable,
	}
}

// 解析mc玩家发言
func ParseMessage(originMsg []byte) *json_struct.ReciveMessage {
	re := regexp.MustCompile(parseMessageRGX)
	match := re.FindStringSubmatch(string(originMsg))
	if len(match) == 4 {
		return &json_struct.ReciveMessage{
			Player:     match[2],
			Time:       match[1],
			Speak:      match[3],
			OriginData: originMsg,
		}
	}
	return &json_struct.ReciveMessage{OriginData: originMsg}
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
			return data, err
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
	addrIp4, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return 0, err
	}
	addrIp6, err := net.ResolveTCPAddr("tcp6", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return 0, err
	}
    // mac 中同一端口可以支持2种ip方式被不同socket监听
	l4, err := net.ListenTCP("tcp4", addrIp4)
	if err != nil {
		return 0, err
	}
	defer l4.Close()
	l6, err := net.ListenTCP("tcp6", addrIp6)
	if err != nil {
		return 0, err
	}
	defer l6.Close()
	return l4.Addr().(*net.TCPAddr).Port, nil
}

// Int转int32
func IntToint32(i int) int32 {
	iStr := strconv.Itoa(i)
	i64, _ := strconv.ParseInt(iStr, 10, 32)
	return int32(i64)
}

// Uint64转float64
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

// 获取cmd命令数组
func GetCommandArr(memory int, runPath string) []string {
	return []string{
		"java",
		"-jar",
		fmt.Sprintf("-Xmx%dM", memory),
		fmt.Sprintf("-Xms%dM", memory),
		runPath,
		"nogui",
	}
}

// md5加密
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// --------------private-----------------
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
