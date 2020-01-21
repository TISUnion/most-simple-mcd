package service

var (
	ChangeConfCallBacks []func()
)

func CallBackInit() {
	// 设置配置更改回调容器
	ChangeConfCallBacks = make([]func(), 0)

	// 日志配置回调
	ChangeConfCallBacks = append(ChangeConfCallBacks, GetLogContainerInstance().ChangeConfCallBack)
}


