package mirror_server

import (
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"sync"
)

// tmpl
const (
	pluginName        = "镜像插件"
	pluginDescription = "备份存档，运行备份存档镜像"
	pluginCommand     = "!!mirror"
	isGlobal          = true
	helpDescription   = "!!mirror list|-l 查看所有镜像服务器\\n!!mirror save|-s <自定义备份镜像名称> 保存一份当前服务器的镜像\\n!!mirror start|-s <备份镜像id> 开启镜像服务器\\n!!mirror stop|-p <备份镜像id> 关闭镜像服务器"
)

// self
const (
	maxLen = 8
)

var (
	stateMap map[int]string
	listHead []string
)

type MirrorServerPlugin struct {
	mirrors      []server.MinecraftServer
	startMirrors []server.MinecraftServer
	stopMirrors  []server.MinecraftServer
	mcSaveState  map[string]bool // 服务端是否保存状态储存
	id           string
	lock         *sync.Mutex
}

func (p *MirrorServerPlugin) ChangeConfCallBack() {
}

func (p *MirrorServerPlugin) DestructCallBack() {
}

func (p *MirrorServerPlugin) InitCallBack() {
	p.mcSaveState = make(map[string]bool)
	stateMap = make(map[int]string)
	// 0.未启动 1.启动  -1.正在启动 -2.正在关闭
	stateMap[0] = "未启动"
	stateMap[1] = "启动"
	stateMap[-1] = "正在启动"
	stateMap[-2] = "正在关闭"
	listHead = []string{"id", "名称", "内存", "版本", " 状态"}
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
	if len(com.Params) == 0 {
		_ = mcServer.TellCommand(message.Player, helpDescription)
	} else {
		p.paramsHandle(message.Player, com, mcServer)
	}
}

// TODO
func (p *MirrorServerPlugin) paramsHandle(player string, pc *json_struct.PluginCommand, mcServer server.MinecraftServer) {
	switch pc.Params[0] {
	case "list", "-l":
		data := make([][]string, 0)
		for _, mcMs := range p.mirrors {
			mcConf := mcMs.GetServerConf()
			data = append(data, []string{utils.Ellipsis(mcConf.EntryId, maxLen), mcConf.Name, strconv.Itoa(mcConf.Memory), mcConf.Version, stateMap[mcConf.State]})
		}
		_ = mcServer.TellCommand(player, utils.FormateTable(listHead, data))
	case "save", "-s":
		if len(pc.Params) < 2 {
			return
		}
		//name := pc.Params[1]


	default:
		_ = mcServer.TellCommand(player, helpDescription)
	}
}

//func (p *MirrorServerPlugin) joinListContent() string {
//	res := ""
//	for _, mcS := range p.mirrors {
//
//	}
//}

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
