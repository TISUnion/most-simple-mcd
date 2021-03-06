package container

import _interface "github.com/TISUnion/most-simple-mcd/interface"

// log
// 日志接口
type LogContainer interface {
	_interface.CallBack

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
	// 第三个参数为日志路径默认在log文件夹下
	AddLog(string, ...string) _interface.Log

	// WriteLogOnChannels
	// 第一个参数为日志信息
	// 第二个参数为日志等级
	// 第三个参数为要写日志名称
	WriteLogOnChannels(string, string, []string)

	// WriteLog
	// 第一个参数为要写的日志内容
	// 第二个参数为日志等级（可选，默认info）
	WriteLog(...string)
}
