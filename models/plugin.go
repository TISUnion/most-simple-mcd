package models

// 插件信息
type PluginInfo struct {
	Name            string `json:"name"`
	Id              string `json:"id"`
	IsBan           bool   `json:"is_ban"`
	CommandName     string `json:"command_name"`
	Description     string `json:"description"`
	HelpDescription string `json:"help_description"`
}

