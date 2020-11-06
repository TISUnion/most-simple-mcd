package constant

const (
	CREATE_CONF_ERROR = "无法创建配置文件"

	READ_CONF_ERROR = "无法读取配置文件"

	WRITE_CONF_ERROR = "无法在配置文件中写入配置"

	PARSE_INI_CONF_ERROR = "解析配置文件失败"

	CREATE_LOG_FAILED = "创建日志文件失败"

	GET_CURRENT_PATH_FAILED = "获取当前执行目录失败"

	WRITE_LOG_FAILED = "日志写入失败"

	CREATE_CONF_FAILED_AND_ROLLBACK = "创建配置文件失败，回滚至默认配置"
)

// http错误码
const (
	HTTP_OK                   = 0
	HTTP_SYSTEM_ERROR         = 2009199999 // 系统错误
	HTTP_SYSTEM_ERROR_MESSAGE = "系统错误"
	HTTP_PARAMS_ERROR         = 2009190000 // 参数错误
	HTTP_PARAMS_ERROR_MESSAGE = "参数错误"
	PASSWORD_ERROR            = 2009190001
	PASSWORD_ERROR_MESSAGE    = "账号或密码错误"
	TOKEN_FAILED              = 2009190002
	TOKEN_FAILED_MESSAGE      = "token过期"
	PARSE_FILE_ERROR          = "上传文件获取错误"
	COPY_FILE_ERROR           = "复制文件错误"
	UNCOMPRESS_FILE_ERROR     = "解压文件错误"
)
