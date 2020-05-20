package config

import "go.uber.org/zap"

type LogInterface interface {
	Env()bool
	GetLogLevel()string
	GetLogPath()string
	SetSugaredLogger(*zap.SugaredLogger)
	GetSugaredLogger()*zap.SugaredLogger
	GetAlarmUrl()string
}
