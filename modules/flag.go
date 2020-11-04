package modules

import (
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/gookit/color"
	"os"
	"regexp"
	"strings"
)

var (
	terminalConfs map[string]string
)

// 自带flag包对当前场景不适用，所以自实现获取命令行参数解析
func InitFlag() map[string]string {
	if terminalConfs != nil {
		return terminalConfs
	}
	terminalConfs = make(map[string]string)
	// 合并命令
	argstr := strings.Join(os.Args[1:], " ")
	// 去除多余空格
	rmBlock, _ := regexp.Compile(" *= *")
	argstr = rmBlock.ReplaceAllString(argstr, "=")
	// 拆分参数块
	paramfeilds := strings.Split(argstr, " ")
	// 去除参数前置横杠
	rmCross, _ := regexp.Compile("^-*")
	for _, row := range paramfeilds {
		if row == "" {
			continue
		}
		feild := rmCross.ReplaceAllString(row, "")
		// 拆分参数名和值
		if strings.Contains(feild, "=") {
			feildArr := strings.Split(feild, "=")
			terminalConfs[feildArr[0]] = feildArr[1]
		} else {
			terminalConfs[feild] = constant.TRUE_STR
		}
	}
	if terminalConfs[string(constant.CLI_THX)] != "" {
		color.Green.Println(string(constant.CLI_THX_TEXT))
		SendExitSign()
	}
	return terminalConfs
}
