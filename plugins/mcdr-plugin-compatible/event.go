package mcdr_plugin_compatible

import (
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"regexp"
)

// TODO
func (p *McdrPlugin) dealEvent(srv server.MinecraftServer, originData string) {
	re := regexp.MustCompile(srv.GetPlayerJoinRegularExpression())
	match := re.FindStringSubmatch(originData)
	if len(match) > 0 {

	}
}

func (p *McdrPlugin) onPlayerJoined(srv server.MinecraftServer) {

}

func (p *McdrPlugin) onPlayerLeft(srv server.MinecraftServer) {

}

func (p *McdrPlugin) onPlayerMadeAdvancement(srv server.MinecraftServer) {

}

func (p *McdrPlugin) onDeathMessage(srv server.MinecraftServer) {

}

func (p *McdrPlugin) onServerStartup(srv server.MinecraftServer) {

}

func (p *McdrPlugin) onServerStop(srv server.MinecraftServer) {

}

func (p *McdrPlugin) onMcdrStop() {

}
