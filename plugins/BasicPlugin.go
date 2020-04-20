package plugins

import (
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	uuid "github.com/satori/go.uuid"
	"strconv"
)

// tmpl
const (
	pluginName        = "基础插件"
	pluginDescription = "提供最基础的命令"
	pluginCommand     = "!!server"
	isGlobal          = false
	helpDescription   = "!!server help|-l <命令> 帮助信息加具体命令查看命令帮助，不加显示所有命令列表\n!!server info|-if 查看当前服务端信息\n!!server infos|-ifs 查看所有服务端信息\n!!server plugins|-ps 查看插件列表\n!!server stop|-sp 停止当前服务端\n!!server restart|-rst 重启当前服务端\n!!server ban|-bn <命令> 禁止使用命令\n!!server unban|-ubn <命令> 解除禁止使用命令"
)

// self
const (
	maxLen = 5
)

var (
	stateMap map[int]string
)

type BasicPlugin struct {
	id       string
	mcServer server.MinecraftServer
}

func (p *BasicPlugin) GetHelpDescription() string {
	return helpDescription
}

func (p *BasicPlugin) ChangeConfCallBack() {}

func (p *BasicPlugin) DestructCallBack() {}

func (p *BasicPlugin) InitCallBack() {
	stateMap = make(map[int]string)
	// 0.未启动 1.启动  -1.正在启动 -2.正在关闭
	stateMap[constant.MC_STATE_STOP] = "未启动"
	stateMap[constant.MC_STATE_START] = "启动"
	stateMap[constant.MC_STATE_STARTIND] = "正在启动"
	stateMap[constant.MC_STATE_STOPING] = "正在关闭"
}

func (p *BasicPlugin) GetId() string {
	return p.id
}

func (p *BasicPlugin) GetName() string {
	return pluginName
}

func (p *BasicPlugin) GetDescription() string {
	return pluginDescription
}

func (p *BasicPlugin) GetCommandName() string {
	return pluginCommand
}

func (p *BasicPlugin) Init(server server.MinecraftServer) {
	p.mcServer = server
}

func (p *BasicPlugin) IsGlobal() bool {
	return isGlobal
}

func (p *BasicPlugin) NewInstance() plugin.Plugin {
	pl := &BasicPlugin{
		id: uuid.NewV4().String(),
	}
	modules.RegisterCallBack(pl)
	return pl
}

func (p *BasicPlugin) HandleMessage(message *json_struct.ReciveMessage) {
	if message.Player == "" {
		return
	}
	commandObj := utils.ParsePluginCommand(message.Speak)
	if commandObj.Command != pluginCommand {
		return
	}
	if len(commandObj.Params) == 0 {
		_ = p.mcServer.TellrawCommand(message.Player, helpDescription)
	} else {
		p.paramsHandle(message.Player, commandObj)
	}
}

func (p *BasicPlugin) Start() {}

func (p *BasicPlugin) Stop() {}

func (p *BasicPlugin) paramsHandle(player string, pc *json_struct.PluginCommand) {
	switch pc.Params[0] {
	case "info", "-if":
		header := []string{"id", "名称", "端口", "内存(单位：M)", "版本", "模式"}
		cfg := p.mcServer.GetServerConf()
		data := [][]string{{utils.Ellipsis(cfg.EntryId, maxLen), cfg.Name, strconv.Itoa(cfg.Port), strconv.Itoa(cfg.Memory), cfg.Version, cfg.GameType}}
		p.mcServer.TellrawCommand(player, utils.FormateTable(header, data))
	case "infos", "-ifs":
		header := []string{"id", "名称", "端口", "内存(单位：M)", "版本", "模式", "运行状态"}
		ctr := modules.GetMinecraftServerContainerInstance()
		aCfg := ctr.GetAllServerConf()
		data := make([][]string, 0)
		for _, cfg := range aCfg {
			// 镜像不展示
			if !cfg.IsMirror {
				data = append(data, []string{utils.Ellipsis(cfg.EntryId, maxLen), cfg.Name, strconv.Itoa(cfg.Port), strconv.Itoa(cfg.Memory), cfg.Version, cfg.GameType, stateMap[cfg.State]})
			}
		}
		p.mcServer.TellrawCommand(player, utils.FormateTable(header, data))
	case "plugins", "-ps":
		aPlcfg := p.mcServer.GetPluginsInfo()
		header := []string{"名称", "是否启用", "命令", "简介"}
		data := make([][]string, 0)
		for _, plcfg := range aPlcfg {
			isBanStr := "是"
			if !plcfg.IsBan {
				isBanStr = "否"
			}
			data = append(data, []string{plcfg.Name, isBanStr, plcfg.CommandName, plcfg.Description})
		}
		p.mcServer.TellrawCommand(player, utils.FormateTable(header, data))
	case "stop", "-sp":
		if err := modules.GetMinecraftServerContainerInstance().StopById(p.mcServer.GetServerEntryId()); err != nil {
			p.mcServer.TellrawCommand(player, "关闭失败")
		}
	case "restart", "-rst":
		if err := modules.GetMinecraftServerContainerInstance().RestartById(p.mcServer.GetServerEntryId()); err != nil {
			p.mcServer.TellrawCommand(player, "重启失败")
		}
	case "help", "-l":
		aPlcfg := p.mcServer.GetPluginsInfo()
		var (
			header []string
			data   [][]string
		)
		if len(pc.Params) < 2 {
			header = []string{"命令", "简介"}
			data = make([][]string, 0)
			for _, plcfg := range aPlcfg {
				data = append(data, []string{plcfg.CommandName, plcfg.Description})
			}

		} else {
			cmd := pc.Params[1]
			header = []string{"命令", "简介", "用法"}
			data = make([][]string, 0)
			for _, plcfg := range aPlcfg {
				if plcfg.CommandName == cmd {
					data = append(data, []string{plcfg.CommandName, plcfg.Description, plcfg.HelpDescription})
				}
			}
		}
		p.mcServer.TellrawCommand(player, utils.FormateTable(header, data))
	case "ban", "-bn":
		if len(pc.Params) < 2 {
			p.mcServer.TellrawCommand(player, "请输入命令")
			return
		}
		cmd := pc.Params[1]
		cmdObj := p.getPluginBYCmd(cmd)
		if cmdObj != nil {
			p.mcServer.BanPlugin(cmdObj.Id)
		}
	case "unban", "-ubn":
		if len(pc.Params) < 2 {
			p.mcServer.TellrawCommand(player, "请输入命令")
			return
		}
		cmd := pc.Params[1]
		cmdObj := p.getPluginBYCmd(cmd)
		if cmdObj != nil {
			p.mcServer.UnbanPlugin(cmdObj.Id)
		}
	}
}

func (p *BasicPlugin) getPluginBYCmd(cmd string) *json_struct.PluginInfo {
	aPlcfg := p.mcServer.GetPluginsInfo()
	for _, plcfg := range aPlcfg {
		if cmd == plcfg.CommandName {
			return plcfg
		}
	}
	return nil
}

var BasicPluginObj = &BasicPlugin{}
