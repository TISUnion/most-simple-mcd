package service

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/utils"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type TerminalType map[string]*string

var (
	_appConf         *Conf
	DefaultConfParam map[string]*_interface.ConfParam
)

// Conf
// 导入配置优先级： 命令行变量 > 环境变量 > 配置文件 > 默认配置
type Conf struct {
	// Confs
	// 配置
	confs map[string]*_interface.ConfParam

	// ConfKeys
	// 所有配置键值
	ConfKeys []string

	// lock
	// 读写锁
	lock *sync.Mutex
}

func ConfInit() {
	// 加载默认配置
	DefaultConfParam = make(map[string]*_interface.ConfParam)
	DefaultConfParam[constant.IS_RELOAD_CONF] = utils.NewConfParam(constant.IS_RELOAD_CONF, "false", constant.IS_RELOAD_CONF_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.RELOAD_CONF_INTERVAL] = utils.NewConfParam(constant.RELOAD_CONF_INTERVAL, "5000", constant.RELOAD_CONF_INTERVAL_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.IS_MANAGE_HTTP] = utils.NewConfParam(constant.IS_MANAGE_HTTP, "true", constant.IS_MANAGE_HTTP_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.MANAGE_HTTP_SERVER_PORT] = utils.NewConfParam(constant.MANAGE_HTTP_SERVER_PORT, "80", constant.MANAGE_HTTP_SERVER_PORT_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.LOG_SAVE_INTERVAL] = utils.NewConfParam(constant.LOG_SAVE_INTERVAL, constant.LOG_SAVE_INTERVAL_TWICEDAY, constant.LOG_SAVE_INTERVAL_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	if workspace, err := utils.GetCurrentPath(); err == nil {
		DefaultConfParam[constant.WORKSPACE] = utils.NewConfParam(constant.WORKSPACE, workspace, constant.WORKSPACE_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	} else {
		DefaultConfParam[constant.WORKSPACE] = utils.NewConfParam(constant.WORKSPACE, "./", constant.WORKSPACE_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	}
	DefaultConfParam[constant.LOG_PATH] = utils.NewConfParam(constant.LOG_PATH, filepath.Join(DefaultConfParam[constant.WORKSPACE].ConfVal, "logs"), constant.LOG_PATH_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.CONF_PATH] = utils.NewConfParam(constant.CONF_PATH, filepath.Join(DefaultConfParam[constant.WORKSPACE].ConfVal, "conf/mcd.ini"), constant.CONF_PATH_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.IS_AUTO_CHANGE_MC_SERVER_REPEAT_PORT] = utils.NewConfParam(constant.IS_AUTO_CHANGE_MC_SERVER_REPEAT_PORT, "true", constant.IS_AUTO_CHANGE_MC_SERVER_REPEAT_PORT_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	DefaultConfParam[constant.MONITOR_INTERVAL] = utils.NewConfParam(constant.MONITOR_INTERVAL, "2s", constant.MONITOR_INTERVAL_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
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
// 重新再加配置
func (c *Conf) ReloadConfig() {
	c.lock.Lock()
	defer c.lock.Unlock()
	// 加载插件配置文件
	c.loadPluginsConf()
	// 加载配置文件
	c.loadFileConf()
	// 加载环境变量
	c.loadEnvConf()
	// 执行配置更改回调
	RunChangeConfCallBacks()
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
		envkvAr := strings.Split(v, "=")
		if len(envkvAr) >= 2 && envkvAr[1] != "" {
			c.SetConfParam(envkvAr[0], envkvAr[1], "", constant.CONF_ENVIRONMENT_LEVEL)
		}
	}
}

// loadTerminalConf
// 加载命令行配置
func (c *Conf) loadTerminalConf(terminalConfs TerminalType) {
	if terminalConfs != nil {
		for k, v := range terminalConfs {
			if *v != "" {
				c.SetConfParam(k, *v, "", constant.CONF_TERMINAL_LEVEL)
			}
		}
	}
}

func (c *Conf) GetConfig() map[string]string {
	result := make(map[string]string)
	for k, v := range c.confs {
		result[k] = v.ConfVal
	}
	return result
}

// 获取配置键
func (c *Conf) GetConfigKeys() []string {
	return c.ConfKeys
}

// 获取单个配置值
func (c *Conf) GetConfVal(key string) string {
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
	// 加载默认配置
	c.loadDefaultConf()

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

func (c *Conf) ChangeConfCallBack() {
	_switch, err := strconv.ParseBool(c.GetConfVal(constant.IS_RELOAD_CONF))
	if err != nil {
		_switch = false
	}
	// 启动自动加载配置任务
	jobContainer := GetJobContainerInstance()
	if _switch {
		interval := fmt.Sprintf("@every %sms", c.GetConfVal(constant.RELOAD_CONF_INTERVAL))
		if !jobContainer.HasJob(constant.RELOAD_CONF_JOB_NAME) {
			jobContainer.RegisterJob(constant.RELOAD_CONF_JOB_NAME, interval, c.ReloadConfig)
		}
		_ = jobContainer.StartJob(constant.RELOAD_CONF_JOB_NAME)
	} else {
		jobContainer.StopJob(constant.RELOAD_CONF_JOB_NAME)
	}
}

func (c *Conf) DestructCallBack() {

}

func (c *Conf) InitCallBack() {
	ConfInit()
}

// 获取配置实例
func GetConfObj(terminalConfs TerminalType) _interface.Conf {
	if _appConf != nil {
		return _appConf
	}
	_appConf = &Conf{
		lock:     &sync.Mutex{},
		ConfKeys: make([]string, 0),
		confs:    make(map[string]*_interface.ConfParam),
	}
	// 注册回调
	RegisterCallBack(_appConf)

	_appConf.Init(terminalConfs)

	// 第一次执行配置更改回调
	_appConf.ChangeConfCallBack()
	return _appConf
}

// 设置ini配置对象
func setIniCfg(data map[string]*_interface.ConfParam) *ini.File {
	cfg := ini.Empty()
	sec, _ := cfg.NewSection("DEFAULT")
	for k, v := range data {
		confVal := "\"" + v.ConfVal + "\""
		_, _ = sec.NewKey(k, confVal)
	}
	return cfg
}
