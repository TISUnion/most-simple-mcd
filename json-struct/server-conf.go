package json_struct

type ServerConf struct {
	// Name
	// 服务器名称
	Name string `json:"name"`

	// CmdStr
	// 执行的完整命令
	// 下标为0： 命令名称
	// 大于0为命令参数
	CmdStr []string `json:"cmdStr"`

	// Port
	// 启动服务器端口
	Port int `json:"port"`

	// RunPath
	// 运行所在工作区间
	RunPath string `json:"runPath"`

	// IsMirror
	// 是否是镜像服务器
	IsMirror bool `json:"isMirror"`

	// IsStartMonitor
	// 是否启动资源监听器
	IsStartMonitor bool `json:"isStartMonitor"`

	// Memory
	// 使用内存大小，单位M
	Memory int `json:"memory"`
}
