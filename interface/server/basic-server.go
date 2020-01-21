// server
// 接口
package server

import _interface "github.com/TISUnion/most-simple-mcd/interface"

// BasicServer
// 服务器基础接口
type BasicServer interface {
	_interface.CallBack
	Start() error
	Stop() error
	Restart() error
}
