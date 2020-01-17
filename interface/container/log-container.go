package container

import _interface "github.com/TISUnion/most-simple-mcd/interface"

// log
// 日志接口
type LogContainer interface {
	// GetLogByName
	// 根据名称获取log结构体
	GetLogByName(string) _interface.Log

	// GetLogById
	// 根据id获取log结构体
	GetLogById(int) _interface.Log

	// AddLog
	// 新建一个日志
	// 第一个参数为日志名称
	// 第二个参数为写入日志等级 不传为info
	// 第三个参数为是否显示 不传为false
	// 第四个参数为日志路径默认在log文件夹下
	AddLog(string, ...interface{}) _interface.Log

	// WriteLogOnChannels
	// 第一个参数为日志信息
	// 后面的参数为要写入日志的名称
	WriteLogOnChannels(string, ...string)
}
