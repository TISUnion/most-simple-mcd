package plugins

import (
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/modules"
	uuid "github.com/satori/go.uuid"
)

// tmpl
const (
	pluginName        = "基础插件"
	pluginDescription = "提供最基础的命令"
	pluginCommand     = "!!server"
	isGlobal          = false
	helpDescription   = "\\n!!server help|-l <命令> 帮助信息加具体命令查看命令帮助，不加显示所有命令列表\\n!!server info|-if 查看当前服务端信息\\n!!server infos|-ifs 查看所有服务端信息\\n!!server plugins|-ps 查看插件列表"
)

type BasicPlugin struct {
	id       string
	mcServer server.MinecraftServer
}

func (p *BasicPlugin) ChangeConfCallBack() {}

func (p *BasicPlugin) DestructCallBack() {}

func (p *BasicPlugin) InitCallBack() {}

func (p *BasicPlugin) GetId() string {
	return p.id
}

func (p *BasicPlugin) GetName() string {
	return pluginName
}

func (p *BasicPlugin) GetDescription() string {
	return pluginDescription
}

func (p *BasicPlugin) GetCommandName() string {
	return pluginCommand
}

func (p *BasicPlugin) Init(server server.MinecraftServer) {
	p.mcServer = server
}

func (p *BasicPlugin) IsGlobal() bool {
	return isGlobal
}

func (p *BasicPlugin) NewInstance() plugin.Plugin {
	pl := &BasicPlugin{
		id: uuid.NewV4().String(),
	}
	modules.RegisterCallBack(pl)
	return pl
}

func (p *BasicPlugin) HandleMessage(messageType *json_struct.ReciveMessage) {
	panic("implement me")
}

func (p *BasicPlugin) Start() {}

func (p *BasicPlugin) Stop() {}
