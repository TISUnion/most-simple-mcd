package utils

import "github.com/TISUnion/most-simple-mcd/wire"

// 致命错误，退出程序
func PanicError(msg string) {
	panic(msg)
}

// 判断字符串是否在数组中
func StrInArr(item string, arr []string) bool{
	for _, v := range arr {
		if item == v {
			return true
		}
	}
	return false
}

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

func WriteLog(msg string, level string) {
	wire.GetLogContainerInstance().WriteLog(msg, level)
}

