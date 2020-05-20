/*
  all server
*/

package bird_proxy

import (
	conf "bird/config"
	"bird/pkg/initliaze/logger"
	"context"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type BirdProxy struct {
	http.Handler
	Addr         string   `json:"addr"`
	ReadTimeout  int64    `json:"read_time_out"`
	WriteTimeout int64    `json:"write_time_out"`
	ShutdownTime int64    `json:"shutdown_time"`
	Informer     chan struct{}
	// log config
	AlarmUrl     string   `json:"alarm_url"`
	IsDebug      bool     `json:"is_debug"`
	LogLevel     string   `json:"log_level"`
	LogPath      string   `json:"log_pah"`
	SugaredLogger *zap.SugaredLogger
}

func (p *BirdProxy)Env()bool {
	return p.IsDebug
}

func (p *BirdProxy)GetLogLevel()string {
	return p.LogLevel
}

func (p *BirdProxy)GetLogPath()string{
	return p.LogPath
}

func (p *BirdProxy)SetSugaredLogger(sugaredLogger *zap.SugaredLogger) {
	p.SugaredLogger = sugaredLogger
	return
}

func (p *BirdProxy)GetSugaredLogger()*zap.SugaredLogger{
	return p.SugaredLogger
}

func (p *BirdProxy)GetAlarmUrl()string {
	return p.AlarmUrl
}

func (p *BirdProxy)Server(){
	server := &http.Server{
		Addr:         p.Addr,
		Handler:      p.Handler,
		ReadTimeout:  time.Duration(p.ReadTimeout) *time.Second,
		WriteTimeout: time.Duration(p.WriteTimeout) * time.Second,
	}
	go func() {
		<- p.Informer
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Errorf(conf.BirdConf,"Could not gracefully shutdown the server: %v\n", err)
		}
		close(p.Informer)
	}()
	if err := server.ListenAndServe(); err != nil {
		logger.Errorf(conf.BirdConf,"server2 is running err %v",err)
	}
}