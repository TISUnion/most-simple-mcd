package reread_chicken

import (
	plugin_interface "github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	uuid "github.com/satori/go.uuid"
)

type RereadChickenPlugin struct {
	mcServer server.MinecraftServer
	id       string
}

const (
	pluginName        = "复读机插件"
	pluginDescription = "复读机，会复读"
	pluginCommand     = "!!repeat"
	isGlobal          = false
	helpDescription   = "使用!!repeat后会复读"
	help              = "-p"
)

func (r *RereadChickenPlugin) GetDescription() string {
	return pluginDescription
}

func (r *RereadChickenPlugin) GetCommandName() string {
	return pluginCommand
}

func (r *RereadChickenPlugin) ChangeConfCallBack() {}
func (r *RereadChickenPlugin) DestructCallBack()   {}
func (r *RereadChickenPlugin) InitCallBack()       {}
func (r *RereadChickenPlugin) GetId() string       { return r.id }
func (r *RereadChickenPlugin) GetName() string     { return pluginName }
func (r *RereadChickenPlugin) IsGlobal() bool      { return isGlobal }
func (r *RereadChickenPlugin) Start()              {}
func (r *RereadChickenPlugin) Stop()               {}
func (r *RereadChickenPlugin) HandleMessage(message *json_struct.ReciveMessage) {
	if message.Player == "" {
		return
	}
	com := utils.ParsePluginCommand(message.Speak)
	if com.Command != pluginCommand {
		return
	}
	tellMsg := ""
	if len(com.Params) == 0 || com.Params[0] == help {
		tellMsg = helpDescription
	} else {
		tellMsg = com.Params[0]
	}
	_ = r.mcServer.TellCommand(message.Player, tellMsg)
}
func (r *RereadChickenPlugin) Init(mcServer server.MinecraftServer) {
	r.mcServer = mcServer
}
func (r *RereadChickenPlugin) NewInstance() plugin_interface.Plugin {
	p := &RereadChickenPlugin{
		id: uuid.NewV4().String(),
	}
	modules.RegisterCallBack(p)
	return p
}

var RereadChickenPluginObj = &RereadChickenPlugin{}
