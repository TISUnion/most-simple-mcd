package main

import (
	"github.com/TISUnion/most-simple-mcd/tool/lib"
	"github.com/gookit/color"
	"github.com/gookit/gcli/v2/interact"
	"strings"
)

func main() {
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
	err := lib.CreatePluginTmplFile(dirname, enName, zhName, description, command, helpDescription, isGlobal)
	if err != nil {
		color.Error.Println("生成失败，因为 ", err)
		return
	}
	color.Success.Println("生成成功")
}
