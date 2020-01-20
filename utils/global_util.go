package utils

import (
	"fmt"
	_interface "github.com/TISUnion/most-simple-mcd/interface"
)

// 致命错误，退出程序
func PanicError(msg string, err error) {
	panic(fmt.Sprintf("%s, reason: %v", msg, err))
}

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
