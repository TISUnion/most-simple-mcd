package mcdr_plugin_compatible

/*
#cgo pkg-config:python3
#include "utils.h"
#define NAME "name"
#define ID "id"
#define PORT "port"
#define MEMORY "memory"
#define VERSION "version"
#define SIDE "side"
#define COMMENT "comment"
#define MCDR "MCDR"
#define MAX_PLUGIN_COUNT 300

Server mcServerInfo(char *id);
PyObject *pluginModules[300];
int pluginModulesCount = 1;
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

// 获取py的server对象
PyObject *GetServer(char *id)
{
	PyObject *server = CreateClass(id, NULL);
	//设置属性
	Server s = mcServerInfo(id);
	if (!strcmp(s.id, "")) {
		return NULL;
	}
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

// 设置插件包字典
int SetPlugin(char *packageName)
{
	pluginModules[pluginModulesCount] = PyModule_GetDict(PyImport_Import(PyUnicode_FromString(packageName)));
	return pluginModulesCount++;
}

// 重载插件包字典
void FreshPlugin(int pluginIndex, char *packageName)
{
	PyObject *pluginModule = pluginModules[pluginIndex];
	Py_XDECREF(pluginModule);
	pluginModules[pluginIndex] = PyImport_Import(PyUnicode_FromString(packageName));
}

// 调用插件回调
void CallBackPlugin(char *id, int pluginIndex, char *funcName)
{
	//PyObject *server = GetServer(id);
	//PyObject *pluginModule = pluginModules[pluginIndex];
	//PyObject *pyFuncName = PyUnicode_FromString(funcName);
	//if (!PyDict_Contains(pluginModule, pyFuncName)) { //不存在函数，就退出
	//	return;
	//}
	//PyObject *pluginFunc = PyDict_GetItemString(pluginModule, pyFuncName);
	//PyObject_CallFunction(pluginFunc, "OOO", server, )
}


*/
import "C"
import (
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/modules"
	"unsafe"
)

type Server C.Server
type Info C.Info

var ctr = modules.GetMinecraftServerContainerInstance()

