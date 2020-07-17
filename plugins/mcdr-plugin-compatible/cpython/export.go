package cpython

/*
#include "Utils.h"
*/
import "C"
import (
	"github.com/TISUnion/most-simple-mcd/modules"
	"unsafe"
)

type Server C.Server

var ctr = modules.GetMinecraftServerContainerInstance()

//export serverInfo
func serverInfo(cid C.CString) Server {
	srvId := C.GoString(cid)
	s, err := ctr.GetServerById(srvId)
	if err != nil {
		return nil
	}
	info := s.GetServerConf()
	id := C.CString(info.EntryId)
	name := C.CString(info.Name)
	memory := int(info.Memory)
	port := int(info.Port)
	version := C.CString(info.Version)
	side := C.CString(info.Side)
	comment := C.CString(info.Comment)
	defer C.free(unsafe.Pointer(id))
	defer C.free(unsafe.Pointer(version))
	defer C.free(unsafe.Pointer(side))
	defer C.free(unsafe.Pointer(comment))
	defer C.free(unsafe.Pointer(name))
	csrv := Server{name: name, id: id, memory: memory, port: port, version: version, side: side, comment: comment}
	return csrv
}

//export start
func start(cid C.CString) {
	id := C.GoString(cid)
}

//export stop
func stop(cid C.CString) {
}

//export restart
func restart(cid C.CString) {
}

//export stopExit
func stopExit(cid C.CString) {
}

//export exit
func exit(cid C.CString) {
}

//export isServerRunning
func isServerRunning(cid C.CString) int {
	return 0
}

//export isServerStartup
func isServerStartup(cid C.CString) int {
	return 0
}

//export isRconRunning
func isRconRunning(cid C.CString) int {
	return 0
}

//export execute
func execute(cid, ctext C.CString) {
}

//export say
func say(cid, ctext C.CString) {
}

//export tell
func tell(cid, cplayer, ctext C.CString) {
}

//export reply
func reply(cid, cplayer, ctext C.CString) {
}

//export loadPlugin
func loadPlugin(cpluginName C.CString) {
}

//export enablePlugin
func enablePlugin(cid C.CString) {
}

//export disablePlugin
func disablePlugin(cid C.CString) {
}

//export refreshAllPlugins
func refreshAllPlugins(cid C.CString) {
}

//export refreshChangedPlugins
func refreshChangedPlugins(cid C.CString) {
}

//export getPluginList
func getPluginList(cid C.CString) {
}
