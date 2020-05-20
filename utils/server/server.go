package server

import (
	"bird/config"
	"bird/pkg/bird_proxy"
	"bird/pkg/initliaze/logger"
	"bird/utils/cache"
)

func Server(readTimeOut,writeTimeOut,shutdownTimeOut int64,addr,url string,isDebug bool,level,path string){
	go func(readTimeOut,writeTimeOut,shutdownTimeOut int64,addr,url string,isDebug bool,level,path string) {
		server := &bird_proxy.BirdProxy{
			Addr:         addr,
			ReadTimeout:  readTimeOut,
			WriteTimeout: writeTimeOut,
			ShutdownTime: shutdownTimeOut,
			Informer:make(chan struct{}),

			AlarmUrl:     url,
			IsDebug:      isDebug,
			LogLevel:     level,
			LogPath:      path,
		}
		server.Handler = server
		err := logger.InitLogger(server)
		if err != nil {
			logger.Errorf(config.BirdConf,"[Server] start port:%v err:%v",server.Addr,err)
			return
		}
		cache.ServiceInformerCache.Set(addr,server)
		server.Server()
	}(readTimeOut,writeTimeOut,shutdownTimeOut,addr,url,isDebug,level,path )
}
