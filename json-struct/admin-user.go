package json_struct

type AdminUser struct {
	Nickname string   `json:"nickname"`
	Account  string   `json:"account"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
	Avatar   string   `json:"avatar"`
}

type UserToken struct {
	Token string `json:"token"`
}
