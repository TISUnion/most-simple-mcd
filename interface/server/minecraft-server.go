package server

import json_struct "github.com/TISUnion/most-simple-mcd/json-struct"

type ReciveMessageType struct {
	Player     string
	Time       string
	Speak      string
	OriginData []byte
	ServerId   int
}

type MinecraftServer interface {
	BasicServer
	// 第一参数为命令，后面为参数
	Command(string) error //执行命令

	// 修改内存使用阈值（单位M）
	// 为0表示不修改
	SetMemory(int)

	// 修改服务器名称
	Rename(string)

	// 获取服务配置
	GetServerConf() *json_struct.ServerConf

	// 获取资源监控服务
	GetServerMonitor() MonitorServer

	// 获取当前全局id,每次启动id不一定相同 TODO 后期优化
	GetServerEntryId() int

	// 启动资源监控服务(只有关闭后才会启动)
	StartMonitorServer()
}
