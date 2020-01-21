package service

import (
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/utils"
	"gopkg.in/ini.v1"
	"io"
	"os/exec"
	"path/filepath"
	"sync"
)

const (
	EULA_FILE_NAME = "eula.txt"
	EULA           = "eula"
	TRUE_STR       = "true"
)

// MinecraftServer
// mc服务器子进程对象
type MinecraftServer struct {
	*json_struct.ServerConf

	// CmdObj
	//子进程实例
	CmdObj *exec.Cmd

	// stdin
	// 用于关闭输入管道
	stdin io.WriteCloser

	// stdout
	// 子进程输出
	stdout io.ReadCloser

	// lock
	// 输入管道同步锁
	lock *sync.Mutex

	// isStart
	// 是否启动
	isStart bool
}

func (m *MinecraftServer) ChangeConfCallBack() {
}

func (m *MinecraftServer) DestructCallBack() {
	_ = m.stdin.Close()
	_ = m.stdout.Close()
}

func (m *MinecraftServer) Start() error {
	return nil
}

func (m *MinecraftServer) Stop() error {
	panic("implement me")
}

func (m *MinecraftServer) Restart() error {
	panic("implement me")
}

func (m *MinecraftServer) Command(string, ...interface{}) error {
	panic("implement me")
}

func (m *MinecraftServer) validatePort() int {
	return 0
}

func (m *MinecraftServer) validateEula() error {
	path := filepath.Join(m.RunPath, EULA_FILE_NAME)
	f, _ := utils.CreateFile(path)
	f.Close()
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}
	sec, err := cfg.GetSection("")
	if err != nil {
		return err
	}
	if sec.HasKey(EULA) {
		eula, err := sec.Key(EULA).Bool()
		if err != nil {
			eula = false
		}
		if !eula {
			sec.Key(EULA).SetValue(TRUE_STR)
			_ = cfg.SaveTo(path)
		}
	} else {
		_, _ = sec.NewKey(EULA, TRUE_STR)
		_ = cfg.SaveTo(path)
	}
	return nil
}

func NewMinecraftServer(serverConf *json_struct.ServerConf) server.MinecraftServer {
	cmdObj := exec.Command(serverConf.CmdStr)
	cmdObj.Dir = serverConf.RunPath
	stdin, err := cmdObj.StdinPipe()
	if err != nil {
		return nil
	}
	stdout, err := cmdObj.StdoutPipe()
	if err != nil {
		return nil
	}
	minecraftServer := &MinecraftServer{
		ServerConf: serverConf,
		CmdObj:     cmdObj,
		stdin:      stdin,
		stdout:     stdout,
		lock:       &sync.Mutex{},
		isStart:    false,
	}
	return minecraftServer
}
