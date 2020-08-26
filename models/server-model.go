package models

// 服务端运行状态
type ServerRunState struct {
	// State
	// 启动状态：0.未启动 1.启动  -1.正在启动 -2.正在关闭
	State int `json:"state"`
}

// 服务器接收消息
type ReciveMessage struct {
	Player       string   `json:"player"`
	Time         string   `json:"time"`         // 时间字符串 H:i:s
	Speak        string   `json:"speak"`        // 玩家发言
	OriginData   string   `json:"origin_data"`  // 服务端原始输出
	ServerId     string   `json:"server_id"`    // 服务端id
	Hour         int      `json:"hour"`         // 时间：小时
	Minute       int      `json:"minute"`       // 时间：分钟
	Second       int      `json:"second"`       // 时间：秒
	Content      string   `json:"content"`      // 事件内容
	Source       int      `json:"source"`       // 命令台为1，服务端输出为2
	LoggingLevel string   `json:"LoggingLevel"` // 日志等级
	IsPlayer     bool     `json:"IsPlayer"`     // 是否有玩家
	IsUser       bool     `json:"isUser"`		// 是否有玩家，或者为命令行输入
	Command      string   `json:"command"`		// 命令
	Params       []string `json:"params"`		// 命令参数
	Event		 int   `json:"event"`		// 事件类型
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
