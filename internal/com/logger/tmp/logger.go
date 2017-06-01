package logger

import (
	"encoding/json"
	"fmt"
	"time"
	"strings"
	"magpie/internal/com/logger/tmp/logs"
	"magpie/internal/com/cfg"
	"magpie/internal/com/utils"
)

type Data map[string]interface{}

var (
	loggerMaps = map[string]*logs.BeeLogger{}
	logFileLevel = logs.LevelDebug

// 错误级别的文件后缀
	LevelExt = map[int]string{
		logs.LevelEmergency: "emergency",
		logs.LevelAlert: "alert",
		logs.LevelCritical: "critical",
		logs.LevelError: "error",
		logs.LevelWarning: "warning", // = LevelWarn
		logs.LevelNotice: "notice",
		logs.LevelInformational: "info", // = LevelInfo
		logs.LevelDebug: "debug", // = LevelTrace
	}
)

func init() {
	initLogConf()
}

/**
 * 日志文件的配置
 */
type LogFileParams struct {
	FileName string `json:"filename"`
	Level    int `json:"level"`
}

/**
 * 获取一个日志实例
 */
func getFileLogger(category string, level int) *logs.BeeLogger {
	key := utils.Md5Str(fmt.Sprintf("%s%d", category, level))

	if _, ok := loggerMaps[key]; !ok {
		// 日志缓冲区间长度
		cl, _ := cfg.Int("logger", "ChannelLen")
		chLen := int64(cl)
		if chLen <= 0 {
			loggerMaps[key] = logs.NewLogger(1000)
		}else {
			loggerMaps[key] = logs.NewLogger(chLen)
		}

		// 文件名格式
		curDate := time.Now().Format("20060102")
		p, _ := cfg.String("logger", "LogFilePath")
		logPath := strings.TrimRight(p, "/") + "/"
		filename := logPath + category + "." + LevelExt[level] + "." + curDate
		pb, _ := json.Marshal(&LogFileParams{
			FileName: filename,
			Level: level,
		})

		loggerMaps[key].SetLogger("file", string(pb))
		loggerMaps[key].SetLogger("console", "")
		loggerMaps[key].EnableFuncCallDepth(true)
		loggerMaps[key].SetLogFuncCallDepth(3) // default: 2
		//		loggerMaps[key].Async()
	}

	return loggerMaps[key]
}

func Info(category string, v ...interface{}) {
	if logs.LevelInformational > logFileLevel {
		return
	}

	getFileLogger(category, logs.LevelInformational).Info(generateFmtStr(1), formatData(v))
}

func Error(category string, v ...interface{}) {
	if logs.LevelError > logFileLevel {
		return
	}

	getFileLogger(category, logs.LevelError).Error(generateFmtStr(1), formatData(v))
}

func Debug(category string, v ...interface{}) {
	if logs.LevelDebug > logFileLevel {
		return
	}

	getFileLogger(category, logs.LevelDebug).Debug(generateFmtStr(1), formatData(v))
}

func generateFmtStr(n int) string {
	return strings.Repeat("\n%v", n) + "\n"
}

func initLogConf() {
	level, err := cfg.Int("logger", "LogFileLevel")
	if err != nil {
		return
	}

	if _, ok := LevelExt[level]; !ok {
		return
	}

	logFileLevel = level
}

func SetLogFileLevel(level int) {
	logFileLevel = level
}

// 格式化日志记录
func formatData(v interface{}) interface{} {
	str, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return v
	}

	return string(str)
}