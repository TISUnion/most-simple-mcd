// _interface
// 接口库
package _interface

import "github.com/TISUnion/most-simple-mcd/interface/server"

// MinecraftContainer
// minecraft服务容器接口
type MinecraftContainer interface {
	GetById(int) (server.MinecraftServer, error)
	GetByName(string) (server.MinecraftServer, error)
	Add(string, server.MinecraftServer) error
	DelById(int) error
	DelByName(int) error

	// Clear
	// 清除所有mc服务器
	Clear() error
}
