package bird_proxy

import (
	"bird/utils/json"
	"net/http"
)

const (
	UrlNotFound = "404"
	MethodError = "400"
	ServiceError = "500"
)

type Response struct {
	Code string         `json:"code"`
	Data string         `json:"data"`
}


func writeJson(w http.ResponseWriter,code,data string){
	res,_ := json.Marshal(Response{
		Code: code,
		Data: data,
	})
	w.Write(res)
	return
}