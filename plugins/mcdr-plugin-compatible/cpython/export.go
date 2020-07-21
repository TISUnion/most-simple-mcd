package cpython

/*
#cgo pkg-config:python3
#include "Utils.h"
#define NAME "name"
#define ID "id"
#define PORT "port"
#define MEMORY "memory"
#define VERSION "version"
#define SIDE "side"
#define COMMENT "comment"
#define MCDR "MCDR"

Server mcServerInfo(char *id);
void mcStart(char *id);
void mcStop(char *id);
void mcRestart(char *id);
void mcStopExit(char *id);
void mcExit(char *id);
int mcIsServerRunning(char *id);
int mcIsServerStartup(char *id);
int mcIsRconRunning(char *id);
void mcExecute(char *id, char *text);
void mcSay(char *id, char *text);
void mcTell(char *id, char *player, char *text);
void mcReply(char *id, char *player, char *text);

static PyObject *py_mc_start(PyObject *self, PyObject *args)
{
	PyObject *pid = PyObject_GetAttrString(self, "id");
	char *id = pyObj2string(pid);
	mcStart(id);
	return Py_None;
}

static PyObject *py_mc_stop(PyObject *self, PyObject *args)
{
	PyObject *pid = PyObject_GetAttrString(self, "id");
	char *id = pyObj2string(pid);
	mcStop(id);
	return Py_None;
}

static PyObject *py_mc_restart(PyObject *self, PyObject *args)
{
	PyObject *pid = PyObject_GetAttrString(self, "id");
	char *id = pyObj2string(pid);
	mcRestart(id);
	return Py_None;
}

static PyObject *py_mc_stopExit(PyObject *self, PyObject *args)
{
	PyObject *pid = PyObject_GetAttrString(self, "id");
	char *id = pyObj2string(pid);
	mcStopExit(id);
	return Py_None;
}

static PyObject *py_mc_exit(PyObject *self, PyObject *args)
{
	PyObject *pid = PyObject_GetAttrString(self, "id");
	char *id = pyObj2string(pid);
	mcExit(id);
	return Py_None;
}

static PyObject *py_mc_isServerRunning(PyObject *self, PyObject *args)
{
	PyObject *pid = PyObject_GetAttrString(self, "id");
	char *id = pyObj2string(pid);
	int res = mcIsServerRunning(id);
	return (PyObject *)Py_BuildValue("i", res);
}

static PyObject *py_mc_isServerStartup(PyObject *self, PyObject *args)
{
	PyObject *pid = PyObject_GetAttrString(self, "id");
	char *id = pyObj2string(pid);
	int res = mcIsServerRunning(id);
	return (PyObject *)Py_BuildValue("i", res);
}

static PyObject *py_mc_isRconRunning(PyObject *self, PyObject *args)
{
	PyObject *pid = PyObject_GetAttrString(self, "id");
	char *id = pyObj2string(pid);
	int res = mcIsRconRunning(id);
	return (PyObject *)Py_BuildValue("i", res);
}

static PyObject *py_mc_execute(PyObject *self, PyObject *args)
{
	PyObject *pid = PyObject_GetAttrString(self, "id");
	char *id = pyObj2string(pid);
	PyObject *ptext = PyObject_GetAttrString(self, "text");
	char *text = pyObj2string(ptext);
	mcExecute(id, text);
	return Py_None;
}

static PyObject *py_mc_say(PyObject *self, PyObject *args)
{
	PyObject *pid = PyObject_GetAttrString(self, "id");
	char *id = pyObj2string(pid);
	PyObject *ptext = PyObject_GetAttrString(self, "text");
	char *text = pyObj2string(ptext);
	mcSay(id, text);
	return Py_None;
}

static PyObject *py_mc_tell(PyObject *self, PyObject *args)
{
	PyObject *pid = PyObject_GetAttrString(self, "id");
	char *id = pyObj2string(pid);
	PyObject *ptext = PyObject_GetAttrString(self, "text");
	char *text = pyObj2string(ptext);
	PyObject *pplayer = PyObject_GetAttrString(self, "player");
	char *player = pyObj2string(pplayer);
	mcTell(id, player, text);
	return Py_None;
}

static PyObject *py_mc_reply(PyObject *self, PyObject *args)
{
	PyObject *pid = PyObject_GetAttrString(self, "id");
	char *id = pyObj2string(pid);
	PyObject *ptext = PyObject_GetAttrString(self, "text");
	char *text = pyObj2string(ptext);
	PyObject *pplayer = PyObject_GetAttrString(self, "player");
	char *player = pyObj2string(pplayer);
	mcReply(id, player, text);
	return Py_None;
}

static PyObject *GetServer(char *id)
{
	PyObject *server = CreateClass(id, NULL);
	//设置属性
	Server s = mcServerInfo(id);
	PyObject_SetAttrString(server, NAME, PyUnicode_FromString(s.name));
	PyObject_SetAttrString(server, ID, PyUnicode_FromString(s.id));
	PyObject_SetAttrString(server, PORT, (PyObject *)Py_BuildValue("i",s.port));
	PyObject_SetAttrString(server, MEMORY, (PyObject *)Py_BuildValue("i",s.memory));
	PyObject_SetAttrString(server, VERSION, PyUnicode_FromString(s.version));
	PyObject_SetAttrString(server, SIDE, PyUnicode_FromString(s.side));
	PyObject_SetAttrString(server, COMMENT, PyUnicode_FromString(s.comment));
	PyObject_SetAttrString(server, MCDR, (PyObject *)Py_BuildValue("i", 1));
	// 设置方法
	PyMethodDef start_def = {id, (PyCFunction)py_mc_start, METH_NOARGS, ""};
	PyObject *Py_start = PyCFunction_New(&start_def, server);
	PyObject_SetAttrString(server, "start", Py_start);

	PyMethodDef stop_def = {id, (PyCFunction)py_mc_stop, METH_NOARGS, ""};
	PyObject *Py_stop = PyCFunction_New(&stop_def, server);
	PyObject_SetAttrString(server, "stop", Py_stop);

	PyMethodDef restart_def = {id, (PyCFunction)py_mc_restart, METH_NOARGS, ""};
	PyObject *Py_restart = PyCFunction_New(&restart_def, server);
	PyObject_SetAttrString(server, "restart", Py_restart);

	PyMethodDef stopExit_def = {id, (PyCFunction)py_mc_stopExit, METH_NOARGS, ""};
	PyObject *Py_stopExit = PyCFunction_New(&stopExit_def, server);
	PyObject_SetAttrString(server, "stop_exit", Py_stopExit);

	PyMethodDef exit_def = {id, (PyCFunction)py_mc_exit, METH_NOARGS, ""};
	PyObject *Py_exit = PyCFunction_New(&exit_def, server);
	PyObject_SetAttrString(server, "exit", Py_exit);

	PyMethodDef isServerRunning_def = {id, (PyCFunction)py_mc_isServerRunning, METH_NOARGS, ""};
	PyObject *Py_isServerRunning = PyCFunction_New(&isServerRunning_def, server);
	PyObject_SetAttrString(server, "is_server_running", Py_isServerRunning);

	PyMethodDef isServerStartup_def = {id, (PyCFunction)py_mc_isServerStartup, METH_NOARGS, ""};
	PyObject *Py_isServerStartup = PyCFunction_New(&isServerStartup_def, server);
	PyObject_SetAttrString(server, "is_server_startup", Py_isServerStartup);

	return server;
}
*/
import "C"
import (
	"github.com/TISUnion/most-simple-mcd/modules"
	"unsafe"
)

