package reread_chicken

import (
	plugin_interface "github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/modules"
)

type RereadChickenPlugin struct{}

func (r *RereadChickenPlugin) GetDescription() string {
	return ""
}

func (r *RereadChickenPlugin) GetCommandName() string {
	return ""
}

func (r *RereadChickenPlugin) ChangeConfCallBack()                                  {}
func (r *RereadChickenPlugin) DestructCallBack()                                    {}
func (r *RereadChickenPlugin) InitCallBack()                                        {}
func (r *RereadChickenPlugin) GetId() string                                        { return "" }
func (r *RereadChickenPlugin) GetName() string                                      { return "" }
func (r *RereadChickenPlugin) IsGlobal() bool                                       { return false }
func (r *RereadChickenPlugin) Start()                                               {}
func (r *RereadChickenPlugin) Stop()                                                {}
func (r *RereadChickenPlugin) HandleMessage(messageType *json_struct.ReciveMessage) {}
func (r *RereadChickenPlugin) Init(server server.MinecraftServer)                   {}
func (r *RereadChickenPlugin) NewInstance() plugin_interface.Plugin {
	p := &RereadChickenPlugin{}
	modules.RegisterCallBack(p)
	return p
}

var RereadChickenPluginObj = &RereadChickenPlugin{}
