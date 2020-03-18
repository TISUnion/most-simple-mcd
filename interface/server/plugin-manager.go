package server

import (
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
)

// 插件管理器，每个mc服务端都会有一个实例
type PluginManager interface {
	_interface.CallBack

	// 获取可用所有插件
	GetPlugins() map[string]_interface.Plugin

	// 根据id禁用插件
	BanPlugin(string)

	// 根据id接触禁用
	UnbanPlugin(string)

	// 处理投递消息
	HandleMessage(messageType *json_struct.ReciveMessageType)
}
