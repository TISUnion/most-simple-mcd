package _interface

type Conf interface {
	// GetConfig
	// 获取所有配置
	GetConfig() map[string]string

	// GetConfigKeys
	// 获取所有配置的键值
	GetConfigKeys() []string

	// GetConfVal
	// 获取单个配置
	GetConfVal(string) string

	// reloadConfig
	// 重载配置
	ReloadConfig()
}
