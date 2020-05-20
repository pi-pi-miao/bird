package basic

import (
	"bird/pkg/bird_basic_service"
)

func InitBasic(){
	bird_basic_service.
		NewBirdConfigBasic().
		Basic()
	return
}