package service

import (
	"github.com/TISUnion/most-simple-mcd/contants"
	"github.com/TISUnion/most-simple-mcd/utils"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

type ConfValType map[string]string

type TerminalType map[string]string

type ConfKeysType []string

var (
	confKeys      ConfKeysType
	appConf       *Conf
	defaultConfig ConfValType
)

// 参数名常量
const (
	IS_RELOAD_CONF          = "config.auto.reload"          // 自动加载配置文件
	RELOAD_CONF_INTERVAL    = "config.auto.reload.interval" // 自动加载配置文件间隔，单位：毫秒
	CONF_PATH               = "config.path"                 // 配置文件地址
	IS_START_MC_GUI         = "server.gui"                  // 启动gui
	IS_MANAGE_HTTP          = "http.manage.server"          // 启动管理后台
	MANAGE_HTTP_SERVER_PORT = "http.manage.server.port"     // 管理后台服务端口
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
}

func init() {
	confKeys = append(confKeys, IS_RELOAD_CONF, RELOAD_CONF_INTERVAL, CONF_PATH, IS_START_MC_GUI, IS_MANAGE_HTTP, MANAGE_HTTP_SERVER_PORT, WORKSPACE, I18N)

	defaultConfig = make(map[string]string)

	defaultConfig[IS_RELOAD_CONF] = "true"
	defaultConfig[RELOAD_CONF_INTERVAL] = "2000"

	defaultConfig[IS_START_MC_GUI] = "false"
	defaultConfig[IS_MANAGE_HTTP] = "true"
	defaultConfig[MANAGE_HTTP_SERVER_PORT] = "80"

	if workspace, err := utils.GetCurrentPath(); err == nil {
		defaultConfig[WORKSPACE] = workspace
	} else {
		defaultConfig[WORKSPACE] = "./"
	}
	defaultConfig[CONF_PATH] = filepath.Join(defaultConfig["workspace"], "conf/mcd.ini")
	defaultConfig[I18N] = "zh"
}

// loadFilePath
// 获取配置文件目录
func (c *Conf) loadFilePath(terminalConfs TerminalType) {
	// 根据优先级获取配置文件目录
	if path, ok := terminalConfs[CONF_PATH]; ok {
		c.confs[CONF_PATH] = path
	} else if path := os.Getenv(CONF_PATH); path != "" {
		c.confs[CONF_PATH] = path
	}

	// 验证文件是否存在
	path := c.confs[CONF_PATH]
	// 没有文件就创建文件
	if !utils.IsFile(path) {
		if err := utils.CreateFile(path); err != nil {
			// TODO 写入日志

			if path != defaultConfig["CONF_PATH"] {
				// 回退至默认配置
				c.confs[CONF_PATH] = defaultConfig["CONF_PATH"]
				if !utils.IsFile(c.confs[CONF_PATH]) {
					if err := utils.CreateFile(c.confs[CONF_PATH]); err != nil {
						panic(contants.READ_OR_CREATE_CONF_ERROR)
					}
				}
			} else {
				panic(contants.READ_OR_CREATE_CONF_ERROR)
			}
		}
	}
}

// reloadConfig
// 重载配置
func (c *Conf) ReloadConfig() {

}

// loadFileConf
// 加载文件配置
func (c *Conf) loadFileConf() {
	cfg, err := ini.Load(c.confs["CONF_PATH"])
	if err != nil {
		// 文件解析错误 TODO 写入日志
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
	for _, v := range env{
		c.confs[v] = os.Getenv(v)
	}
}

// loadTerminalConf
// 加载命令行配置
func (c *Conf) loadTerminalConf(terminalConfs TerminalType) {

}

// 获取配置
func GetConfObj(terminalConfs TerminalType) *Conf {
	if appConf != nil {
		return appConf
	}

	appConf = &Conf{
		confs:    defaultConfig,
		confKeys: confKeys,
	}
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
