// server
// 接口
package server

// BasicServer
// 服务器基础接口
type BasicServer interface {
	Start() error
	Stop() error
	Restart() error
}
