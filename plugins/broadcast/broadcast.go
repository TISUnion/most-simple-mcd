package broadcast

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/modules"
	uuid "github.com/satori/go.uuid"
)

const (
	pluginName        = "全服广播插件"
	pluginDescription = "全服广播，使你成为全服最靓的仔"
	pluginCommand     = "!!broadcast"
	isGlobal          = true
	helpDescription   = "使用方式：!!broadcast <广播内容>"
)

type BroadcastPlugin struct {
	id string
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

func (p *BroadcastPlugin) GetId() string {
	return p.id
}

func (p *BroadcastPlugin) GetName() string {
	return pluginName
}

func (p *BroadcastPlugin) Init(mcServer server.MinecraftServer) {

}

/* ------------------回调接口-------------------- */
func (p *BroadcastPlugin) ChangeConfCallBack() {}
func (p *BroadcastPlugin) DestructCallBack()   {}
func (p *BroadcastPlugin) InitCallBack()       {}

/* --------------------------------------------- */

/* ---------非全局插件，服务端启动，关闭回调--------- */
func (p *BroadcastPlugin) Start() {}
func (p *BroadcastPlugin) Stop()  {}

/* --------------------------------------------- */

func (p *BroadcastPlugin) HandleMessage(message *models.ReciveMessage) {
	if !message.IsPlayer {
		return
	}
	if message.Command != pluginCommand {
		return
	}
	mcServer, err := modules.GetMinecraftServerContainerInstance().GetServerById(message.ServerId)
	if err != nil {
		return
	}
	if len(message.Params) == 0 {
		_ = mcServer.TellrawCommand(message.Player, helpDescription)
	} else {
		p.paramsHandle(message.Player, message, mcServer)
	}
}

func (p *BroadcastPlugin) paramsHandle(player string, pc *models.ReciveMessage, mcServer server.MinecraftServer) {
	switch pc.Params[0] {
	case "help", "-h":
		_ = mcServer.TellrawCommand(player, helpDescription)
	default:
		broadcast := fmt.Sprintf("[%s]<%s> %s", mcServer.GetServerConf().Name, player, pc.Params[0])
		ctr := modules.GetMinecraftServerContainerInstance()
		aMcSrv := ctr.GetAllServerObj()
		for _, mcS := range aMcSrv {
			_ = mcS.TellrawCommand(constant.MC_ALL_PLAYER, broadcast)
		}
	}
}

func (*BroadcastPlugin) NewInstance() plugin.Plugin {
	return nil
}

var broadcastPluginObj plugin.Plugin

func GetBroadcastPluginInstance() plugin.Plugin {
	if broadcastPluginObj != nil {
		return broadcastPluginObj
	}
	broadcastPluginObj = &BroadcastPlugin{
		id: uuid.NewV4().String(),
	}
	modules.RegisterCallBack(broadcastPluginObj)
	return broadcastPluginObj
}
