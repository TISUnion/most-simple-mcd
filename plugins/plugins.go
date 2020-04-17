package plugins

import (
	plugin_interface "github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/plugins/broadcast"
	mirror_server "github.com/TISUnion/most-simple-mcd/plugins/mirror-server"
	reread_chicken "github.com/TISUnion/most-simple-mcd/plugins/reread-chicken"
)

var plugins = []plugin_interface.Plugin{
	reread_chicken.RereadChickenPluginObj,
	mirror_server.GetMirrorServerPluginInstance(),
	BasicPluginObj,
	broadcast.BroadcastPluginObj,
}
