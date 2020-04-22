package modules

import (
	"errors"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/utils"
	"github.com/gookit/color"
	"github.com/gookit/gcli/v2/interact"
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
		if constant.CLI_MODE == param {
			CliInteraction()
			SendExitSign()
			break
		}
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

func CliInteraction() {
	ans := interact.SelectOne(
		"请选择操作",
		[]string{"创建新插件"},
		"0",
	)
	switch ans {
	case "0":
		var (
			dirname, enName, zhName, description, command, helpDescription string
			isGlobal                                                       bool
		)
		dirname, _ = interact.ReadLine("输入插件包名(文件名)：")
		if dirname == "" {
			dirname, _ = interact.ReadLine("输入插件包名(文件名)：")
		}

		enName, _ = interact.ReadLine("输入插件struct的名称(类名)：")
		if enName == "" {
			enName, _ = interact.ReadLine("输入插件struct的名称(类名)：")
		}

		zhName, _ = interact.ReadLine("输入插件名称(展示用)：")
		if zhName == "" {
			zhName, _ = interact.ReadLine("输入插件名称(展示用)：")
		}

		description, _ = interact.ReadLine("输入插件简述(展示用)：")
		if description == "" {
			description, _ = interact.ReadLine("输入插件简述(展示用)：")
		}

		command, _ = interact.ReadLine("输入插件命令名称：")
		if command == "" {
			command, _ = interact.ReadLine("输入插件命令名称：")
		}
		if interact.Confirm("是否是全局插件：") {
			isGlobal = true
		} else {
			isGlobal = false
		}
		helpDescription = "插件使用方法"

		enName = strings.Title(enName)
		err := CreatePluginTmplFile(dirname, enName, zhName, description, command, helpDescription, isGlobal)
		if err != nil {
			color.Error.Println("生成失败，因为 ", err)
			return
		}
		color.Success.Println("生成成功")
	}
}
