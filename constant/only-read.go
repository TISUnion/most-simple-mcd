package constant

// 只读变量
var (
	CLI_THX      = []byte{109, 99, 100, 101, 103, 103}
	CLI_THX_TEXT = []byte{116, 104, 97, 110, 107, 32, 121, 111, 117, 32, 102, 111, 114, 32, 117, 115, 105, 110, 103, 32, 77, 79, 83, 84, 45, 83, 73, 77, 80, 76, 69, 45, 77, 67, 68}

	PLUGIN_TMPL = `package {{filename2packagename .Dirname}}

import (
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	uuid "github.com/satori/go.uuid"
)

const (
	pluginName        = "{{.ZHName}}"
	pluginDescription = "{{.Description}}"
	pluginCommand     = "{{.Command}}"
	isGlobal          = {{.IsGlobal}}
	helpDescription   = "{{.HelpDescription}}"
)

type {{.ENName}}Plugin struct {
	{{if .IsGlobal}}{{else}}mcServer server.MinecraftServer{{end}}
    id       string
}

func (p *{{.ENName}}Plugin) GetDescription() string {
	return pluginDescription
}

func (p *{{.ENName}}Plugin) GetHelpDescription() string {
	return helpDescription
}

func (p *{{.ENName}}Plugin) GetCommandName() string {
	return pluginCommand
}

func (p *{{.ENName}}Plugin) IsGlobal() bool {
    return isGlobal
}

func (p *{{.ENName}}Plugin) GetId() string {
    return p.id
}

func (p *{{.ENName}}Plugin) GetName() string     {
    return pluginName
}

func (p *{{.ENName}}Plugin) Init(mcServer server.MinecraftServer) {
    {{if .IsGlobal}}{{else}}p.mcServer = mcServer{{end}}
}

/* ------------------回调接口-------------------- */
func (p *{{.ENName}}Plugin) ChangeConfCallBack() {}
func (p *{{.ENName}}Plugin) DestructCallBack()   {}
func (p *{{.ENName}}Plugin) InitCallBack()       {}
/* --------------------------------------------- */

/* ---------非全局插件，服务端启动，关闭回调--------- */
func (p *{{.ENName}}Plugin) Start()              {}
func (p *{{.ENName}}Plugin) Stop()               {}
/* --------------------------------------------- */

func (p *{{.ENName}}Plugin) HandleMessage(message *models.ReciveMessage) {
	if !message.IsPlayer {
		return
	}
	if message.Command != pluginCommand {
		return
	}
	{{if .IsGlobal}}
	mcServer, err := modules.GetMinecraftServerContainerInstance().GetServerById(message.ServerId)
    if err != nil {
    	return
    }
    {{end}}
	if len(message.Params) == 0 {
		_ = {{if .IsGlobal}}mcServer{{else}}p.mcServer{{end}}.TellrawCommand(message.Player, helpDescription)
	} else {
        p.paramsHandle(message.Player, message{{if .IsGlobal}}, mcServer{{end}})
	}
}

func (p *{{.ENName}}Plugin) paramsHandle(player string, pc *models.ReciveMessage{{if .IsGlobal}}, mcServer server.MinecraftServer{{end}}) {
    switch pc.Params[0] {
    // write code...
    default:
    	_ = {{if .IsGlobal}}mcServer{{else}}p.mcServer{{end}}.TellrawCommand(player, helpDescription)
    }
}

func (*{{.ENName}}Plugin) NewInstance() plugin.Plugin {
    {{if .IsGlobal}}return nil{{else}}plg := &{{.ENName}}Plugin{
		id: uuid.NewV4().String(),
	}
	modules.RegisterCallBack(plg)
	return plg{{end}}
}
{{if .IsGlobal}}
var {{.ENName}}PluginObj plugin.Plugin

func Get{{.ENName}}PluginInstance() plugin.Plugin {
	if {{.ENName}}PluginObj != nil {
		return {{.ENName}}PluginObj
	}
	{{.ENName}}PluginObj = &{{.ENName}}Plugin{
		id: uuid.NewV4().String(),
	}
	modules.RegisterCallBack({{.ENName}}PluginObj)
	return {{.ENName}}PluginObj
}

{{else}}var {{.ENName}}PluginObj = &{{.ENName}}Plugin{}{{end}}`
)
