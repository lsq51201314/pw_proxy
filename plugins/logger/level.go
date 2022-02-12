package logger

const (
	Log_Level_DEBUG = 0
	Log_Level_INFO  = 1
	Log_Level_WARN  = 2
	Log_Level_ERROR = 3
	Log_Level_FATAL = 4
	Log_Level_TRACE = 5
)

func GetLevelStr(code int) string {
	switch code {
	case Log_Level_DEBUG:
		return "调试信息"
	case Log_Level_INFO:
		return "提示信息"
	case Log_Level_WARN:
		return "警告信息"
	case Log_Level_ERROR:
		return "错误信息"
	case Log_Level_FATAL:
		return "致命错误"
	case Log_Level_TRACE:
		return "跟踪信息"
	default:
		return "未知信息"
	}
}
