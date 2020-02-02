package service

import _interface "github.com/TISUnion/most-simple-mcd/interface"

var (
	// 组件初始化回调
	initConfCallBacks = make([]func(), 0)

	// 设置配置更改回调容器
	changeConfCallBacks = make([]func(), 0)

	// 退出后释放资源
	destructCallBacks = make([]func(), 0)
)

func RegisterCallBack(callback _interface.CallBack) {
	// 日志配置回调
	changeConfCallBacks = append(changeConfCallBacks, callback.ChangeConfCallBack)

	// 释放资源回调
	destructCallBacks = append(destructCallBacks, callback.ChangeConfCallBack)

	// 初始化回调
	initConfCallBacks = append(initConfCallBacks, callback.InitCallBack)
}

func RunChangeConfCallBacks() {
	for _, callback := range changeConfCallBacks {
		callback()
	}
}

func RunDestructCallBacks() {
	for _, callback := range destructCallBacks {
		callback()
	}
}

func RunInitCallBacks() {
	for _, callback := range initConfCallBacks {
		callback()
	}
}
