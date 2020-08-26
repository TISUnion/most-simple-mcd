package mcdr_plugin_compatible

import (
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"regexp"
)

// TODO
func (p *McdrPlugin) dealEvent(srv server.MinecraftServer, originData string) {
	var (
		re          *regexp.Regexp
		match       []string
		player      string
		advancement string
	)
	// 加入游戏
	srv.GetPlayerJoinRegularExpression(originData)
	if len(match) > 0 {
		player = match[0]
		p.onPlayerJoined(srv, player)
		return
	}

	// 离开游戏
	re = regexp.MustCompile(srv.GetPlayerLeftRegularExpression())
	match = re.FindStringSubmatch(originData)
	if len(match) > 0 {
		player = match[0]
		p.onPlayerLeft(srv, player)
		return
	}

	// 获得成就
	re = regexp.MustCompile(srv.GetPlayerAdvancementRegularExpression())
	match = re.FindStringSubmatch(originData)
	if len(match) > 1 {
		player = match[0]
		advancement = match[1]
		p.onPlayerMadeAdvancement(srv, player, advancement)
		return
	}
}

func (p *McdrPlugin) onPlayerJoined(srv server.MinecraftServer, player string) {

}

func (p *McdrPlugin) onPlayerLeft(srv server.MinecraftServer, player string) {

}

func (p *McdrPlugin) onPlayerMadeAdvancement(srv server.MinecraftServer, player , advancement string) {

}

func (p *McdrPlugin) onDeathMessage(srv server.MinecraftServer, player, deathMsg string) {

}

func (p *McdrPlugin) onServerStartup(srv server.MinecraftServer) {

}

func (p *McdrPlugin) onServerStop(srv server.MinecraftServer) {

}

func (p *McdrPlugin) onMcdrStop() {

}
