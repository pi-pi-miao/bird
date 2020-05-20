package bird_controller

import (
	conf "bird/config"
	"bird/pkg/bird_basic_service"
	"bird/pkg/initliaze/logger"
	"bird/utils/cache"
	"bird/utils/json"
	"io/ioutil"
	"net/http"
)

type LoginData struct {
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

// login and return token
// first login then edit config
func Login(w http.ResponseWriter,r *http.Request){
	loginBody := &LoginData{}
	body,err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warnf(conf.BirdConf,"login err,request is :%v",*r)
		response(w,"read body err",ReadRequestError)
		return
	}
	if err := json.Unmarshal(body,loginBody);err != nil {
		logger.Warnf(conf.BirdConf,"%v login err,request is :%v",string(body),*r)
		response(w,"unmarshal body err",UnmarshalRequestBodyError)
		return
	}
	if value,ok := cache.BirdAccountCache.Get(loginBody.Name);!ok {
		logger.Warnf(conf.BirdConf,"%v login err,request is :%v",loginBody.Name,*r)
		response(w,"the Name is not register",PasswordError)
		return
	}else {
		if value.(*LoginData).Password != loginBody.Password {
			logger.Warnf(conf.BirdConf,"%v login err,request is :%v",loginBody.Name,*r)
			response(w,"password is not right",PasswordError)
			return
		}
	}

	// todo test hash set   cache.BirdAccountCache.Len()
	token := bird_basic_service.CreateToken(loginBody.Name,loginBody.Password)
	cache.BirdAccountCache.Set(loginBody.Name,&LoginData{
		Name:     loginBody.Name,
		Password: loginBody.Password,
		Token:    token,
	})
	response(w,token,Success)
	return
}


func response(w http.ResponseWriter,data,code string){
	responseData,_ := json.Marshal(BirdResponse{
		Code: code,
		Data: data,
	})
	w.Write(responseData)
	return
}