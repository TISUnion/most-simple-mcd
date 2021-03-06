package constant

import "time"

// 配置名常量
const (
	IS_RELOAD_CONF                       = "config_auto_reload"          // 自动加载配置文件
	RELOAD_CONF_INTERVAL                 = "config_auto_reload_interval" // 自动加载配置文件间隔，单位：毫秒
	CONF_PATH                            = "config_path"                 // 配置文件地址
	IS_MANAGE_HTTP                       = "http_manage_server"          // 启动管理后台
	MANAGE_HTTP_SERVER_PORT              = "http_manage_server_port"     // 管理后台服务端口
	LOG_PATH                             = "log_path"                    // 日志写入目录
	LOG_SAVE_INTERVAL                    = "log_interval"                // 日志保存间隔，例如: 每2天对久日志压缩，日志写入新日志中
	WORKSPACE                            = "workspace"                   // 工作目录
	TMP_PATH                             = "tmp_path"                    // 临时文件存放目录
	I18N                                 = "i18n"                        // 国际化
	IS_AUTO_CHANGE_MC_SERVER_REPEAT_PORT = "mc_server_port_auto_change"  // 是否自动更换mc服务器重复使用的端口
	MONITOR_INTERVAL                     = "monitor_interval"            // 进程资源读取间隔，单位为秒
	WEBSOCKET_HOST                       = "websocket_host"              // websocket 连接host
)

// 配置说明
const (
	IS_RELOAD_CONF_DESCREPTION                       = "自动加载配置文件"
	RELOAD_CONF_INTERVAL_DESCREPTION                 = "自动加载配置文件间隔，单位：毫秒"
	CONF_PATH_DESCREPTION                            = "配置文件地址"
	IS_MANAGE_HTTP_DESCREPTION                       = "启动管理后台"
	MANAGE_HTTP_SERVER_PORT_DESCREPTION              = "管理后台服务端口"
	LOG_PATH_DESCREPTION                             = "日志写入目录"
	LOG_SAVE_INTERVAL_DESCREPTION                    = "日志保存间隔，例如: 每2天对旧日志压缩，日志写入新日志中"
	IS_AUTO_CHANGE_MC_SERVER_REPEAT_PORT_DESCREPTION = "是否自动更换mc服务器重复使用的端口"
	WORKSPACE_DESCREPTION                            = "工作目录"
	I18N_DESCREPTION                                 = "国际化"
	MONITOR_INTERVAL_DESCREPTION                     = "进程资源读取间隔，单位: 秒"
	TMP_PATH_DESCREPTION                             = "临时文件存放目录"
	WEBSOCKET_HOST_DESCREPTION                       = "websocket的连接host"
)

// 配置常量
const (
	RELOAD_CONF_JOB_NAME = "reload-conf"
	// 配置覆盖优先级
	CONF_SYSTEM_LEVEL      = 4
	CONF_TERMINAL_LEVEL    = 3
	CONF_ENVIRONMENT_LEVEL = 2
	CONF_FILE_LEVEL        = 1
	CONF_DEFAULT_LEVEL     = 0
)

// 日志常量
const (
	DEFAULT_LOG_NAME              = "default"
	EVERYDAY_JOB_NAME             = "everyday-add-log"
	LOG_SAVE_INTERVAL_EVERYDAY    = "0 0 * * *"   // cron每日表达式
	LOG_SAVE_INTERVAL_TWICEDAY    = "0 0 1/2 * ?" // cron每隔2天表达式
	LOG_SAVE_INTERVAL_EVERYMONDAY = "0 0 * * 0"   // cron每周一表达式
	LOG_SAVE_INTERVAL_EVERYMONTH  = "0 0 1 * ?"   // cron每月1日表达式
	LOG_DEBUG                     = "debug"
	LOG_INFO                      = "info"
	LOG_ERROR                     = "error"
	LOG_WARNING                   = "warning"
	LOG_FATAL                     = "fatal"
	LOG_FORMAT                    = "%s [%s]: %s\r\n"
	LOG_CODELINE_FORMAT           = "%s [%s] %s : %s\r\n"
	LOG_TIME_FORMAT               = "2006-01-02 15:04:05.000000"
)

