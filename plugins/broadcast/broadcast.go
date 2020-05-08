package broadcast

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
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
		_ = mcServer.TellrawCommand(message.Player, helpDescription)
	} else {
		p.paramsHandle(message.Player, com, mcServer)
	}
}

func (p *BroadcastPlugin) paramsHandle(player string, pc *models.PluginCommand, mcServer server.MinecraftServer) {
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

var BroadcastPluginObj plugin.Plugin

func GetBroadcastPluginInstance() plugin.Plugin {
	if BroadcastPluginObj != nil {
		return BroadcastPluginObj
	}
	BroadcastPluginObj = &BroadcastPlugin{
		id: uuid.NewV4().String(),
	}
	modules.RegisterCallBack(BroadcastPluginObj)
	return BroadcastPluginObj
}
