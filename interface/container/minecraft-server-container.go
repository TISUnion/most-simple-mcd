package container

import (
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
)

// MinecraftContainer
// minecraft服务容器接口
type MinecraftContainer interface {
	_interface.CallBack
	GetServerById(int) (server.MinecraftServer, bool)
	GetMirrorServerById(int) (server.MinecraftServer, bool)
	StartById(int) error
	StopById(int) error
	RestartById(int) error
	GetAllServerConf() []*json_struct.ServerConf

	Add(string, server.MinecraftServer) error

	// StopAll
	// 关闭所有mc服务器
	StopAll() error
}
