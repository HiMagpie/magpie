package logger

import (
	"strings"
	"github.com/KowloonZh/flog"
	"fmt"
	"magpie/internal/com/cfg"
)

type Data map[string]interface{}

var (
	logger *flog.Flog
	logLevel = flog.LEVEL_DEBUG
	LevelExt = map[int]string{
		flog.LEVEL_DEBUG: "debug",
		flog.LEVEL_INFO: "info",
		flog.LEVEL_WARNING: "warning",
		flog.LEVEL_ERROR: "error",
	}
)

func init() {
	initLogConf()
}

func initLogConf() {
	if _, ok := LevelExt[cfg.C.Log.Level]; !ok {
		fmt.Println("日志级别错误!")
		return
	}

	logger = flog.New(cfg.C.Log.Path)
	logger.LogMode = flog.LOGMODE_CATE_LEVEL
	logger.Level = logLevel
	logger.OpenConsoleLog = true
	logger.LogFlags = []int{flog.LF_DATETIME, flog.LF_LEVEL, flog.LF_SHORTFILE} // flog.LF_CATE
	logger.LogFlagSeparator = "`"
	logger.LogFunCallDepth = 4
}

func Info(category string, v ...interface{}) {
	logger.Info(category, v...)
}

func Error(category string, v ...interface{}) {
	logger.Error(category, v...)
}

func Debug(category string, v ...interface{}) {
	logger.Debug(category, v...)
}

func SetLogLevel(level int) {
	logLevel = level
}

func Format(d ...interface{}) string {
	s := ""
	for i, _ := range d {
		if i % 2 == 1 {
			s += fmt.Sprintf("%v=%v`", d[i - 1], d[i])
		}
	}

	if len(d) % 2 != 0 {
		s += fmt.Sprintf("%v", d[len(d) - 1])
	}

	return strings.TrimRight(s, "`")
}

func colorful(s string, status string) string {
	out := ""
	switch status {
	case "succ":
		out = "\033[32;1m" // Blue
	case "fail":
		out = "\033[31;1m" // Red
	case "warn":
		out = "\033[33;1m" // Yellow
	case "note":
		out = "\033[34;1m" // Green
	default:
		out = "\033[0m"// Default
	}
	return out + s + "\033[0m"
}