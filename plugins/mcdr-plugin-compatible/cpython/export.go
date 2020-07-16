package cpython

import "C"

//export serverInfo
func serverInfo(id C.CString) {
}

//export start
func start(id C.CString) {
}

//export stop
func stop(id C.CString) {
}

//export restart
func restart(id C.CString) {
}

//export stopExit
func stopExit(id C.CString) {
}

//export exit
func exit(id C.CString) {
}

//export isServerRunning
func isServerRunning(id C.CString) int {
	return 0
}

//export isServerStartup
func isServerStartup(id C.CString) int {
	return 0
}

//export isRconRunning
func isRconRunning(id C.CString) int {
	return 0
}

//export execute
func execute(id, text C.CString)  {
}

//export say
func say(id, text C.CString)  {
}

//export tell
func tell(id, player, text C.CString)  {
}

//export reply
func reply(id, player, text C.CString)  {
}

//export loadPlugin
func loadPlugin(plugin_name C.CString)  {
}

//export enablePlugin
func enablePlugin(id C.CString)  {
}

//export disablePlugin
func disablePlugin(id C.CString)  {
}

//export refreshAllPlugins
func refreshAllPlugins(id C.CString)  {
}

//export refreshChangedPlugins
func refreshChangedPlugins(id C.CString)  {
}

//export getPluginList
func getPluginList(id C.CString)  {
}
