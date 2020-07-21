package cpython

/**
//#cgo pkg-config:python3
//#include "Utils.h"
//

*/
//import "C"

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
