package service

import (
	_interface "github.com/TISUnion/most-simple-mcd/interface"
)

var (
	// 设置配置更改回调容器
	changeConfCallBacks = make([]func(), 0)

	// 退出后释放资源
	destructCallBacks = make([]func(), 0)
)

func RegisterCallBack(callback _interface.CallBack) {

	// 日志配置回调
	changeConfCallBacks = append(changeConfCallBacks, callback.ChangeConfCallBack)

	// 释放资源回调
	destructCallBacks = append(destructCallBacks, callback.DestructCallBack)

	// 运行初始化回调
	callback.InitCallBack()
}

func RunChangeConfCallBacks() {
	for _, callback := range changeConfCallBacks {
		callback()
	}
}

func RunDestructCallBacks() {
	// 反转数组
	// 这里的销毁回调，必须倒序执行，否则依赖销毁，会报错
	length := len(destructCallBacks)
	for i := 0; i < length/2; i++ {
		temp := destructCallBacks[length-1-i]
		destructCallBacks[length-1-i] = destructCallBacks[i]
		destructCallBacks[i] = temp
	}

	for _, callback := range destructCallBacks {
		callback()
	}
}
