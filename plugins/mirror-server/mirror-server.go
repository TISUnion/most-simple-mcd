package mirror_server

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/container"
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	uuid "github.com/satori/go.uuid"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// tmpl
const (
	pluginName        = "镜像插件"
	pluginDescription = "备份存档，运行备份存档镜像"
	pluginCommand     = "!!mirror"
	isGlobal          = true
	helpDescription   = "\\n使用方式：!!mirror list|-l 查看所有镜像服务器\\n!!mirror save|-s <自定义备份镜像名称> 保存一份当前服务器的镜像\\n!!mirror start|-s <备份镜像id> 开启镜像服务器\\n!!mirror stop|-p <备份镜像id> 关闭镜像服务器"
)

// self
const (
	maxLen        = 8
	MC_MIRROR_DIR = "minecraft-mirrors"
)

var (
	stateMap map[int]string
	listHead []string
)

type MirrorServerPlugin struct {
	mcContainer  container.MinecraftContainer
	mirrors      []server.MinecraftServer
	// 服务端是否保存状态储存 false: 不储存 true：期望储存
	mcSaveState map[string]bool
	// 通知保存完成管道
	savedChan chan struct{}
	id        string
	lock      *sync.Mutex
}

func (p *MirrorServerPlugin) ChangeConfCallBack() {
}

func (p *MirrorServerPlugin) DestructCallBack() {
}

func (p *MirrorServerPlugin) InitCallBack() {
	p.mcSaveState = make(map[string]bool)
	stateMap = make(map[int]string)
	p.savedChan = make(chan struct{})
	p.lock = &sync.Mutex{}
	// 0.未启动 1.启动  -1.正在启动 -2.正在关闭
	stateMap[constant.MC_STATE_STOP] = "未启动"
	stateMap[constant.MC_STATE_START] = "启动"
	stateMap[constant.MC_STATE_STARTIND] = "正在启动"
	stateMap[constant.MC_STATE_STOPING] = "正在关闭"
	listHead = []string{"id", "名称", "内存", "版本", " 状态"}

	// 注册保存回调
	p.mcContainer = modules.GetMinecraftServerContainerInstance()
	p.mcContainer.RegisterAllServerSaveCallback(p.saveCallback)
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
		//_ = mcServer.TellCommand(message.Player, helpDescription)
	} else {
		p.paramsHandle(message.Player, com, mcServer)
	}
}

func (p *MirrorServerPlugin) paramsHandle(player string, pc *json_struct.PluginCommand, mcServer server.MinecraftServer) {
	switch pc.Params[0] {
	case "list", "-l":
		data := make([][]string, 0)
		for _, mcMs := range p.mirrors {
			mcConf := mcMs.GetServerConf()
			data = append(data, []string{utils.Ellipsis(mcConf.EntryId, maxLen), mcConf.Name, strconv.Itoa(mcConf.Memory), mcConf.Version, stateMap[mcConf.State]})
		}
		//_ = mcServer.TellCommand(player, utils.FormateTable(listHead, data))
		fmt.Println(utils.FormateTable(listHead, data))
		_ = mcServer.SayCommand(utils.FormateTable(listHead, data))
	case "save", "-s":
		if len(pc.Params) < 2 {
			return
		}
		name := pc.Params[1]
		mcServerConf := mcServer.GetServerConf()
		p.saveServer(mcServerConf.EntryId, mcServer)
		mirrorId := uuid.NewV4().String()
		runPath, ok := p.buildMirror(mcServerConf, mirrorId)
		if !ok {
			return
		}
		mirrorSrvConf := &json_struct.ServerConf{
			EntryId:  mirrorId,
			Name:     name,
			CmdStr:   utils.GetCommandArr(constant.MC_DEFAULT_MEMORY, runPath),
			RunPath:  runPath,
			IsMirror: true,
			Memory:   constant.MC_DEFAULT_MEMORY,
		}
		p.mcContainer.AddServer(mirrorSrvConf, true)
	default:
		//_ = mcServer.TellCommand(player, helpDescription)
		_ = mcServer.SayCommand(helpDescription)
	}
}

// 保存服务端（save-all）
func (p *MirrorServerPlugin) saveServer(id string, mcServer server.MinecraftServer) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.mcSaveState[id] = true
	_ = mcServer.Command("/save-all")
	<-p.savedChan
}

func (p *MirrorServerPlugin) saveCallback(id string) {
	if p.mcSaveState[id] {
		p.savedChan <- struct{}{}
		p.mcSaveState[id] = false
	}
}

// 构建镜像
func (p *MirrorServerPlugin) buildMirror(conf *json_struct.ServerConf, id string) (path string, ok bool) {
	serverPath, filename := filepath.Split(conf.RunPath)
	mirrorPath := filepath.Join(modules.GetConfVal(constant.WORKSPACE), MC_MIRROR_DIR, id)
	mirrorRunPath := filepath.Join(mirrorPath, id+constant.JAR_SUF)
	oldFilePath := filepath.Join(mirrorPath, filename)
	err := utils.CopyDir(serverPath, mirrorPath)
	if err != nil {
		modules.WriteLogToDefault(fmt.Sprintf("构建镜像失败: %v", err))
		return "", false
	}
	err = os.Rename(oldFilePath, mirrorRunPath)
	if err != nil {
		modules.WriteLogToDefault(fmt.Sprintf("构建镜像失败: %v", err))
		return "", false
	}
	return mirrorRunPath, true
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
