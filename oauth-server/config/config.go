package config

import (
	"fmt"
	"github.com/clintwan/armory"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"

	"github.com/tidwall/gjson"
	)

var config gjson.Result

type logrusClassicTextFormatter struct {
	base logrus.JSONFormatter
}

func (f *logrusClassicTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	colored := Get("app.environment").String() == "staging"
	levelColor := "36m"
	if entry.Level == logrus.DebugLevel {
		levelColor = "32m"
	} else if entry.Level == logrus.ErrorLevel {
		levelColor = "31m"
	}

	caller := entry.Caller
	timeNow := entry.Time.Format("15:04:05")
	level := entry.Level.String()
	filePath := caller.File + ":" + strconv.Itoa(caller.Line)
	funcTitle := caller.Function
	msg := entry.Message

	if entry.Data["flag"] == "gorm" {
		filePath = entry.Data["file"].(string)
		funcTitle = entry.Data["flag"].(string) + ".func"
		msg = entry.Data["msg"].(string)
	}

	r := fmt.Sprintf("[%s] [%-5s] %s %s %s\n",
		timeNow, level, filePath, funcTitle, msg,
	)
	if colored {
		r = fmt.Sprintf("[%s] \033[%s[%-5s]\033[0m %s \033[34m%s\033[0m \033[36;1m%s\033[0m\n",
			timeNow, levelColor, level, filePath, funcTitle, msg,
		)
	}
	// show all color
	// tr := []string{}
	// for idx := 0; idx < 100; idx++ {
	// 	tr = append(tr, fmt.Sprintf("\033[%dm [%dm] \033[0m", idx, idx))
	// 	tr = append(tr, fmt.Sprintf("\033[%d;1m [%d;1m] \033[0m", idx, idx))
	// }
	// r = strings.Join(tr, "")
	return []byte(r), nil
}

// Init Init
func Init(buf []byte) {
	config = gjson.ParseBytes(buf)
	initLogger()
}

func initLogger() {
	env := Get("app.environment").String()
	if env == "production" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(armory.Log.DailyRotateLog(armory.Pilot.AppPath(Get("log.runtimeLog").String())))
	} else {
		logrus.SetFormatter(&logrusClassicTextFormatter{})
		logrus.SetOutput(os.Stdout)
	}

	logLevel := logrus.InfoLevel
	switch Get("log.level").String() {
	case "panic":
		logLevel = logrus.PanicLevel
	case "fatal":
		logLevel = logrus.FatalLevel
	case "error":
		logLevel = logrus.ErrorLevel
	case "warn":
		logLevel = logrus.WarnLevel
	case "info":
		logLevel = logrus.InfoLevel
	case "debug":
		logLevel = logrus.DebugLevel
	case "trace":
		logLevel = logrus.TraceLevel
	}
	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)
}

// Get Get
func Get(path string) gjson.Result {
	return config.Get(path)
}
