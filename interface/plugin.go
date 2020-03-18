package _interface

import (
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
)

type Plugin interface {
	CallBack

	// 插件ID（uuid）
	GetId() string

	// 获取插件名称
	GetName() string

	// 在服务端插件初始化
	Init(server server.MinecraftServer)

	// 是否是全局插件
	IsGlobal() bool

	// 处理投递消息
	HandleMessage(messageType *json_struct.ReciveMessageType)

	// 开启插件
	Start()

	// 停止插件
	Stop()
}
