package mcdr_plugin_compatible

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	uuid "github.com/satori/go.uuid"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	pluginName                         = "兼容mcdr插件"
	pluginDescription                  = "兼容mcdr插件"
	pluginCommand                      = "!!mcdr"
	isMcdrPluginCompatiblePluginGlobal = true
	helpDescription                    = "插件使用方法"
	pluginDir                          = "mcdr-plugins"
	pythonExt                          = ".py"
)

type McdrPluginCompatiblePlugin struct {
	pluginRootPath string
	plugins        map[string]*McdrPlugin
	id             string
}

func (p *McdrPluginCompatiblePlugin) GetDescription() string {
	return pluginDescription
}

func (p *McdrPluginCompatiblePlugin) GetHelpDescription() string {
	return helpDescription
}

func (p *McdrPluginCompatiblePlugin) GetCommandName() string {
	return pluginCommand
}

func (p *McdrPluginCompatiblePlugin) IsGlobal() bool {
	return isMcdrPluginCompatiblePluginGlobal
}

func (p *McdrPluginCompatiblePlugin) GetId() string {
	return p.id
}

func (p *McdrPluginCompatiblePlugin) GetName() string {
	return pluginName
}

func (p *McdrPluginCompatiblePlugin) Init(mcServer server.MinecraftServer) {

}

/* ------------------回调接口-------------------- */
func (p *McdrPluginCompatiblePlugin) ChangeConfCallBack() {}
func (p *McdrPluginCompatiblePlugin) DestructCallBack()   {
	PyVmEnd()
}
func (p *McdrPluginCompatiblePlugin) InitCallBack() {
	if !PyVmStart() {
		modules.WriteLogToDefault("python虚拟机开启失败")
		modules.SendExitSign()
	}
	p.plugins = make(map[string]*McdrPlugin)
	p.scanPlugin()
}

/* --------------------------------------------- */

/* ---------非全局插件，服务端启动，关闭回调--------- */
func (p *McdrPluginCompatiblePlugin) Start() {}
func (p *McdrPluginCompatiblePlugin) Stop()  {}

/* --------------------------------------------- */

func (p *McdrPluginCompatiblePlugin) HandleMessage(message *models.ReciveMessage) {
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

func (p *McdrPluginCompatiblePlugin) paramsHandle(player string, pc *models.PluginCommand, mcServer server.MinecraftServer) {
	switch pc.Params[0] {
	// write code...
	default:
		_ = mcServer.TellrawCommand(player, helpDescription)
	}
}

func (*McdrPluginCompatiblePlugin) NewInstance() plugin.Plugin {
	return nil
}

var McdrPluginCompatiblePluginObj plugin.Plugin

func GetMcdrPluginCompatiblePluginInstance() plugin.Plugin {
	if McdrPluginCompatiblePluginObj != nil {
		return McdrPluginCompatiblePluginObj
	}
	McdrPluginCompatiblePluginObj = &McdrPluginCompatiblePlugin{
		id:             uuid.NewV4().String(),
		pluginRootPath: filepath.Join(modules.GetConfVal(constant.WORKSPACE), pluginDir),
	}
	modules.RegisterCallBack(McdrPluginCompatiblePluginObj)
	return McdrPluginCompatiblePluginObj
}

func (p *McdrPluginCompatiblePlugin) scanPlugin() {
	err := filepath.Walk(p.pluginRootPath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		name := info.Name()
		if filepath.Ext(name) == pythonExt {
			// 兼容windows
			path = strings.ReplaceAll(path, "\\", "/")

			re := regexp.MustCompile(fmt.Sprintf("(.*%s/)|(.py)", pluginDir))
			pluginName := re.ReplaceAllString(path, "")
			packageName := strings.ReplaceAll(pluginName, "/", ".")
			pluginName = strings.ReplaceAll(pluginName, "/", "_")
			// 已经加载过的，则不重复加载
			if _, ok := p.plugins[pluginName]; ok {
				return nil
			}
			packageName = fmt.Sprint(pluginDir, ".", packageName)
			pluginInfo := fmt.Sprint("mcdr插件: ", pluginName)
			plp := &McdrPlugin{
				pluginPath:        path,
				packageName:       packageName,
				pluginName:        pluginName,
				id:                uuid.NewV4().String(),
				pluginDescription: pluginInfo,
				helpDescription:   pluginInfo,
			}
			modules.RegisterCallBack(plp)
			p.plugins[pluginName] = plp
		}
		return nil
	})
	if err != nil {
		modules.WriteLogToDefault("加载插件失败！")
	}
}

// TODO 分发事件
//func ()  {
//
//}