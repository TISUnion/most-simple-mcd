package here

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	uuid "github.com/satori/go.uuid"
	"regexp"
)

const (
	pluginName        = "玩家的坐标高亮"
	pluginDescription = "玩家的坐标将被显示并且将被发光效果高亮"
	pluginCommand     = "!!here"
	isGlobal          = false
	helpDescription   = "输入!!here玩家的坐标将被显示并且将被发光效果高亮"

	LOW_VERSION_DESCRIPTION = "该插件只能在1.13+使用"
	LOW_VERSION             = "1.13"
	highlightTime = 15
)

type HerePlugin struct {
	mcServer server.MinecraftServer
	id       string
}

func (p *HerePlugin) GetDescription() string {
	return pluginDescription
}

func (p *HerePlugin) GetHelpDescription() string {
	return helpDescription
}

func (p *HerePlugin) GetCommandName() string {
	return pluginCommand
}

func (p *HerePlugin) IsGlobal() bool {
	return isGlobal
}

func (p *HerePlugin) GetId() string {
	return p.id
}

func (p *HerePlugin) GetName() string {
	return pluginName
}

func (p *HerePlugin) Init(mcServer server.MinecraftServer) {
	p.mcServer = mcServer
}

/* ------------------回调接口-------------------- */
func (p *HerePlugin) ChangeConfCallBack() {}
func (p *HerePlugin) DestructCallBack()   {}
func (p *HerePlugin) InitCallBack()       {}

/* --------------------------------------------- */

/* ---------非全局插件，服务端启动，关闭回调--------- */
func (p *HerePlugin) Start() {}
func (p *HerePlugin) Stop()  {}

/* --------------------------------------------- */

func (p *HerePlugin) HandleMessage(message *models.ReciveMessage) {
	// 匹配是否是坐标json
	p.PosDataHandle([]byte(message.OriginData))

	if !message.IsPlayer {
		return
	}
	if message.Command != pluginCommand {
		return
	}

	// 判断版本是否符合
	if utils.CompareMcVersion(p.mcServer.GetServerConf().Version, LOW_VERSION) == constant.COMPARE_LT {
		_ = p.mcServer.TellrawCommand(message.Player, LOW_VERSION_DESCRIPTION)
		return
	}

	if len(message.Params) == 0 {
		_ = p.mcServer.RunCommand("data", "get", "entity", message.Player)
	} else {
		_ = p.mcServer.TellrawCommand(message.Player, helpDescription)
	}
}

// 获取玩家坐标信息
func (p *HerePlugin) PosDataHandle(originData []byte) {
	reg, _ := regexp.Compile("\\[Server thread/INFO\\]: (.+) has the following entity data: \\{.+Dimension: (-?\\d), Rotation: \\[.+\\[(.+)d, (.+)d, (.+)d\\],.+\\}")
	match := reg.FindStringSubmatch(string(originData))
	if len(match) != 6 {
		return
	}
	player := match[1]
	dimension := match[2]
	posData := match[3:]
	position_show := fmt.Sprintf("[x:%s, y:%s, z:%s]", posData[0], posData[1], posData[2])
	dimension_display := ""
	switch dimension {
	case "0":
		dimension_display = "§2主世界"
	case "-1":
		dimension_display = "§4地狱"
	case "1":
		dimension_display = "§5末地"
	}
	tellMsg := fmt.Sprintf("§e%s§r @ %s §r%s", player, dimension_display, position_show)
	_ = p.mcServer.TellrawCommand(constant.MC_ALL_PLAYER, tellMsg)
	_ = p.mcServer.TellrawCommand(player, fmt.Sprintf("你将会被高亮%d秒", highlightTime))
}

func (*HerePlugin) NewInstance() plugin.Plugin {
	plg := &HerePlugin{
		id: uuid.NewV4().String(),
	}
	modules.RegisterCallBack(plg)
	return plg
}

var HerePluginObj = &HerePlugin{}
