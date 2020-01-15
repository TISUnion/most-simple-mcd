package service

import (
	"bufio"
	"io"
	"os/exec"
	"sync"
)

type NameType string

type CmdStrType string

type PortType string

// MinecraftServer
// mc服务器
type MinecraftServer struct {
	// Name
	// 服务器名称
	Name NameType
	// CmdObj
	//子进程实例
	CmdObj *exec.Cmd
	// CmdStr
	// 执行的完整命令
	CmdStr CmdStrType
	// stdin
	// 用于关闭输入管道
	stdin io.WriteCloser
	// stdout
	// 子进程输出
	stdout *bufio.Reader
	// lock
	// 输入管道同步锁
	lock sync.Mutex
	// Port
	// 启动服务器端口
	Port PortType
}

func (*MinecraftServer) init(serverName NameType, CmdStr CmdStrType, Port PortType) {

}
