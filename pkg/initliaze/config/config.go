package config

import (
	conf "bird/config"
	"bird/pkg/initliaze/logger"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"time"
)

func InitConfig(path string)error{
	f,_ := os.Stat(path)
	conf.BirdConf.BirdConfig.LastModifyTime = f.ModTime().Unix()
	if _, err := toml.DecodeFile(path, conf.BirdConf); err != nil {
		fmt.Printf("toml decode err %v", err)
		panic(err)
	}
	// check
	if conf.BirdConf.BirdConfig.IsDebug {
		fmt.Println("config data is ",conf.BirdConf.BirdConfig)
	}
	go reload(path)
	return nil
}

func reload(path string){
	defer func() {
		if err := recover();err != nil {
			logger.Errorf(conf.BirdConf,"[reload] goroutine panic err %v",err)
		}
	}()
	ticker := time.NewTicker(time.Duration(conf.BirdConf.BirdConfig.Reload)*time.Second)
	for {
		select {
		case <- ticker.C:
			func(){
				fmt.Println("[reload] run")
				f,_ := os.Stat(path)
				conf.BirdConf.BirdConfig.Lock.Lock()
				defer conf.BirdConf.BirdConfig.Lock.Unlock()
				lastTime := f.ModTime().Unix()
				if lastTime > conf.BirdConf.BirdConfig.LastModifyTime {
					conf.BirdConf.BirdConfig.LastModifyTime = lastTime
					fmt.Println("[reload] run reload")
					if _, err := toml.DecodeFile(path, conf.BirdConf); err != nil {
						fmt.Printf("toml decode err %v", err)
					}
				}
			}()
		}
	}
}