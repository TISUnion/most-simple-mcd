package service

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/utils"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	_appConf         *Conf
	DefaultConfParam map[string]*json_struct.ConfParam
)

// Conf
// 导入配置优先级：管理后台设置配置 > 命令行变量 > 环境变量 > 配置文件 > 默认配置
type Conf struct {
	// Confs
	// 配置
	confs map[string]*json_struct.ConfParam

	// ConfKeys
	// 所有配置键值
	ConfKeys []string

	// lock
	// 读写锁
	lock *sync.Mutex
}

func (c *Conf) GetConfigObj() map[string]*json_struct.ConfParam {
	return c.confs
}

func ConfInit() {
	// 注册默认配置
	DefaultConfParam = make(map[string]*json_struct.ConfParam)
	RegisterConfig(constant.IS_RELOAD_CONF, "false", constant.IS_RELOAD_CONF_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	RegisterConfig(constant.RELOAD_CONF_INTERVAL, "5000", constant.RELOAD_CONF_INTERVAL_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	RegisterConfig(constant.IS_MANAGE_HTTP, "true", constant.IS_MANAGE_HTTP_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	RegisterConfig(constant.MANAGE_HTTP_SERVER_PORT, "80", constant.MANAGE_HTTP_SERVER_PORT_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	RegisterConfig(constant.LOG_SAVE_INTERVAL, constant.LOG_SAVE_INTERVAL_TWICEDAY, constant.LOG_SAVE_INTERVAL_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	if workspace, err := utils.GetCurrentPath(); err == nil {
		RegisterConfig(constant.WORKSPACE, workspace, constant.WORKSPACE_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	} else {
		RegisterConfig(constant.WORKSPACE, "./", constant.WORKSPACE_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	}
	RegisterConfig(constant.LOG_PATH, filepath.Join(DefaultConfParam[constant.WORKSPACE].ConfVal, "logs"), constant.LOG_PATH_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	RegisterConfig(constant.CONF_PATH, filepath.Join(DefaultConfParam[constant.WORKSPACE].ConfVal, "conf/mcd.ini"), constant.CONF_PATH_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	RegisterConfig(constant.IS_AUTO_CHANGE_MC_SERVER_REPEAT_PORT, "true", constant.IS_AUTO_CHANGE_MC_SERVER_REPEAT_PORT_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	RegisterConfig(constant.MONITOR_INTERVAL, "2s", constant.MONITOR_INTERVAL_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
	RegisterConfig(constant.I18N, "zh", constant.I18N_DESCREPTION, constant.CONF_DEFAULT_LEVEL)
}

// loadFilePath
// 获取配置文件目录
func (c *Conf) loadFilePath(terminalConfs map[string]*string) {
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

// reloadConfig
// 重新再加配置
func (c *Conf) ReloadConfig() {
	c.lock.Lock()
	defer c.lock.Unlock()
	// 加载配置文件
	c.loadFileConf()
	// 加载环境变量
	c.loadEnvConf()
	// 执行配置更改回调
	RunChangeConfCallBacks()
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
func (c *Conf) loadTerminalConf(terminalConfs map[string]*string) {
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

func (c *Conf) Init(terminalConfs map[string]*string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// 加载文件配置文件路径
	c.loadFilePath(terminalConfs)

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
			confParam.Description = description
		}
	}
}

// 注册配置
func (c *Conf) RegisterConfParam(Name, ConfVal, description string, level int) {
	if _, ok := c.confs[Name]; !ok {
		c.ConfKeys = append(c.ConfKeys, Name)
		c.confs[Name] = &json_struct.ConfParam{
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
	terminalConfs := InitFlag()
	c.Init(terminalConfs)
	c.ChangeConfCallBack()
}

// 获取配置实例
func GetConfInstance() _interface.Conf {
	if _appConf != nil {
		return _appConf
	}
	_appConf = &Conf{
		lock:     &sync.Mutex{},
		ConfKeys: make([]string, 0),
		confs:    make(map[string]*json_struct.ConfParam),
	}
	// 注册回调
	RegisterCallBack(_appConf)
	return _appConf
}

// 设置ini配置对象
func setIniCfg(data map[string]*json_struct.ConfParam) *ini.File {
	cfg := ini.Empty()
	sec, _ := cfg.NewSection(ini.DefaultSection)
	for k, v := range data {
		confVal := "\"" + v.ConfVal + "\""
		_, _ = sec.NewKey(k, confVal)
	}
	return cfg
}

// 获取配置值帮助函数
func GetConfVal(confKey string) string {
	return GetConfInstance().GetConfVal(confKey)
}

func RegisterConfig(Name, ConfVal, description string, level int) {
	conObj := GetConfInstance()
	DefaultConfParam[Name] = utils.NewConfParam(Name, ConfVal, description, level)
	conObj.RegisterConfParam(Name, ConfVal, description, level)
}
