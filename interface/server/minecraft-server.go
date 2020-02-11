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

	// 修改最大最小使用内存阈值（单位M）
	// 第一个参数为最大值
	// 第二个参数为最小值
	// 为0表示不修改
	SetMaxMinMemory(int, int)

	// 修改服务器名称
	Rename(string)

	// 获取服务配置
	GetServerConf() *json_struct.ServerConf
}
