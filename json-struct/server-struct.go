package json_struct

// 服务器配置
type ServerConf struct {
	// EntryId
	// 实例唯一id
	EntryId string

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

	// Version
	// 服务端版本
	Version string `json:"version"`

	// GameType
	// 服务器模式
	GameType string `json:"gameType"`
}

// 服务器接收消息
type ReciveMessageType struct {
	Player     string `json:"player"`
	Time       string `json:"time"`
	Speak      string `json:"speak"`
	OriginData []byte `json:"origin_data"`
	ServerId   string `json:"server_id"`
}
