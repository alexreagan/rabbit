package g

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path"
	"time"
)

func InitLog() {
	if viper.GetBool("log.enabled") == false {
		return
	}

	// set rotate log
	logPath := viper.GetString("log.path")
	logFile := path.Join(logPath, "stdout.log")
	rotateLog, _ := rotatelogs.New(
		logFile+".%Y%m%d",
		rotatelogs.WithLinkName(logFile),
		rotatelogs.WithMaxAge(time.Duration(24*7)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)

	// set log format
	log.SetOutput(rotateLog)
	log.SetReportCaller(true)
	log.SetFormatter(&nested.Formatter{
		CallerFirst:     true,
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:03:04",
	})

	// set log level
	switch viper.GetString("log.level") {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.Fatal("log conf only allow [info, debug, warn], please check your confguire")
	}

	return
}
