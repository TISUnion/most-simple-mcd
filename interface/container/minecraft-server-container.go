package container

import (
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/models"
)

// MinecraftContainer
// minecraft服务容器接口
type MinecraftContainer interface {
	// 回调
	_interface.CallBack
	// 根据id获取服务端实例
	GetServerById(string) (server.MinecraftServer, error)

	// 根据id获取镜像服务端实例
	GetMirrorServerById(string) (server.MinecraftServer, error)

	// 根据id开启服务端
	StartById(string) error

	// 启动所有服务端
	StartAll() error

	// 根据id停止服务端
	StopById(string) error

	// 根据id重启服务端
	RestartById(string) error

	// 获取所有服务端配置
	GetAllServerConf() []*json_struct.ServerConf

	// 添加服务端
	AddServer(*json_struct.ServerConf, bool)

	// 生成服务端参数对象
	HandleMcFile(string, string, int, int) *json_struct.ServerConf

	// StopAll
	// 关闭所有mc服务器
	StopAll() error

	// 获取所有服务端对象实例
	GetAllServerObj() map[string]server.MinecraftServer

	// 所有服务端配置保存到数据库中
	SaveToDb()

	// 注册所有服务端关闭回调, 回调函数会传入服务端id
	RegisterAllServerCloseCallback(func(string))

	// 注册所有服务端开启回调, 回调函数会传入服务端id
	RegisterAllServerOpenCallback(func(string))

	// 注册所有服务端保存回调, 回调函数会传入服务端id
	RegisterAllServerSaveCallback(func(string))
}
