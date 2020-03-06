package _interface

import json_struct "github.com/TISUnion/most-simple-mcd/json-struct"

type Conf interface {
	CallBack

	// GetConfig
	// 获取所有配置
	GetConfig() map[string]string

	// GetConfigObj
	// 获取所有配置对象
	GetConfigObj()  map[string]*json_struct.ConfParam

	// GetConfigKeys
	// 获取所有配置的键值
	GetConfigKeys() []string

	// GetConfVal
	// 获取单个配置
	GetConfVal(string) string

	// SetConfig
	// 设置一个配置变量（若存在则覆盖, 不存在不会创建，需要先注册）
	SetConfig(string, string)

	// ReloadConfig
	// 重载配置
	ReloadConfig()

	// RegisterConfParam
	// 注册配置，已拥有，则不覆盖
	RegisterConfParam(Name, ConfVal, description string, level int)
}
