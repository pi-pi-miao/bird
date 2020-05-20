package config

import (
	"go.uber.org/zap"
	"sync"
)

var (
	BirdConf *Config
	BirdCh   = make(chan struct{},100)
	BasicChStop  = make(chan struct{})
)

type Config struct{
	BirdConfig *BirdConfig
}

type BirdConfig struct {
	Lock         *sync.RWMutex
	Addr         string   `toml:"addr" json:"addr"`
	ReadTimeout  int64    `toml:"read_time_out" json:"read_time_out"`
	WriteTimeout int64    `toml:"write_time_out" json:"write_time_out"`
	ShutdownTime int64    `toml:"shutdown_time" json:"shutdown_time"`
	BirdConfAddr string   `toml:"bird_conf_addr" json:"bird_conf_addr"`
	Reload       int      `toml:"reload" json:"reload"`
	LastModifyTime int64
	// log config
	AlarmUrl     string   `toml:"alarm_url" json:"alarm_url"`
	IsDebug      bool     `toml:"is_debug" json:"is_debug"`
	LogLevel     string   `toml:"log_level" json:"log_level"`
	LogPath      string   `toml:"log_pah" json:"log_pah"`
	SugaredLogger *zap.SugaredLogger
	// auth
	Name         string   `toml:"name" json:"name"`
	Password     string   `toml:"password" json:"password"`
	PasswordSlat string   `toml:"password_slat" json:"password_slat"`
}


func (c *Config)Env()bool{
	return c.BirdConfig.IsDebug
}

func (c *Config)GetLogLevel()string{
	return c.BirdConfig.LogLevel
}

func (c *Config)GetLogPath()string {
	return c.BirdConfig.LogPath
}

func (c *Config)SetSugaredLogger(sugaredLogger *zap.SugaredLogger) {
	c.BirdConfig.SugaredLogger = sugaredLogger
	return
}

func (c *Config)GetSugaredLogger()*zap.SugaredLogger{
	return c.BirdConfig.SugaredLogger
}

func (c *Config)GetAlarmUrl()string{
	return c.BirdConfig.AlarmUrl
}