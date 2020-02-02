package service

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/utils"
	"gopkg.in/ini.v1"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
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

	// Pid
	// 进程pid
	Pid int

	// lock
	// 输入管道同步锁
	lock *sync.Mutex

	// isStart
	// 是否启动
	isStart bool
}

func (m *MinecraftServer) ChangeConfCallBack() {
}

func (m *MinecraftServer) InitCallBack() {
}

func (m *MinecraftServer) DestructCallBack() {
	_ = m.stdin.Close()
	_ = m.stdout.Close()
}

func (m *MinecraftServer) runProcess() error {
	if err := m.validateEula(); err != nil {
		return err
	}
	m.CmdObj.Dir = m.RunPath
	if err := m.CmdObj.Start(); err != nil {
		return err
	}
	return nil
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

// validatePort
// 校验mc的端口 TODO
func (m *MinecraftServer) ValidatePort() (int, error) {
	path := filepath.Join(m.RunPath, constant.MC_CONF_NAME)
	cfg, err := ini.Load(path)
	// 没有配置文件
	if err != nil || !cfg.Section("").HasKey(constant.MAX_PLAYER_PARAM) {
		// 开启进程，自动创建
		if err := m.runProcess(); err != nil {
			return 0, err
		}

		// 接收进程信息
		for {
			data, err := ioutil.ReadAll(m.stdout)
			if len(data) > 0 {
				fmt.Println(string(data))
				break
			}

			if err != nil {
				fmt.Println(err)
				break
			}
		}
	} else {
		port, _ := cfg.Section("").Key(constant.PORT).Int()
		// 开启的服务端的端口已被占用
		if p, _ := utils.GetFreePort(port); p == 0 {
			// 如果可以自动更换端口就自动更换端口
			if isChange, _ := strconv.ParseBool(GetConfInstance().GetConfVal(constant.IS_AUTO_CHANGE_MC_SERVER_REPEAT_PORT)) ; isChange {

			}
		}
	}
	return 0, nil
}

// validateEula
// 校验mc的eula文件
func (m *MinecraftServer) validateEula() error {
	path := filepath.Join(m.RunPath, constant.EULA_FILE_NAME)
	f, _ := utils.CreateFile(path)
	_ = f.Close()
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}
	sec, err := cfg.GetSection("")
	if err != nil {
		return err
	}
	if sec.HasKey(constant.EULA) {
		eula, err := sec.Key(constant.EULA).Bool()
		if err != nil {
			eula = false
		}
		if !eula {
			sec.Key(constant.EULA).SetValue(constant.TRUE_STR)
			_ = cfg.SaveTo(path)
		}
	} else {
		_, _ = sec.NewKey(constant.EULA, constant.TRUE_STR)
		_ = cfg.SaveToIndent(path, "\t")
	}
	return nil
}

// NewMinecraftServer
// 新建一个mc服务端进程
func NewMinecraftServer(serverConf *json_struct.ServerConf) *MinecraftServer {
	cmdObj := exec.Command(serverConf.CmdStr[0], serverConf.CmdStr[1:]...)
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
