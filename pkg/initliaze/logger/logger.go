package logger

import (
	conf "bird/config"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"bird/utils/json"
)

func InitLogger(logConfig conf.LogInterface) error {
	err := initLogger(logConfig)
	if err != nil {
		return err
	}
	log.SetFlags(log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
	return nil
}

func ZnTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func initLogger(logConfig conf.LogInterface) error {
	var js string
	if logConfig.Env()  {
		js = fmt.Sprintf(`{
      "level": "%s",
      "encoding": "console",
      "outputPaths": ["stdout", "%s"],
      "errorOutputPaths": ["stdout", "%s"]
      }`, logConfig.GetLogLevel(),logConfig.GetLogPath(),logConfig.GetLogPath())
	} else {
		js = fmt.Sprintf(`{
      "level": "%s",
      "encoding": "console",
      "outputPaths": ["%s"],
      "errorOutputPaths": ["%s"]
      }`, logConfig.GetLogLevel(),logConfig.GetLogPath(),logConfig.GetLogPath())
	}

	cfg := zap.Config{}
	if err := json.Unmarshal([]byte(js), &cfg); err != nil {
		log.Fatal("init logger error: ", err)
		return err
	}
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = ZnTimeEncoder

	tlog, err := cfg.Build()
	if err != nil {
		log.Fatal("init logger error: ", err)
		return err
	}
	logConfig.SetSugaredLogger(tlog.Sugar())
	return nil
}

func Debug(logConfig conf.LogInterface,args ...interface{}) {
	logConfig.GetSugaredLogger().Debug(args...)
}

func Info(logConfig conf.LogInterface,args ...interface{}) {
	logConfig.GetSugaredLogger().Info(args...)
}

func Warn(logConfig conf.LogInterface,args ...interface{}) {
	SendMonitor2DingDing(logConfig.GetAlarmUrl(),fmt.Sprintf("%v",args))
	logConfig.GetSugaredLogger().Warn(args...)
}

func Error(logConfig conf.LogInterface,args ...interface{}) {
	SendMonitor2DingDing(logConfig.GetAlarmUrl(),fmt.Sprintf("%v",args))
	logConfig.GetSugaredLogger().Error(args...)
}

func Debugf(logConfig conf.LogInterface,format string, args ...interface{}) {
	logConfig.GetSugaredLogger().Debugf(format, args...)
}

func Infof(logConfig conf.LogInterface,format string, args ...interface{}) {
	logConfig.GetSugaredLogger().Infof(format, args...)
}

func Warnf(logConfig conf.LogInterface,format string, args ...interface{}) {
	SendMonitor2DingDing(logConfig.GetAlarmUrl(),fmt.Sprintf(format, args))
	logConfig.GetSugaredLogger().Warnf(format, args...)
}

func Errorf(logConfig conf.LogInterface,format string, args ...interface{}) {
	SendMonitor2DingDing(logConfig.GetAlarmUrl(),fmt.Sprintf(format, args))
	logConfig.GetSugaredLogger().Errorf(format, args...)
}

func SendMonitor2DingDing(alarmUrl,args string) {
	if len(alarmUrl) != 0 {
		http.Post(alarmUrl, "application/json", strings.NewReader(args))
	}
}