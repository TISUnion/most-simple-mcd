package server

import json_struct "github.com/TISUnion/most-simple-mcd/json-struct"

type MinecraftServer interface {
	BasicServer
	// 执行命令
	Command(string) error //执行命令

	// 执行tell命令
	TellCommand(string, string) error //执行tell命令

	// 修改内存使用阈值（单位M）
	// 为0表示不修改
	SetMemory(int)

	// 修改服务器名称
	Rename(string)

	// 获取服务配置
	GetServerConf() *json_struct.ServerConf

	// 修改服务配置
	SetServerConf(*json_struct.ServerConf)

	// 获取资源监控服务
	GetServerMonitor() MonitorServer

	// 获取当前服务端uuid
	GetServerEntryId() string

	// 注册订阅该服务端消息管道，PS：必须保证管道必须一直是被消费中！
	RegisterSubscribeMessageChan(chan *json_struct.ReciveMessage)

	// 启动资源监控服务(只有关闭后才会启动)
	StartMonitorServer()

	// 根据插件名称或者id，ban插件
	BanPlugin(string)

	// 根据插件名称或者id，解ban插件
	UnbanPlugin(string)

	// 获取插件配置
	GetPluginsInfo()[]*json_struct.PluginInfo

	// 关闭资源监控服务
	StopMonitorServer()

	// 写入服务器的日志
	WriteLog(string, string)

	// 注册服务端关闭回调, 传入服务端id
	RegisterCloseCallback(func(string))

	// 注册服务端开启回调, 传入服务端id
	RegisterOpenCallback(func(string))

	// 注册服务端保存回调, 传入服务端id
	RegisterSaveCallback(func(string))
}
