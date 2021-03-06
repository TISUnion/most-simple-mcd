package constant

/**
	参考 Fallen-Breath/MCDReforged
	https://github.com/Fallen-Breath/MCDReforged/tree/master/utils/parser
	TODO
*/
const (
	PLAYER_MESSAGE = `^<(\w{1,16})> (.*)$`
	//=======================原版端
	VANILLA_SERVER     = "vanilla"
	VANILLA_VERSION    = "Starting minecraft server version ([0-9]*\\.?[0-9]*\\.?[0-9]*\\.?)"
	VANILLA_GAME_TYPE  = "Default game type: (?P<type>[a-zA-Z]+)"
	VANILLA_GAME_START = "Done \\([0-9.]*s\\)! For help, type \"help\"( or \"\\?\")?"
	VANILLA_GAME_SAVE  = "\\[Server thread/INFO\\]: Saved the (game|world)"
	// [09:00:00] [Server thread/INFO]: <Steve> Hello
	VANILLA_MESSAGE = `\[(\d+:\d+:\d+)]\s+\[Server thread/(\w+)\]: (.*)`
	// Steve[/127.0.0.1:9864] logged in with entity id 131 at (187.2703, 146.79014, 404.84718)
	VANILLA_PLAYER_JOIN = `(\w{1,16})\[/[\d.:]+\] logged in with entity id \d+ at \([\dE\-., ]+\)`
	// Steve left the game
	VANILLA_PLAYER_LEFT = `(\w{1,16}) left the game`
	// Steve has made the advancement [Stone Age]
	VANILLA_PLAYER_ADVANCEMENT = `(\w{1,16}) has made the advancement \[(.+)\]`
	//======================Bukkit端
	BUKKIT_SERVER = "bukkit"
	// [09:00:01 INFO]: <Steve> hi
	BUKKIT_MESSAGE = `\[(\d+:\d+:\d+) (\w+)]: (.*)`
	// Fallen_Breath[/127.0.0.1:50099] logged in with entity id 11 at ([lobby]0.7133817548136454, 4.0, 5.481879061970788)
	BUKKIT_JOIN = `（\w{1,16}）\[/[\d.:]+\] logged in with entity id \d+ at \((\[\w+\])?[\dE\-., ]+\)`
	//======================Bukkit端，1.14
	BUKKIT14 = "bukkit14"
	//======================Cat端
	CAT_SERVER = "cat"
	//======================forge端
	FORGE_SERVER = "forge"
	//======================bungeecord转发
	BUNGEECORD_PROXY = "bungeecord"
	//======================waterfall转发
	WATERFALL_PROXY = "waterfall"

	SPOGOT = "spigot"
)

