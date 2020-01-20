package service

var (
	ChangeConfCallBacks []func()
)

func CallBackInit() {
	// 设置配置更改回调
	ChangeConfCallBacks = make([]func(), 0)

	ChangeConfCallBacks = append(ChangeConfCallBacks, GetLogContainerInstance().ChangeConfCallBack)
}


