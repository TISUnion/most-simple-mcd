package constant

// 日志常量
const (
	DEFAULT_CHANNEL               = "default"
	EVERYDAY_JOB_NAME             = "everyday-add-log"
	LOG_SAVE_INTERVAL_EVERYDAY    = "59 23 * * *"   // cron每日表达式
	LOG_SAVE_INTERVAL_TWICEDAY    = "59 23 1/2 * ?" // cron每隔2天表达式
	LOG_SAVE_INTERVAL_EVERYMONDAY = "59 23 * * 0"   // cron每周一表达式
	LOG_SAVE_INTERVAL_EVERYMONTH  = "59 23 1 * ?"   // cron每月1日表达式
	LOG_DEBUG                     = "debug"
	LOG_INFO                      = "info"
	LOG_ERROR                     = "error"
	LOG_WARNING                   = "warning"
	LOG_FATAL                     = "fatal"
	LOG_FORMAT                    = "%s [%s]: %s"
	LOG_CODELINE_FORMAT           = "%s [%s] %s : %s"
	LOG_TIME_FORMAT               = "2006-01-02 15:04:05.000000"
)
