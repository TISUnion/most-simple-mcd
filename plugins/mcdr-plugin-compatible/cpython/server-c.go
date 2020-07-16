package cpython

/**
#include <Python.h>
#include "Utils.h"

#define NAME "name"
#define ID "id"
#define PORT "port"
#define MEMORY "memory"
#define VERSION "version"
#define SIDE "side"
#define COMMENT "comment"

extern void serverInfo(char *id);
extern void start(char *id);
extern void stop(char *id);
extern void restart(char *id);
extern void stopExit(char *id);
extern void exit(char *id);
extern int isServerRunning(char *id);
extern int isServerStartup(char *id);
extern int isRconRunning(char *id);
extern void execute(char *id, char *text);
extern void say(char *id, char *text);
extern void tell(char *id, char *player, char *text);
extern void reply(char *id, char *player, char *text);

PyObject *start(PyObject *self, PyObject *args)
{
  char *id = PyObject_GetAttrString(self, "id");
  start(id);
  return Py_None;
}

PyObject *stop(PyObject *self, PyObject *args)
{
  char *id = PyObject_GetAttrString(self, "id");
  stop(id);
  return Py_None;
}

PyObject *restart(PyObject *self, PyObject *args)
{
  char *id = PyObject_GetAttrString(self, "id");
  restart(id);
  return Py_None;
}

PyObject *stopExit(PyObject *self, PyObject *args)
{
  char *id = PyObject_GetAttrString(self, "id");
  stopExit(id);
  return Py_None;
}

PyObject *exit(PyObject *self, PyObject *args)
{
  char *id = PyObject_GetAttrString(self, "id");
  exit(id);
  return Py_None;
}

PyObject *isServerRunning(PyObject *self, PyObject *args)
{
  char *id = PyObject_GetAttrString(self, "id");
  int res = isServerRunning(id);
  return (PyObject *)Py_BuildValue("i", res);
}

PyObject *isServerStartup(PyObject *self, PyObject *args)
{
  char *id = PyObject_GetAttrString(self, "id");
  int res = isServerRunning(id);
  return (PyObject *)Py_BuildValue("i", res);
}

PyObject *isRconRunning(PyObject *self, PyObject *args)
{
  char *id = PyObject_GetAttrString(self, "id");
  int res = isRconRunning(id);
  return (PyObject *)Py_BuildValue("i", res);
}

PyObject *execute(PyObject *self, PyObject *args)
{
  char *id = PyObject_GetAttrString(self, "id");
  char *text = PyObject_GetAttrString(self, "text");
  int res = execute(id, text);
  return Py_None;
}

PyObject *say(PyObject *self, PyObject *args)
{
  char *id = PyObject_GetAttrString(self, "id");
  char *text = PyObject_GetAttrString(self, "text");
  int res = say(id, text);
  return Py_None;
}

PyObject *tell(PyObject *self, PyObject *args)
{
  char *id = PyObject_GetAttrString(self, "id");
  char *text = PyObject_GetAttrString(self, "text");
  char *player = PyObject_GetAttrString(self, "player");
  int res = tell(id, player, text);
  return Py_None;
}

PyObject *reply(PyObject *self, PyObject *args)
{
  char *id = PyObject_GetAttrString(self, "id");
  char *text = PyObject_GetAttrString(self, "text");
  char *player = PyObject_GetAttrString(self, "player");
  int res = reply(id, player, text);
  return Py_None;
}

PyObject *GetServer(char *id)
{
	CreateClass(id, NULL);
}
*/
import "C"

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
