package modules

import (
	"errors"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/utils"
	"github.com/gookit/color"
	"os"
	"strings"
)

var (
	terminalConfs map[string]string
)

// 自带flag包对当前场景不适用，所以自实现获取命令行参数
func InitFlag() map[string]string {
	if terminalConfs != nil {
		return terminalConfs
	}
	terminalConfs = make(map[string]string)
	// 拷贝Args，不污染源数据
	cmdParams := make([]string, len(os.Args))
	copy(cmdParams, os.Args)
	cmdParams = cmdParams[1:]
	cmdParamsLen := len(cmdParams)
	for i := 0; i < cmdParamsLen; i += 2 {
		param := cmdParams[i]
		if string(constant.CLI_THX) == param {
			color.Green.Println(string(constant.CLI_THX_TEXT))
			SendExitSign()
			break
		}
		if strings.HasPrefix(param, constant.PIX) {
			terminalConfs[param[1:]] = cmdParams[i+1]
		} else {
			utils.PanicError(constant.CMD_ERROR, errors.New(fmt.Sprintf("%s格式错误", param)))
		}
	}
	return terminalConfs
}