package mirror_server

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/container"
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	uuid "github.com/satori/go.uuid"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const (
	pluginName        = "镜像插件"
	pluginDescription = "备份存档，运行备份存档镜像"
	pluginCommand     = "!!mirror"
	isGlobal          = true
	helpDescription   = "使用方式：!!mirror list|-l 查看所有镜像服务器\n!!mirror save|-s <自定义备份镜像名称> 保存一份当前服务器的镜像\n!!mirror start|-st <备份镜像id> 开启镜像服务器\n!!mirror stop|-sp <备份镜像id> 关闭镜像服务器"

	maxLen        = 5
	MC_MIRROR_DIR = "minecraft-mirrors"
)

var (
	stateMap map[int64]string
	listHead []string
)

type MirrorServerPlugin struct {
	mcContainer container.MinecraftContainer
	mirrors     []server.MinecraftServer
	// 服务端是否保存状态储存 false: 不储存 true：期望储存
	mcSaveState map[string]bool
	// 通知保存完成管道
	savedChan chan struct{}
	id        string
	lock      *sync.Mutex
}

func (p *MirrorServerPlugin) GetDescription() string {
	return pluginDescription
}

func (p *MirrorServerPlugin) GetHelpDescription() string {
	return helpDescription
}

func (p *MirrorServerPlugin) GetCommandName() string {
	return pluginCommand
}

func (p *MirrorServerPlugin) IsGlobal() bool {
	return isGlobal
}

func (p *MirrorServerPlugin) GetId() string {
	return p.id
}

func (p *MirrorServerPlugin) GetName() string {
	return pluginName
}

func (p *MirrorServerPlugin) Init(mcServer server.MinecraftServer) {

}

/* ------------------回调接口-------------------- */
func (p *MirrorServerPlugin) ChangeConfCallBack() {}
func (p *MirrorServerPlugin) DestructCallBack()   {}
func (p *MirrorServerPlugin) InitCallBack() {
	p.mcSaveState = make(map[string]bool)
	stateMap = make(map[int64]string)
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

	// 获取已有镜像
	p.getMirrors()
}

/* --------------------------------------------- */

/* ---------非全局插件，服务端启动，关闭回调--------- */
func (p *MirrorServerPlugin) Start() {}
func (p *MirrorServerPlugin) Stop()  {}

/* --------------------------------------------- */

func (p *MirrorServerPlugin) HandleMessage(message *models.ReciveMessage) {
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

func (p *MirrorServerPlugin) paramsHandle(player string, pc *models.ReciveMessage, mcServer server.MinecraftServer) {
	switch pc.Params[0] {
	case "list", "-l":
		data := make([][]string, 0)
		for _, mcMs := range p.mirrors {
			mcConf := mcMs.GetServerConf()
			data = append(data, []string{utils.Ellipsis(mcConf.EntryId, maxLen), mcConf.Name, strconv.FormatInt(mcConf.Memory, 10), mcConf.Version, stateMap[mcConf.State]})
		}
		_ = mcServer.TellrawCommand(player, utils.FormateTable(listHead, data))
	case "save", "-s":
		if len(pc.Params) < 2 {
			_ = mcServer.TellrawCommand(player, "缺少备份名称！")
			return
		}
		_ = mcServer.TellrawCommand(player, "开始备份......")
		name := pc.Params[1]
		mcServerConf := mcServer.GetServerConf()
		p.saveServer(mcServerConf.EntryId, mcServer)
		mirrorId := uuid.NewV4().String()
		runPath, ok := p.buildMirrorPath(mcServerConf, mirrorId)
		if !ok {
			return
		}
		mirrorSrvConf := &models.ServerConf{
			EntryId:  mirrorId,
			Name:     name,
			CmdStr:   utils.GetCommandArr(constant.MC_DEFAULT_MEMORY, runPath),
			RunPath:  runPath,
			IsMirror: true,
			Memory:   constant.MC_DEFAULT_MEMORY,
			Side:     mcServerConf.Side,
		}
		p.mcContainer.AddServer(mirrorSrvConf, true)
		p.getMirrors()
		_ = mcServer.TellrawCommand(player, name+" 备份完成！")
	case "start", "-st":
		if len(pc.Params) < 2 {
			_ = mcServer.TellrawCommand(player, "缺少启动服务端id！")
			return
		}
		id := pc.Params[1]
		if mirrorSvr, err := p.mcContainer.GetMirrorServerById(id); err != nil {
			if err == modules.REPEAT_ID {
				_ = mcServer.TellrawCommand(player, "请输入完整的ID！")
			} else {
				_ = mcServer.TellrawCommand(player, "ID不存在！")
			}
			return
		} else {
			_ = p.mcContainer.StartById(mirrorSvr.GetServerEntryId())
			_ = mcServer.TellrawCommand(player, "启动成功，可通过-l查看服务端是否完成启动")
		}
	case "stop", "-sp":
		if len(pc.Params) < 2 {
			_ = mcServer.TellrawCommand(player, "缺少启动服务端id！")
			return
		}
		id := pc.Params[1]
		if mirrorSvr, err := p.mcContainer.GetMirrorServerById(id); err != nil {
			if err == modules.REPEAT_ID {
				_ = mcServer.TellrawCommand(player, "请输入完整的ID！")
			} else {
				_ = mcServer.TellrawCommand(player, "ID不存在！")
			}
			return
		} else {
			_ = p.mcContainer.StopById(mirrorSvr.GetServerEntryId())
			_ = mcServer.TellrawCommand(player, "关闭成功，可通过-l查看服务端是否完成关闭")
		}
	default:
		_ = mcServer.TellrawCommand(player, helpDescription)
	}
}

// 构建镜像路径
func (p *MirrorServerPlugin) buildMirrorPath(conf *models.ServerConf, id string) (path string, ok bool) {
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

func (p *MirrorServerPlugin) getMirrors() {
	allMcSrv := p.mcContainer.GetAllServerObj()
	for _, mcMS := range allMcSrv {
		if mcMS.GetServerConf().IsMirror {
			p.mirrors = append(p.mirrors, mcMS)
		}
	}
}

func (*MirrorServerPlugin) NewInstance() plugin.Plugin {
	return nil
}

var MirrorServerPluginObj plugin.Plugin

func GetMirrorServerPluginInstance() plugin.Plugin {
	if MirrorServerPluginObj != nil {
		return MirrorServerPluginObj
	}
	MirrorServerPluginObj = &MirrorServerPlugin{
		id: uuid.NewV4().String(),
	}
	modules.RegisterCallBack(MirrorServerPluginObj)
	return MirrorServerPluginObj
}
