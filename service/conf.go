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

type TerminalType map[string]*string

var (
	appConf          *Conf
	DefaultConfParam map[string]*_interface.ConfParam
)

// Conf
// 首次导入配置优先级： 命令行变量 > 环境变量 > 配置文件 > 默认配置
// 非首次加载：		 配置文件 > 环境变量 > 默认配置
type Conf struct {
	// Confs
	// 配置
	confs map[string]*_interface.ConfParam

	// ConfKeys
	// 所有配置键值
	ConfKeys []string

	// lock
	// 读写锁
	lock *sync.RWMutex
}

func init() {
	DefaultConfParam = make(map[string]*_interface.ConfParam)
	DefaultConfParam[constant.IS_RELOAD_CONF] = utils.NewConfParam(constant.IS_RELOAD_CONF, "true", constant.IS_RELOAD_CONF_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.RELOAD_CONF_INTERVAL] = utils.NewConfParam(constant.RELOAD_CONF_INTERVAL, "2000", constant.RELOAD_CONF_INTERVAL_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.IS_START_MC_GUI] = utils.NewConfParam(constant.IS_START_MC_GUI, "false", constant.IS_START_MC_GUI_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.IS_MANAGE_HTTP] = utils.NewConfParam(constant.IS_MANAGE_HTTP, "true", constant.IS_MANAGE_HTTP_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.MANAGE_HTTP_SERVER_PORT] = utils.NewConfParam(constant.MANAGE_HTTP_SERVER_PORT, "80", constant.MANAGE_HTTP_SERVER_PORT_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.LOG_SAVE_INTERVAL] = utils.NewConfParam(constant.LOG_SAVE_INTERVAL, constant.LOG_SAVE_INTERVAL_TWICEDAY, constant.LOG_SAVE_INTERVAL_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.LOG_SHOW_CODELINE] = utils.NewConfParam(constant.LOG_SHOW_CODELINE, "false", constant.LOG_SHOW_CODELINE_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	if workspace, err := utils.GetCurrentPath(); err == nil {
		DefaultConfParam[constant.WORKSPACE] = utils.NewConfParam(constant.WORKSPACE, workspace, constant.WORKSPACE_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	} else {
		DefaultConfParam[constant.WORKSPACE] = utils.NewConfParam(constant.WORKSPACE, "./", constant.WORKSPACE_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	}
	DefaultConfParam[constant.LOG_PATH] = utils.NewConfParam(constant.LOG_PATH, filepath.Join(DefaultConfParam[constant.WORKSPACE].ConfVal, "logs"), constant.LOG_PATH_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.CONF_PATH] = utils.NewConfParam(constant.CONF_PATH, filepath.Join(DefaultConfParam[constant.WORKSPACE].ConfVal, "conf/mcd.ini"), constant.CONF_PATH_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.I18N] = utils.NewConfParam(constant.I18N, "zh", constant.I18N_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
}

// loadFilePath
// 获取配置文件目录
func (c *Conf) loadFilePath(terminalConfs TerminalType) {
	// 根据优先级获取配置文件目录
	if path, ok := terminalConfs[constant.CONF_PATH]; ok && *path != "" {
		c.SetConfParam(constant.CONF_PATH, *path, constant.CONF_PATH_DESCREPTION, constant.CONF_TERMINAL_LEVEL)
	} else if path := os.Getenv(constant.CONF_PATH); path != "" {
		c.SetConfParam(constant.CONF_PATH, path, constant.CONF_PATH_DESCREPTION, constant.CONF_ENVIRONMENT_LEVEL)
	}
	// 验证文件是否存在
	path := c.confs[constant.CONF_PATH].ConfVal
	// 没有文件就创建文件
	if !utils.ExistsResource(path) {
		var (
			f   *os.File
			err error
		)
		if f, err = utils.CreateFile(path); err != nil {
			if path != c.confs[constant.CONF_PATH].DefaultConfVal {
				confParam := c.confs[constant.CONF_PATH]
				// 回退至默认配置
				confParam.ConfVal = confParam.DefaultConfVal
				confParam.Level = constant.CONF_DEFAULT_LEVEL
				if !utils.ExistsResource(confParam.ConfVal) {
					if f, err = utils.CreateFile(confParam.ConfVal); err != nil {
						utils.PanicError(constant.CREATE_CONF_ERROR, err)
					}
				}
			} else {
				utils.PanicError(constant.CREATE_CONF_ERROR, err)
			}
		}
		defer f.Close()
	}
	// 配置文件内容为空就写入默认配置
	data, err := ioutil.ReadFile(c.confs[constant.CONF_PATH].ConfVal)
	if err != nil {
		utils.PanicError(constant.READ_CONF_ERROR, err)
	}
	if len(data) == 0 {
		cfg := setIniCfg(c.confs)
		if err := cfg.SaveTo(c.confs[constant.CONF_PATH].ConfVal); err != nil {
			utils.PanicError(constant.WRITE_CONF_ERROR, err)
		}
	}
}

// loadPluginsConfKeysType
// TODO 加载插件的所有配置键
func (c *Conf) loadPluginsConf() {

}

// reloadConfig
// TODO 重载配置
func (c *Conf) ReloadConfig() {
	c.lock.Lock()
	defer c.lock.Unlock()
}

// loadDefaultConf
// 加载默认配置
func (c *Conf) loadDefaultConf() {
	for _, confParam := range DefaultConfParam {
		c.SetConfParam(confParam.Name, confParam.ConfVal, confParam.Description, confParam.Level)
	}
}

// loadFileConf
// 加载文件配置
func (c *Conf) loadFileConf() {
	cfg, err := ini.Load(c.confs[constant.CONF_PATH].ConfVal)
	if err != nil {
		// 文件解析错误
		utils.PanicError(constant.PARSE_INI_CONF_ERROR, err)
	} else {
		sec := cfg.Section("")
		keys := sec.KeyStrings()
		for _, v := range keys {
			c.SetConfParam(v, sec.Key(v).String(), "", constant.CONF_FILE_LEVEL)
		}
	}
}

// loadEnvConf
// 加载环境变量
func (c *Conf) loadEnvConf() {
	env := os.Environ()
	for _, v := range env {
		c.SetConfParam(v, os.Getenv(v), "", constant.CONF_ENVIRONMENT_LEVEL)
	}
}

// loadTerminalConf
// 加载命令行配置
func (c *Conf) loadTerminalConf(terminalConfs TerminalType) {
	if terminalConfs != nil {
		for k, v := range terminalConfs {
			c.SetConfParam(k, *v, "", constant.CONF_TERMINAL_LEVEL)
		}
	}
}

func (c *Conf) GetConfig() map[string]string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	result := make(map[string]string)
	for k, v := range c.confs {
		result[k] = v.ConfVal
	}
	return result
}

// 获取配置键
func (c *Conf) GetConfigKeys() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.ConfKeys
}

// 获取单个配置值
func (c *Conf) GetConfVal(key string) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if val, ok := c.confs[key]; ok {
		return val.ConfVal
	}
	return ""
}

func (c *Conf) SetConfig(key string, val string) {
	if key != "" {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.confs[key].ConfVal = val
	}
}

func (c *Conf) Init(terminalConfs TerminalType) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// 加载文件配置文件路径
	c.loadFilePath(terminalConfs)

	// 加载插件配置文件
	c.loadPluginsConf()

	// 加载配置文件
	c.loadFileConf()

	// 加载环境变量
	c.loadEnvConf()

	// 加载命令行参数
	c.loadTerminalConf(terminalConfs)

}

