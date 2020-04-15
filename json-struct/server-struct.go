package json_struct

// 服务器配置
type ServerConf struct {
	// EntryId
	// 实例唯一id
	EntryId string `json:"id"`

	// Name
	// 服务器名称
	Name string `json:"name"`

	// CmdStr
	// 执行的完整命令
	// 下标为0： 命令名称
	// 大于0为命令参数
	CmdStr []string `json:"cmd_str"`

	// Port
	// 启动服务器端口
	Port int `json:"port"`

	// RunPath
	// 运行所在工作区间
	RunPath string `json:"run_rath"`

	// IsMirror
	// 是否是镜像服务器
	IsMirror bool `json:"is_mirror"`

	// IsStartMonitor
	// 是否启动资源监听器
	IsStartMonitor bool `json:"is_start_monitor"`

	// Memory
	// 使用内存大小，单位M
	Memory int `json:"memory"`

	// Version
	// 服务端版本
	Version string `json:"version"`

	// GameType
	// 服务器模式
	GameType string `json:"game_type"`

	// State
	// 启动状态：0.未启动 1.启动  -1.正在启动 -2.正在关闭
	State int `json:"state"`

	// 本机的ip
	Ips []string `json:"ips"`
}

// 服务端运行状态
type ServerRunState struct {
	// State
	// 启动状态：0.未启动 1.启动  -1.正在启动 -2.正在关闭
	State int `json:"state"`
}

// 服务器接收消息
type ReciveMessage struct {
	Player     string `json:"player"`
	Time       string `json:"time"`
	Speak      string `json:"speak"`
	OriginData []byte `json:"origin_data"`
	ServerId   string `json:"server_id"`
}

// 插件信息
type PluginInfo struct {
	Name            string `json:"name"`
	Id              string `json:"id"`
	IsBan           bool   `json:"is_ban"`
	CommandName     string `json:"command_name"`
	Description     string `json:"description"`
	HelpDescription string `json:"help_description"`
}

// 服务器详情
type ServerDetail struct {
	ServInfo *ServerConf   `json:"server_info"`
	PlgnInfo []*PluginInfo `json:"plugin_info"`
}

// 操作服务器
type OperateServer struct {
	ServerId []string `json:"id"`
	// 操作类型：1. 启动  2. 停止  3.重启
	OperateType int `json:"operate_type"`
}

// 操作插件
type OperatePlugin struct {
	ServerId string   `json:"server_id"`
	PluginId []string `json:"plugin_id"`
	// 操作类型：1. 启动  2. 停止
	OperateType int `json:"operate_type"`
}

// 插件命令行解析对象
type PluginCommand struct {
	Command  string
	Params   []string
	ParamMap map[string]string
}
