package broadcast

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
)

// tmpl
const (
	pluginName        = "全服广播插件"
	pluginDescription = "全服广播，使你成为全服最靓的仔"
	pluginCommand     = "!!broadcast"
	isGlobal          = true
	helpDescription   = "\\n使用方式：!!broadcast <广播内容>"
)

type BroadcastPlugin struct {
	id string
}

func (p *BroadcastPlugin) ChangeConfCallBack() {}

func (p *BroadcastPlugin) DestructCallBack() {}

func (p *BroadcastPlugin) InitCallBack() {}

func (p *BroadcastPlugin) GetId() string {
	return p.id
}
func (p *BroadcastPlugin) GetName() string {
	return pluginName
}

func (p *BroadcastPlugin) GetDescription() string {
	return pluginDescription
}

func (p *BroadcastPlugin) GetHelpDescription() string {
	return helpDescription
}

func (p *BroadcastPlugin) GetCommandName() string {
	return pluginCommand
}

func (p *BroadcastPlugin) IsGlobal() bool {
	return isGlobal
}
func (p *BroadcastPlugin) Init(server server.MinecraftServer) {}

func (p *BroadcastPlugin) NewInstance() plugin.Plugin { return nil }

func (p *BroadcastPlugin) HandleMessage(message *json_struct.ReciveMessage) {
	if message.Player == "" {
		return
	}
	com := utils.ParsePluginCommand(message.Speak)
	if com.Command != pluginCommand {
		return
	}
	mcServer, err := modules.GetMinecraftServerContainerInstance().GetServerById(message.ServerId)
	if err != nil {
		return
	}
	if len(com.Params) == 0 {
		_ = mcServer.TellCommand(message.Player, helpDescription)
	} else {
		p.paramsHandle(message.Player, com, mcServer)
	}
}

func (p *BroadcastPlugin) Start() {}

func (p *BroadcastPlugin) Stop() {}

func (p *BroadcastPlugin) paramsHandle(player string, pc *json_struct.PluginCommand, mcServer server.MinecraftServer) {
	switch pc.Params[0] {
	case "help", "-h":
		_ = mcServer.TellCommand(player, helpDescription)
	default:
		broadcast := fmt.Sprintf("[%s]<%s> %s", mcServer.GetServerConf().Name, player, pc.Params[0])
		ctr := modules.GetMinecraftServerContainerInstance()
		aMcSrv := ctr.GetAllServerObj()
		for _, mcS := range aMcSrv {
			_ = mcS.SayCommand(broadcast)
		}
	}
}

var BroadcastPluginObj = &BroadcastPlugin{}