// 设置配置
func (c *Conf) SetConfParam(Name, ConfVal, description string, level int) {
	if confParam, ok := c.confs[Name]; ok {
		// 如果配置优先级不低于与存在配置，就修改
		if confParam.Level <= level {
			confParam.Level = level
			confParam.ConfVal = ConfVal
			if description != "" {
				confParam.Description = description
			}
		}
	} else {
		// 如果不存在配置，就记录配置名并新建配置对象
		c.ConfKeys = append(c.ConfKeys, Name)
		c.confs[Name] = &_interface.ConfParam{
			ConfVal:        ConfVal,
			Name:           Name,
			Level:          level,
			Description:    description,
			DefaultConfVal: ConfVal,
		}
	}
}

// 获取配置实例
func GetConfObj(terminalConfs TerminalType) _interface.Conf {
	if appConf != nil {
		return appConf
	}
	rwLock := &sync.RWMutex{}
	appConf = &Conf{
		lock:     rwLock,
		ConfKeys: make([]string, 0),
		confs:    make(map[string]*_interface.ConfParam),
	}
	appConf.Init(terminalConfs)
	return appConf
}

// 这是ini配置对象
func setIniCfg(data map[string]*_interface.ConfParam) *ini.File {
	cfg := ini.Empty()
	sec, _ := cfg.NewSection("DEFAULT")
	for k, v := range data {
		confVal := "\"" + v.ConfVal + "\""
		_, _ = sec.NewKey(k, confVal)
	}
	return cfg
}
