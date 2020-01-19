package service

import (
	"github.com/TISUnion/most-simple-mcd/constant"
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/utils"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type ConfValType map[string]string

type TerminalType map[string]*string

type ConfKeysType []string

var (
	DefaultConfKeys ConfKeysType
	appConf         *Conf
	DefaultConfig   ConfValType
)

// 参数名常量
const (
	IS_RELOAD_CONF          = "config.auto.reload"          // 自动加载配置文件
	RELOAD_CONF_INTERVAL    = "config.auto.reload.interval" // 自动加载配置文件间隔，单位：毫秒
	CONF_PATH               = "config.path"                 // 配置文件地址
	IS_MANAGE_HTTP          = "http.manage.server"          // 启动管理后台
	MANAGE_HTTP_SERVER_PORT = "http.manage.server.port"     // 管理后台服务端口
	LOG_PATH                = "log.path"                    // 日志写入目录
	LOG_SAVE_INTERVAL		= "log.interval"				// 日志保存间隔，例如: 每2天对久日志压缩，日志写入新日志中
	IS_START_MC_GUI         = "server.gui"                  // 启动gui
	WORKSPACE               = "workspace"                   // 工作目录
	I18N                    = "i18n"                        // 国际化
)

// Conf
// 首次导入配置优先级： 命令行变量 > 环境变量 > 配置文件 > 默认配置
// 非首次加载：		 配置文件 > 环境变量 > 默认配置
type Conf struct {
	// Confs
	// 配置
	confs ConfValType

	// ConfKeys
	// 所有配置键值
	confKeys ConfKeysType

	// lock
	// 读写锁
	lock *sync.RWMutex
}

func init() {
	// 定义默认配置
	DefaultConfKeys = append(DefaultConfKeys, IS_RELOAD_CONF, RELOAD_CONF_INTERVAL, CONF_PATH, IS_START_MC_GUI, IS_MANAGE_HTTP, MANAGE_HTTP_SERVER_PORT, WORKSPACE, I18N)

	DefaultConfig = make(map[string]string)

	DefaultConfig[IS_RELOAD_CONF] = "true"
	DefaultConfig[RELOAD_CONF_INTERVAL] = "2000"

	DefaultConfig[IS_START_MC_GUI] = "false"
	DefaultConfig[IS_MANAGE_HTTP] = "true"
	DefaultConfig[MANAGE_HTTP_SERVER_PORT] = "80"

	if currentPath, err := utils.GetCurrentPath(); err == nil {
		DefaultConfig[LOG_PATH]  = filepath.Join(currentPath, "logs")
	} else {
		DefaultConfig[LOG_PATH]  = filepath.Join("./", "logs")
	}

	DefaultConfig[LOG_SAVE_INTERVAL]  = constant.LOG_SAVE_INTERVAL_TWICEDAY

	if workspace, err := utils.GetCurrentPath(); err == nil {
		DefaultConfig[WORKSPACE] = workspace
	} else {
		DefaultConfig[WORKSPACE] = "./"
	}
	DefaultConfig[CONF_PATH] = filepath.Join(DefaultConfig["workspace"], "conf/mcd.ini")
	DefaultConfig[I18N] = "zh"
}

// loadFilePath
// 获取配置文件目录
func (c *Conf) loadFilePath(terminalConfs TerminalType) {
	// 根据优先级获取配置文件目录
	if path, ok := terminalConfs[CONF_PATH]; ok {
		c.confs[CONF_PATH] = *path
	} else if path := os.Getenv(CONF_PATH); path != "" {
		c.confs[CONF_PATH] = path
	}

	// 验证文件是否存在
	path := c.confs[CONF_PATH]
	// 没有文件就创建文件
	if !utils.IsFile(path) {
		var (
			f   *os.File
			err error
		)
		if f, err = utils.CreateFile(path); err != nil {
			// 写入日志
			utils.WriteLog(constant.CREATE_CONF_FAILED_AND_ROLLBACK, constant.LOG_ERROR)
			if path != DefaultConfig["CONF_PATH"] {
				// 回退至默认配置
				c.confs[CONF_PATH] = DefaultConfig["CONF_PATH"]
				if !utils.IsFile(c.confs[CONF_PATH]) {
					if f, err = utils.CreateFile(c.confs[CONF_PATH]); err != nil {
						utils.PanicError(constant.CREATE_CONF_ERROR)
					}
				}
			} else {
				utils.PanicError(constant.CREATE_CONF_ERROR)
			}
		}
		defer f.Close()
	}

	// 配置文件内容为空就写入默认配置
	data, err := ioutil.ReadFile(c.confs[CONF_PATH])
	if err != nil {
		utils.PanicError(constant.READ_CONF_ERROR)
	}
	if len(data) == 0 {
		cfg := setIniCfg(DefaultConfig)
		if err := cfg.SaveTo(c.confs[CONF_PATH]); err != nil {
			utils.PanicError(constant.WRITE_CONF_ERROR)
		}
	}
}

// loadPluginsConfKeysType
// TODO 加载插件的所有配置键
func (c *Conf) loadPluginsConfKeysType() {

}

// reloadConfig
// 重载配置
func (c *Conf) ReloadConfig() {
	c.lock.Lock()
	defer c.lock.Unlock()

	// 加载环境变量
	appConf.loadEnvConf()

	// 加载配置文件
	appConf.loadFileConf()
}

// loadFileConf
// 加载文件配置
func (c *Conf) loadFileConf() {
	cfg, err := ini.Load(c.confs["CONF_PATH"])
	if err != nil {
		// 文件解析错误
		utils.PanicError(constant.PARSE_INI_CONF_ERROR)
	} else {
		sec := cfg.Section("")
		keys := sec.KeyStrings()
		for _, v := range keys {
			c.confs[v] = sec.Key(v).String()
		}
	}
}

// loadEnvConf
// 加载环境变量
func (c *Conf) loadEnvConf() {
	env := os.Environ()
	for _, v := range env {
		c.confs[v] = os.Getenv(v)
	}
}

// loadTerminalConf
// 加载命令行配置
func (c *Conf) loadTerminalConf(terminalConfs TerminalType) {
	if terminalConfs != nil {
		for k, v := range terminalConfs {
			c.confs[k] = *v
		}
	}
}

func (c *Conf) GetConfig() map[string]string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.confs
}

func (c *Conf) GetConfigKeys() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.confKeys
}

func (c *Conf) GetConfVal(key string) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if val, ok := c.confs[key]; ok {
		return val
	}
	return ""
}

func (c *Conf) SetConfig(key string, val string) {
	if key != "" {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.confs[key] = val
	}
}

// 获取配置实例
func GetConfObj(terminalConfs TerminalType) _interface.Conf {
	if appConf != nil {
		return appConf
	}
	rwLock := &sync.RWMutex{}
	appConf = &Conf{
		confs:    DefaultConfig,
		confKeys: DefaultConfKeys,
		lock:     rwLock,
	}

	rwLock.Lock()
	defer rwLock.Unlock()
	// 加载文件配置文件路径
	appConf.loadFilePath(terminalConfs)

	// 加载配置文件
	appConf.loadFileConf()

	// 加载环境变量
	appConf.loadEnvConf()

	// 加载命令行参数
	appConf.loadTerminalConf(terminalConfs)

	return appConf
}

// 这是ini配置对象
func setIniCfg(data map[string]string) *ini.File {
	cfg := ini.Empty()
	sec, _ := cfg.NewSection("")
	for k, v := range data {
		_, _ = sec.NewKey(k, v)
	}
	return cfg
}
