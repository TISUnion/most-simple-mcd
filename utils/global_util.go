package utils

import (
	"fmt"
)

// 致命错误，退出程序
func PanicError(msg string, err error) {
	panic(fmt.Sprintf("%s, reason: %v", msg, err))
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
