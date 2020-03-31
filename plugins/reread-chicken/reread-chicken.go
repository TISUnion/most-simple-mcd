package reread_chicken

import (
	plugin_interface "github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/modules"
	uuid "github.com/satori/go.uuid"
	"regexp"
)

type RereadChickenPlugin struct {
	mcServer server.MinecraftServer
}

const reg = `!!repeat\s+(.+)`

func (r *RereadChickenPlugin) GetDescription() string {
	return "复读鸡"
}

func (r *RereadChickenPlugin) GetCommandName() string {
	return "!!repeat"
}

func (r *RereadChickenPlugin) ChangeConfCallBack() {}
func (r *RereadChickenPlugin) DestructCallBack()   {}
func (r *RereadChickenPlugin) InitCallBack()       {}
func (r *RereadChickenPlugin) GetId() string       { return uuid.NewV4().String() }
func (r *RereadChickenPlugin) GetName() string     { return "repeatMachine" }
func (r *RereadChickenPlugin) IsGlobal() bool      { return false }
func (r *RereadChickenPlugin) Start()              {}
func (r *RereadChickenPlugin) Stop()               {}
func (r *RereadChickenPlugin) HandleMessage(message *json_struct.ReciveMessage) {
	if message.Speak == "" {
		return
	}
	regObj := regexp.MustCompile(reg)
	match := regObj.FindStringSubmatch(message.Speak)
	if len(match) == 2 {
		_ = r.mcServer.Command("/say " + match[1])
	}
}
func (r *RereadChickenPlugin) Init(mcServer server.MinecraftServer) {
	r.mcServer = mcServer
}
func (r *RereadChickenPlugin) NewInstance() plugin_interface.Plugin {
	p := &RereadChickenPlugin{}
	modules.RegisterCallBack(p)
	return p
}

var RereadChickenPluginObj = &RereadChickenPlugin{}
