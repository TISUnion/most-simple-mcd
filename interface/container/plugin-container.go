package container

import (
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/interface/server"
)

// 插件容器
type PluginContainer interface {
	_interface.CallBack

	// 注册插件
	RegisterPlugin(_interface.Plugin)

	// 新建一个管理器实例
	NewPluginManager(server server.MinecraftServer)

}