package log

const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

//大小为8的数组
var levelPrefix = [LevelDebug + 1]string{"[M] ", "[A] ", "[C] ", "[E] ", "[W] ", "[N] ", "[I] ", "[D] "}

//LogLevelMsg- 日志等级输出
func LogLevelMsg(level int) string {
	if level < 0 || level > LevelDebug {
		return ""
	}
	return levelPrefix[level]
}