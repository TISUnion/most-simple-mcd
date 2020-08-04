package server

import (
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/models"
)

type MinecraftServer interface {
	// 服务端适配器
	_interface.ServerAdapter
	// 通用服务器接口
	BasicServer
	// 执行传入命令
	Command(string) error //执行命令

	// 第一个参数为命令，之后的为命令参数
	RunCommand(string, ...string) error

	// 执行tellraw命令
	// 第一个参数为发送对象
	// 第二个参数若为string 自动封装成json
	// 为utils.TellrowMessage 则调用json方法
	// 非string，自动转为json字符串
	TellrawCommand(string, interface{}) error

	// 执行tell命令
	// 第一个参数为发送对象
	// 第二个参数为发送内容
	TellCommand(string, string) error

	// 执行say命令
	SayCommand(string) error

	// 修改内存使用阈值（单位M）
	// 为0表示不修改
	SetMemory(int64)

	// 修改服务器名称
	Rename(string)

	// 获取服务配置
	GetServerConf() *models.ServerConf

	// 修改服务配置
	SetServerConf(*models.ServerConf)

	// 获取资源监控服务
	GetServerMonitor() MonitorServer

	// 获取当前服务端uuid
	GetServerEntryId() string

	// 注册订阅该服务端消息管道，PS：必须保证管道必须一直是被消费中！
	RegisterSubscribeMessageChan(chan *models.ReciveMessage)

	// 启动资源监控服务(只有关闭后才会启动)
	StartMonitorServer()

	// 根据插件名称或者id，ban插件
	BanPlugin(string)

	// 根据插件名称或者id，解ban插件
	UnbanPlugin(string)

	// 获取插件配置
	GetPluginsInfo()[]*models.PluginInfo

	// 关闭资源监控服务
	StopMonitorServer()

	// 写入服务器的日志
	// 第一个参数为日志内容
	// 第二个参数为日志等级
	WriteLog(string, string)

	// 注册服务端关闭回调, 回调传入服务端id
	RegisterCloseCallback(func(string))

	// 注册服务端开启回调, 回调传入服务端id
	RegisterOpenCallback(func(string))

	// 注册服务端保存回调, 回调传入服务端id
	RegisterSaveCallback(func(string))
}
