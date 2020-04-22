package reread_chicken

import (
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	uuid "github.com/satori/go.uuid"
)

const (
	pluginName        = "复读机插件"
	pluginDescription = "复读机，会复读"
	pluginCommand     = "!!repeat"
	isGlobal          = false
	helpDescription   = "使用!!repeat <内容> 后会复读内容"
)

type RereadChickenPlugin struct {
	mcServer server.MinecraftServer
	id       string
}

func (p *RereadChickenPlugin) GetDescription() string {
	return pluginDescription
}

func (p *RereadChickenPlugin) GetHelpDescription() string {
	return helpDescription
}

func (p *RereadChickenPlugin) GetCommandName() string {
	return pluginCommand
}

func (p *RereadChickenPlugin) IsGlobal() bool {
	return isGlobal
}

func (p *RereadChickenPlugin) GetId() string {
	return p.id
}

func (p *RereadChickenPlugin) GetName() string {
	return pluginName
}

func (p *RereadChickenPlugin) Init(mcServer server.MinecraftServer) {
	p.mcServer = mcServer
}

/* ------------------回调接口-------------------- */
func (p *RereadChickenPlugin) ChangeConfCallBack() {}
func (p *RereadChickenPlugin) DestructCallBack()   {}
func (p *RereadChickenPlugin) InitCallBack()       {}

/* --------------------------------------------- */

/* ---------非全局插件，服务端启动，关闭回调--------- */
func (p *RereadChickenPlugin) Start() {}
func (p *RereadChickenPlugin) Stop()  {}

/* --------------------------------------------- */

func (p *RereadChickenPlugin) HandleMessage(message *json_struct.ReciveMessage) {
	if message.Player == "" {
		return
	}
	com := utils.ParsePluginCommand(message.Speak)
	if com.Command != pluginCommand {
		return
	}

	if len(com.Params) == 0 {
		_ = p.mcServer.TellrawCommand(message.Player, helpDescription)
	} else {
		p.paramsHandle(message.Player, com)
	}
}

func (p *RereadChickenPlugin) paramsHandle(player string, pc *json_struct.PluginCommand) {
	switch pc.Params[0] {
	case "help", "-h":
		_ = p.mcServer.TellrawCommand(player, helpDescription)
	default:
		_ = p.mcServer.TellrawCommand(player, pc.Params[0])
	}
}

func (*RereadChickenPlugin) NewInstance() plugin.Plugin {
	plg := &RereadChickenPlugin{
		id: uuid.NewV4().String(),
	}
	modules.RegisterCallBack(plg)
	return plg
}

var RereadChickenPluginObj = &RereadChickenPlugin{}
