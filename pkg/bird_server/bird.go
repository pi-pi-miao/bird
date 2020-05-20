package bird_server

import (
	"bird/apis/bird_apis"
	conf "bird/config"
	"bird/pkg/bird_proxy"
	"bird/pkg/initliaze"
	"bird/pkg/initliaze/logger"
	"bird/utils/cache"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)


func Run(path string){
	if err := initliaze.Initialize(path);err != nil {
		panic(err)
	}
	server := &bird_proxy.BirdProxy{
		Addr:         conf.BirdConf.BirdConfig.Addr,
		ReadTimeout:  conf.BirdConf.BirdConfig.ReadTimeout,
		WriteTimeout: conf.BirdConf.BirdConfig.WriteTimeout,
		ShutdownTime: conf.BirdConf.BirdConfig.ShutdownTime,
		Informer:make(chan struct{}),
		Handler : http.NewServeMux(),
		AlarmUrl:     conf.BirdConf.BirdConfig.AlarmUrl,
		IsDebug:      conf.BirdConf.BirdConfig.IsDebug,
		LogLevel:     conf.BirdConf.BirdConfig.LogLevel,
		LogPath:      conf.BirdConf.BirdConfig.LogPath,
	}
	go func(message *bird_proxy.BirdProxy) {
		port := strings.Split(conf.BirdConf.BirdConfig.Addr,":")[1]
		cache.ServiceInformerCache.Set(fmt.Sprintf(":%v",port),server)
		bird_apis.BirdApis(server.Handler)
		server.Server()
	}(server)
	logger.Debugf(conf.BirdConf,"bird is running at %v",conf.BirdConf.BirdConfig.Addr)

	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<- exit
	server.Informer <- struct{}{}
	<- server.Informer
	logger.Error(conf.BirdConf,"gracefully shutdown the http server...")
}

