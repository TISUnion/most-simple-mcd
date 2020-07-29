package mcdr_plugin_compatible

import (
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	uuid "github.com/satori/go.uuid"
)

const (
	isGlobal = true
)

type McdrPlugin struct {
	id                string
	pluginDescription string
	pluginCommand     string
	helpDescription   string
	pluginPath        string
	packageName       string
	pluginName        string
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
	return isGlobal
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
func (p *McdrPlugin) InitCallBack()       {

}

/* --------------------------------------------- */

/* ---------非全局插件，服务端启动，关闭回调--------- */
func (p *McdrPlugin) Start() {}
func (p *McdrPlugin) Stop()  {}

/* --------------------------------------------- */

func (p *McdrPlugin) HandleMessage(message *models.ReciveMessage) {
	if message.Player == "" {
		return
	}
	com := utils.ParsePluginCommand(message.Speak)
	if com.Command != p.pluginCommand {
		return
	}

	mcServer, err := modules.GetMinecraftServerContainerInstance().GetServerById(message.ServerId)
	if err != nil {
		return
	}

	if len(com.Params) == 0 {
		_ = mcServer.TellrawCommand(message.Player, p.helpDescription)
	} else {
		p.paramsHandle(message.Player, com, mcServer)
	}
}

func (p *McdrPlugin) paramsHandle(player string, pc *models.PluginCommand, mcServer server.MinecraftServer) {
	switch pc.Params[0] {
	// write code...
	default:
		_ = mcServer.TellrawCommand(player, p.helpDescription)
	}
}

func (*McdrPlugin) NewInstance() plugin.Plugin {
	return nil
}

var McdrPluginCompatiblePluginObj plugin.Plugin

func GetMcdrPluginCompatiblePluginInstance() plugin.Plugin {
	if McdrPluginCompatiblePluginObj != nil {
		return McdrPluginCompatiblePluginObj
	}
	McdrPluginCompatiblePluginObj = &McdrPlugin{
		id: uuid.NewV4().String(),
	}
	modules.RegisterCallBack(McdrPluginCompatiblePluginObj)
	return McdrPluginCompatiblePluginObj
}
