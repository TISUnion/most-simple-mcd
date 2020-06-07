package constant

const (
	VANILLA         = "vanilla"
	VANILLA_VERSION = "\\[Server thread/INFO\\]: Starting minecraft server version ([0-9]*\\.?[0-9]*\\.?[0-9]*\\.?)"
	VANILLA_GAME_TYPE = "\\[Server thread/INFO\\]: Default game type: (?P<type>[a-zA-Z]+)"
	VANILLA_GAME_START = "\\[Server thread/INFO\\]: Done \\(.*\\)! For help, type \"help\""
	VANILLA_GAME_SAVE = "\\[Server thread/INFO\\]: Saved the world"
	VANILLA_MESSAGE = `\[(\d+:\d+:\d+)]\s+\[Server thread/INFO\]:\s+[<|\[]{1}(.+)[>|\]]{1}\s+(.+)`
)

const (
	SPOGOT = "spigot"
)