//export mcServerInfo
func mcServerInfo(cid *C.char) Server {
	srvId := C.GoString(cid)
	defer C.free(unsafe.Pointer(cid))
	s, err := ctr.GetServerById(srvId)
	// 获取失败
	if err != nil {
		id := C.CString("")
		defer C.free(unsafe.Pointer(id))
		return Server{id: id}
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
	return Server{name: name, id: id, memory: C.int(memory), port: C.int(port), version: version, side: side, comment: comment}
}

//export mcStart
func mcStart(cid *C.char) {
	srvId := C.GoString(cid)
	srv, err := ctr.GetServerById(srvId)
	if err == nil {
		_ = srv.Start()
	}
}

//export mcStop
func mcStop(cid *C.char) {
	srvId := C.GoString(cid)
	srv, err := ctr.GetServerById(srvId)
	if err == nil {
		_ = srv.Stop()
	}
}

//export mcRestart
func mcRestart(cid *C.char) {
	srvId := C.GoString(cid)
	srv, err := ctr.GetServerById(srvId)
	if err == nil {
		_ = srv.Restart()
	}
}

//export mcStopExit
func mcStopExit(cid *C.char) {
	modules.SendExitSign()
}

//export mcExit
func mcExit(cid *C.char) {
	modules.SendExitSign()
}

//export mcIsServerRunning
func mcIsServerRunning(cid *C.char) C.int {
	srvId := C.GoString(cid)
	srv, err := ctr.GetServerById(srvId)
	if err == nil {
		if srv.GetServerConf().State == constant.MC_STATE_START || srv.GetServerConf().State == constant.MC_STATE_STARTIND {
			return constant.MC_STATE_START
		}
		return constant.MC_STATE_STOP
	}
	return constant.MC_STATE_STOP
}

//export mcIsServerStartup
func mcIsServerStartup(cid *C.char) C.int {
	srvId := C.GoString(cid)
	srv, err := ctr.GetServerById(srvId)
	if err == nil {
		if srv.GetServerConf().State == constant.MC_STATE_START {
			return constant.MC_STATE_START
		}
		return constant.MC_STATE_STOP
	}
	return constant.MC_STATE_STOP
}

//export mcIsRconRunning
func mcIsRconRunning(cid *C.char) C.int {
	return 0
}

//export mcExecute
func mcExecute(cid, ctext *C.char) {
	srvId := C.GoString(cid)
	text := C.GoString(ctext)
	srv, err := ctr.GetServerById(srvId)
	if err == nil {
		_ = srv.Command(text)
	}
}

//export mcSay
func mcSay(cid, ctext *C.char) {
	srvId := C.GoString(cid)
	text := C.GoString(ctext)
	srv, err := ctr.GetServerById(srvId)
	if err == nil {
		_ = srv.SayCommand(text)
	}
}

//export mcTell
func mcTell(cid, cplayer, ctext *C.char) {
	srvId := C.GoString(cid)
	text := C.GoString(ctext)
	player := C.GoString(cplayer)
	srv, err := ctr.GetServerById(srvId)
	if err == nil {
		_ = srv.TellCommand(player, text)
	}
}

//export mcReply
func mcReply(cid, cplayer, ctext *C.char) {
	srvId := C.GoString(cid)
	text := C.GoString(ctext)
	player := C.GoString(cplayer)
	srv, err := ctr.GetServerById(srvId)
	if err == nil {
		if player == "" {
			modules.WriteLogToDefault(text)
		} else {
			_ = srv.TellCommand(player, text)
		}
	}
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

func PyVmStart() bool {
	res := C.PyVmStart()
	if res == 0 {
		return false
	}
	return true
}

func PyVmEnd() {
	C.PyVmEnd()
}

func SetPyPlugin(packageName string) int {
	CPackageName := C.CString(packageName)
	defer C.free(unsafe.Pointer(CPackageName))
	return C.SetPlugin(CPackageName)
}

func FreshPyPlugin(packageName string, index int) {
	CPackageName := C.CString(packageName)
	defer C.free(unsafe.Pointer(CPackageName))
	C.FreshPlugin(CPackageName, index)
}

// TODO
func TriggerEvent(event int, command *models.ReciveMessage) {
	id := modules.GetIncreateId()
	switch event {
	case C.OnLoad:
	case C.OnUnload:
	case C.OnInfo:
	case C.OnUserInfo:
	case C.OnPlayerJoined:
	}
}

//func StartTest() {
//	C.PyVmStart()
//	p := C.GetServer(C.CString("test-123"))
//	if p == nil {
//		fmt.Println(123)
//	}
//	C.PyVmEnd()
//}

/**
* PyObject *server有以下属性：
*
* server.name                                      服务端名称
* server.id                                        服务端实例id
* server.memory                                    服务端使用内存
* server.version                                   服务端版本
* server.side                                      服务端类型
* server.comment                                   服务端备注
* server.MCDR                                      是否运行在MCDR，固定为true
* server.logger                                    命令行输出
* server.start()                                   开启服务端
* server.stop()                                    关闭服务端
* server.wait_for_start()                          等待直至服务端完全关闭（没有实现，空函数）
* server.restart()                                 重启服务端
* server.stop_exit()                               关闭服务端以及 MCDR，也就是退出整个程序
* server.exit()                                    关闭 MCDR
* server.is_server_running()                       返回一个 bool 代表服务端是否在运行
* server.is_server_startup()                       返回一个 bool 代表服务端是否已经启动完成
* server.is_rcon_running()                         返回一个 bool 代表 rcon 是否在运行(固定为false)
* server.execute(text, encoding=None)              发送字符串 text 至服务端的标准输入流，并自动在其后方追加一个\n
* server.say(text, encoding=None)                  使用 tellraw @a 来在服务端中广播消息
* server.tell(player, text, encoding=None)         使用 tellraw <player> 来在对玩家 <player> 发送消息
* server.reply(info, text, encoding=None)          向消息源发生消息: 如果消息来自玩家则调用 tell(info.player, text); 如果不是则调用 MCDR 的 logger 来将 text 告示至控制台
* server.load_plugin(plugin_name)                  加载名为 plugin_name 的插件。如果该插件已加载则重载它
* server.enable_plugin(plugin_name)                启用名为 plugin_name 的插件。该插件需已被禁用，即文件名后缀为 .py.disabled
* server.disable_plugin(plugin_name)	            禁用名为 plugin_name 的插件
* server.refresh_all_plugins()                     重载所有插件，加载新的插件并卸载移除的插件
* server.refresh_changed_plugins()                 重载所有文件有变化的插件，加载新的插件并卸载移除的插件
* server.get_plugin_list()                         返回一个 str 列表代表已加载的插件的文件名，如 ["pluginA.py", "pluginB.py"]
 */