// 玩家死亡信息
// 参考 Minecraft wiki
// https://minecraft.gamepedia.com/Death_messages
var DeathMessage = []string{
	"(\\w{1,16}) blew up",
	"(\\w{1,16}) burned to death",
	"(\\w{1,16}) didn't want to live in the same world as .+",
	"(\\w{1,16}) died",
	"(\\w{1,16}) died because of .+",
	"(\\w{1,16}) discovered floor was lava",
	"(\\w{1,16}) discovered the floor was lava",
	"(\\w{1,16}) drowned",
	"(\\w{1,16}) drowned whilst trying to escape .+",
	"(\\w{1,16}) experienced kinetic energy",
	"(\\w{1,16}) experienced kinetic energy whilst trying to escape .+",
	"(\\w{1,16}) fell from a high place",
	"(\\w{1,16}) fell off a ladder",
	"(\\w{1,16}) fell off a scaffolding",
	"(\\w{1,16}) fell off some twisting vines",
	"(\\w{1,16}) fell off some vines",
	"(\\w{1,16}) fell off some weeping vines",
	"(\\w{1,16}) fell out of the water",
	"(\\w{1,16}) fell out of the world",
	"(\\w{1,16}) fell too far and was finished by .+",
	"(\\w{1,16}) fell too far and was finished by .+ using .+",
	"(\\w{1,16}) fell while climbing",
	"(\\w{1,16}) hit the ground too hard",
	"(\\w{1,16}) hit the ground too hard whilst trying to escape .+",
	"(\\w{1,16}) starved to death",
	"(\\w{1,16}) starved to death whilst fighting .+",
	"(\\w{1,16}) suffocated in a wall",
	"(\\w{1,16}) suffocated in a wall whilst fighting .+",
	"(\\w{1,16}) tried to swim in lava",
	"(\\w{1,16}) tried to swim in lava to escape .+",
	"(\\w{1,16}) walked into a cactus whilst trying to escape .+",
	"(\\w{1,16}) walked into danger zone due to .+",
	"(\\w{1,16}) walked into fire whilst fighting .+",
	"(\\w{1,16}) walked on danger zone due to .+",
	"(\\w{1,16}) was blown up by .+",
	"(\\w{1,16}) was blown up by .+ using .+",
	"(\\w{1,16}) was burnt to a crisp whilst fighting .+",
	"(\\w{1,16}) was doomed to fall",
	"(\\w{1,16}) was doomed to fall by .+",
	"(\\w{1,16}) was doomed to fall by .+ using .+",
	"(\\w{1,16}) was fireballed by .+",
	"(\\w{1,16}) was fireballed by .+ using .+",
	"(\\w{1,16}) was impaled by Trident",
	"(\\w{1,16}) was impaled by .+",
	"(\\w{1,16}) was impaled by .+ with .+",
	"(\\w{1,16}) was killed by [Intentional Game Design]",
	"(\\w{1,16}) was killed by .+ trying to hurt .+",
	"(\\w{1,16}) was killed by .+ using .+",
	"(\\w{1,16}) was killed by .+ using magic",
	"(\\w{1,16}) was killed by even more magic",
	"(\\w{1,16}) was killed by magic",
	"(\\w{1,16}) was killed by magic whilst trying to escape .+",
	"(\\w{1,16}) was killed trying to hurt .+",
	"(\\w{1,16}) was poked to death by a sweet berry bush",
	"(\\w{1,16}) was poked to death by a sweet berry bush whilst trying to escape .+",
	"(\\w{1,16}) was pricked to death",
	"(\\w{1,16}) was pummeled by .+",
	"(\\w{1,16}) was pummeled by .+ using .+",
	"(\\w{1,16}) was roasted in dragon breath",
	"(\\w{1,16}) was roasted in dragon breath by .+",
	"(\\w{1,16}) was shot by Arrow",
	"(\\w{1,16}) was shot by .+",
	"(\\w{1,16}) was shot by .+ using .+",
	"(\\w{1,16}) was slain by Arrow",
	"(\\w{1,16}) was slain by Small Fireball",
	"(\\w{1,16}) was slain by Trident",
	"(\\w{1,16}) was slain by .+ and \\w{1,16} was slain by \\w{1,16}.",
	"(\\w{1,16}) was slain by .+ using .+ and \\w{1,16} was slain by \\w{1,16} using .+.",
	"(\\w{1,16}) was slain by .+",
	"(\\w{1,16}) was slain by .+ using .+",
	"(\\w{1,16}) was slain by \\w{1,16} using .+",
	"(\\w{1,16}) was squashed by .+",
	"(\\w{1,16}) was squashed by a falling anvil",
	"(\\w{1,16}) was squashed by a falling anvil whilst fighting .+",
	"(\\w{1,16}) was squashed by a falling block",
	"(\\w{1,16}) was squashed by a falling block whilst fighting .+",
	"(\\w{1,16}) was squished too much",
	"(\\w{1,16}) was struck by lightning",
	"(\\w{1,16}) was struck by lightning whilst fighting .+",
	"(\\w{1,16}) was stung to death",
	"(\\w{1,16}) was stung to death by .+",
	"(\\w{1,16}) went off with a bang",
	"(\\w{1,16}) went off with a bang whilst fighting .+",
	"(\\w{1,16}) went up in flames",
	"(\\w{1,16}) withered away",
	"(\\w{1,16}) withered away whilst fighting .+",
}
