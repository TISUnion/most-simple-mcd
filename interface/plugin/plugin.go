package plugin

import (
	"github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
)

type Plugin interface {
	_interface.CallBack

	// 插件ID（uuid）
	GetId() string

	// 获取插件名称
	GetName() string

	// 获取简介
	GetDescription() string

	// 获取使用说明
	GetHelpDescription() string

	// 获取命令
	GetCommandName() string

	// 非全局插件，会服务端插件初始化
	Init(server server.MinecraftServer)

	// 是否是全局插件
	IsGlobal() bool

	// 如果是非全局插件就要提供一个新建插件实体的函数
	NewInstance() Plugin

	// 处理投递消息
	HandleMessage(messageType *json_struct.ReciveMessage)

	// 开启插件
	Start()

	// 停止插件
	Stop()
}
