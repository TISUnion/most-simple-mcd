package plugins

import (
	plugin_interface "github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/plugins/broadcast"
	"github.com/TISUnion/most-simple-mcd/plugins/here"
	mcdr_plugin_compatible "github.com/TISUnion/most-simple-mcd/plugins/mcdr-plugin-compatible"
	mirror_server "github.com/TISUnion/most-simple-mcd/plugins/mirror-server"
	reread_chicken "github.com/TISUnion/most-simple-mcd/plugins/reread-chicken"
)

var plugins = []plugin_interface.Plugin{
	reread_chicken.RereadChickenPluginObj,
	mirror_server.GetMirrorServerPluginInstance(),
	BasicPluginObj,
	broadcast.GetBroadcastPluginInstance(),
	here.HerePluginObj,
	mcdr_plugin_compatible.GetMcdrPluginCompatiblePluginInstance(),
}
