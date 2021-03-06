package mcdr_plugin_compatible

import (
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/modules"
	uuid "github.com/satori/go.uuid"
)

const (
	isMcdrPluginGlobal = true
)

type McdrPlugin struct {
	id                string
	pluginDescription string
	pluginCommand     string
	helpDescription   string
	pluginPath        string
	packageName       string
	pluginName        string
	CArrIndex         int // c中插件数组索引
}

func (p *McdrPlugin) GetDescription() string {
	return p.pluginDescription
}

func (p *McdrPlugin) GetHelpDescription() string {
	return p.helpDescription
}

func (p *McdrPlugin) GetCommandName() string {
	return p.pluginCommand
}

func (p *McdrPlugin) IsGlobal() bool {
	return isMcdrPluginGlobal
}

func (p *McdrPlugin) GetId() string {
	return p.id
}

func (p *McdrPlugin) GetName() string {
	return p.pluginName
}

func (p *McdrPlugin) Init(mcServer server.MinecraftServer) {

}

/* ------------------回调接口-------------------- */
func (p *McdrPlugin) ChangeConfCallBack() {}
func (p *McdrPlugin) DestructCallBack()   {}
func (p *McdrPlugin) InitCallBack() {

}

/* --------------------------------------------- */

/* ---------非全局插件，服务端启动，关闭回调--------- */
func (p *McdrPlugin) Start() {}
func (p *McdrPlugin) Stop()  {}

/* --------------------------------------------- */

func (p *McdrPlugin) HandleMessage(message *models.ReciveMessage) {
	if !message.IsPlayer {
		return
	}
	if message.Command != p.pluginCommand {
		return
	}

	mcServer, err := modules.GetMinecraftServerContainerInstance().GetServerById(message.ServerId)
	if err != nil {
		return
	}

	if len(message.Params) == 0 {
		_ = mcServer.TellrawCommand(message.Player, p.helpDescription)
	} else {
		p.paramsHandle(message.Player, message, mcServer)
	}
}

func (p *McdrPlugin) paramsHandle(player string, pc *models.ReciveMessage, mcServer server.MinecraftServer) {
	switch pc.Params[0] {
	// write code...
	default:
		_ = mcServer.TellrawCommand(player, p.helpDescription)
	}
}

func (*McdrPlugin) NewInstance() plugin.Plugin {
	return nil
}

var McdrPluginObj plugin.Plugin

func GetMcdrPluginInstance() plugin.Plugin {
	if McdrPluginObj != nil {
		return McdrPluginObj
	}
	McdrPluginObj = &McdrPlugin{
		id: uuid.NewV4().String(),
	}
	modules.RegisterCallBack(McdrPluginObj)
	return McdrPluginObj
}
