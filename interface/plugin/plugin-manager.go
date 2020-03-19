package plugin

import (
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
)

// 插件管理服务器，每个mc服务端都会有一个实例
type PluginManager interface {
	_interface.CallBack

	// 获取可用所有插件
	GetAblePlugins() map[string]Plugin

	// 获取可用所有插件
	GetDisablePlugins() map[string]Plugin

	// 根据id禁用插件
	BanPlugin(string)

	// 根据id接触禁用
	UnbanPlugin(string)

	// 处理投递消息
	HandleMessage(*json_struct.ReciveMessage)

	// 通知添加插件(已拥有，则不添加)，可用于动态
	AddPlugin(Plugin)

	// 关闭服务端回调
	CloseMcServerCallBack()

	// 开启服务端回调
	OpenMcServerCallBack()
}
