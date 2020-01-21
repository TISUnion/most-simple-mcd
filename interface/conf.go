package _interface

type ConfParam struct {
	ConfVal string
	DefaultConfVal string
	Name string
	Level int
	Description string
}

type Conf interface {
	CallBack

	// GetConfig
	// 获取所有配置
	GetConfig() map[string]string

	// GetConfigKeys
	// 获取所有配置的键值
	GetConfigKeys() []string

	// GetConfVal
	// 获取单个配置
	GetConfVal(string) string

	// SetConfig
	// 设置一个配置变量（若存在则覆盖）
	SetConfig(string, string)

	// ReloadConfig
	// 重载配置
	ReloadConfig()
}
