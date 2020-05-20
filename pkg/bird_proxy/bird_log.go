package bird_proxy

import "bird/pkg/initliaze/logger"

func (p *BirdProxy)Log(){
	err := logger.InitLogger(p)
	if err != nil {

	}
}