type Server C.Server
type PyObject *C.PyObject

var ctr = modules.GetMinecraftServerContainerInstance()

//export mcServerInfo
func mcServerInfo(cid *C.char) Server {
	srvId := C.GoString(cid)
	defer C.free(unsafe.Pointer(cid))
	s, err := ctr.GetServerById(srvId)
	if err != nil {
		return Server{}
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
	csrv := Server{name: name, id: id, memory: C.int(memory), port: C.int(port), version: version, side: side, comment: comment}
	return csrv
}

//export mcStart
func mcStart(cid *C.char) {
	//id := C.GoString(cid)
}

//export mcStop
func mcStop(cid *C.char) {
}

//export mcRestart
func mcRestart(cid *C.char) {
}

//export mcStopExit
func mcStopExit(cid *C.char) {
}

//export mcExit
func mcExit(cid *C.char) {
}

//export mcIsServerRunning
func mcIsServerRunning(cid *C.char) C.int {
	return 0
}

//export mcIsServerStartup
func mcIsServerStartup(cid *C.char) C.int {
	return 0
}

//export mcIsRconRunning
func mcIsRconRunning(cid *C.char) C.int {
	return 0
}

//export mcExecute
func mcExecute(cid, ctext *C.char) {
}

//export mcSay
func mcSay(cid, ctext *C.char) {
}

//export mcTell
func mcTell(cid, cplayer, ctext *C.char) {
}

//export mcReply
func mcReply(cid, cplayer, ctext *C.char) {
}

//export mcLoadPlugin
func mcLoadPlugin(cpluginName *C.char) {
}

//export mcEnablePlugin
func mcEnablePlugin(cid *C.char) {
}

//export mcDisablePlugin
func mcDisablePlugin(cid *C.char) {
}

//export mcRefreshAllPlugins
func mcRefreshAllPlugins(cid *C.char) {
}

//export mcRefreshChangedPlugins
func mcRefreshChangedPlugins(cid *C.char) {
}

//export mcGetPluginList
func mcGetPluginList(cid *C.char) {
}
