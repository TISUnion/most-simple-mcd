package json_struct

// 用户信息
type AdminUser struct {
	Nickname string   `json:"nickname"`
	Account  string   `json:"account"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
	Avatar   string   `json:"avatar"`
}

// 登陆token
type UserToken struct {
	Token string `json:"token"`
}

// 执行操作返回结果
type OperateResult struct {
	Status uint8 `json:"status"`
}

// 修改配置结构
type Config struct {
	ConfVal        string `json:"config_val"`
	ConfKey        string `json:"config_key"`
}

// 配置信息
type ConfParam struct {
	ConfVal        string `json:"config_val"`
	DefaultConfVal string `json:"-"`
	Name           string `json:"config_key"`
	Level          int    `json:"level"`
	Description    string `json:"description"`
}