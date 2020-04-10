package mirror_server

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	uuid "github.com/satori/go.uuid"
)

const (
	pluginName        = "镜像插件"
	pluginDescription = "备份存档，运行备份存档镜像"
	pluginCommand     = "!!mirror"
	isGlobal          = true
	helpDescription   = "使用!!repeat后会复读"
	help              = "help"
)

type MirrorServerPlugin struct {
	mirrors      []server.MinecraftServer
	startMirrors []server.MinecraftServer
	stopMirrors  []server.MinecraftServer
	id           string
}

func (p *MirrorServerPlugin) ChangeConfCallBack() {
}

func (p *MirrorServerPlugin) DestructCallBack() {
}

func (p *MirrorServerPlugin) InitCallBack() {
}

func (p *MirrorServerPlugin) GetId() string {
	return p.id
}

func (p *MirrorServerPlugin) GetName() string {
	return pluginName
}

func (p *MirrorServerPlugin) GetDescription() string {
	return pluginDescription
}

func (p *MirrorServerPlugin) GetCommandName() string {
	return pluginCommand
}

func (p *MirrorServerPlugin) IsGlobal() bool {
	return isGlobal
}

func (p *MirrorServerPlugin) NewInstance() plugin.Plugin { return nil }

func (p *MirrorServerPlugin) HandleMessage(message *json_struct.ReciveMessage) {
	if message.Player == "" {
		return
	}
	com := utils.ParsePluginCommand(message.Speak)
	if com.Command != pluginCommand {
		return
	}
	mcServer, ok := modules.GetMinecraftServerContainerInstance().GetServerById(message.ServerId)
	if !ok {
		return
	}
	if len(com.Params) == 0 || com.Params[0] == help {
		_ = mcServer.Command(fmt.Sprintf("/tell %s %s"))
	} else {

	}
}

// -------------非全局插件需实现方法--------------
func (p *MirrorServerPlugin) Init(server server.MinecraftServer) {}

func (p *MirrorServerPlugin) Start() {}

func (p *MirrorServerPlugin) Stop() {}
// --------------------------------------------

var mirrorServerPluginObj plugin.Plugin

func GetMirrorServerPluginInstance() plugin.Plugin {
	if mirrorServerPluginObj != nil {
		return mirrorServerPluginObj
	}
	mirrorServerPluginObj = &MirrorServerPlugin{
		id: uuid.NewV4().String(),
	}
	modules.RegisterCallBack(mirrorServerPluginObj)
	return mirrorServerPluginObj
}
