package initliaze

import (
	conf "bird/config"
	"bird/pkg/initliaze/basic"
	"bird/pkg/initliaze/config"
	"bird/pkg/initliaze/logger"
	"errors"
	"fmt"
	"sync"
	"time"
)

func Initialize(path string)error{
	conf.BirdConf = &conf.Config{
		BirdConfig:&conf.BirdConfig{
			Lock:&sync.RWMutex{},
		},
	}
	if err := config.InitConfig(path);err != nil {
		return errors.New(fmt.Sprintf("%v [InitConfig] init config err:%v",time.Now().String(),err))
	}
	if err := logger.InitLogger(conf.BirdConf);err != nil {
		panic(err)
	}

	go basic.InitBasic()
	return nil
}