// mc服务端常量
const (
	EULA_FILE_NAME       = "eula.txt"
	EULA                 = "eula"
	TRUE_STR             = "true"
	MC_CONF_NAME         = "server.properties"
	MC_PORT_TEXT         = "server-port"
	MC_SERVER_DIR        = "minecraft-servers"
	MC_SERVER_BACK       = "minecraft-servers-back"
	MC_SERVER_DB_KEY     = "minecraft:server:configs"
	MC_LAST_UTF8_VERSION = "1.12"
	MC_ALL_PLAYER        = "@a"
	MC_DEFAULT_PORT      = 25565
	MC_DEFAULT_MEMORY    = 1024
	UUID_LENGTH          = 36
	MAX_RESIVE_BUFF_SIZE = 1024

	MC_STATE_STOPING  = -2
	MC_STATE_STARTIND = -1
	MC_STATE_STOP     = 0
	MC_STATE_START    = 1

	MC_SERVER_START   = 1
	MC_SERVER_STOP    = 2
	MC_SERVER_RESTART = 3

	MC_PLUGIN_START = 1
	MC_PLUGIN_STOP  = 2

	MC_OPEN_CALLBACK  = "open"
	MC_CLOSE_CALLBACK = "close"
	MC_SAVE_CALLBACK  = "save"

	MC_STDOUT     = 0
	MC_SYSTEM_OUT = 1
)

// DB常量
const (
	DEFAULT_DATABASE_NAME = "default-db"
)

// web后台常量
const (
	LOACALHOST_URL              = "http://localhost"
	COMPRESS_FILE_NAME          = "webfile.zip"
	Web_FILE_DIR_NAME           = "web-static-file"
	DEFAULT_ACCOUNT_DB_KEY      = "web:admin:account"
	DEFAULT_TOKEN_DB_KEY        = "web:admin:token"
	DEFAULT_TOKEN_DB_KEY_EXPIRE = 4 * time.Hour
	DEFAULT_ACCOUNT             = "steve"
	DEFAULT_PASSWORD            = "123456"
	QUERY_ID                    = "id"
	QUERY_NAME                  = "name"
	QUERY_TYPE                  = "type"

	LOG_TYPE_SERVER  = "1"
	LOG_TYPE_GIN     = "2"
	LOG_TYPE_DEFAULT = "3"

	JAR_SUF           = ".jar"
	TOKEN_HEADER_NAME = "X-Token"
)

// 各模块常量
const (
	TIME_FORMAT                     = "2006-01-02 15:04:05"
	DEFAULT_MANAGE_HTTP_SERVER_PORT = 80
	PLUGIN_COMMAND_TYPE             = 1 //插件运行命令
	SERVER_COMMAND_TYPE             = 2 //服务端运行命令
	ALL_COMMAND_TYPE                = 3 //插件、服务端都运行
	LOG_DIR                         = "logs"
	PLUGIN_DIR                      = "plugins"
	GIN_LOG_NAME                    = "gin-server"
	UPLOAD_FILE_NAME                = "file"
	UPLOAD_PORT_TEXT                = "port"
	UPLOAD_MEMORY_TEXT              = "memory"
	UPLOAD_NAME_TEXT                = "name"
	UPLOAD_SIDE_TEXT                = "side"
	UPLOAD_COMMENT_TEXT             = "comment"

	OS_WINDOWS = "windows"
	OS_LINUX   = "linux"
	OS_DARWIN  = "darwin"

	CLI_MODE  = "-cli"
	PIX       = "-"
	CMD_ERROR = "解析命令行出错"

	COMPARE_GT = 1
	COMPARE_EQ = 0
	COMPARE_LT = -1
)

// mc事件常量
const (
	ON_LOAD = iota + 1
	ON_UNlOAD
	ON_INFO
	ON_USER_INFO
	ON_PLAYER_JOINED
	ON_PLAYER_LEFT
	ON_DEATH
	ON_PLAYER_MADE_ADVANCEMENT
	ON_SERVER_STARTUP
	ON_SERVER_STOP
	ON_MCDR_STOP
	VERSION
	GAME_TYPE
	GAME_START
	GAME_SAVE
)
