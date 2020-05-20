package bird_proxy

import (
	"bird/pkg/bird_basic_service"
	"bird/pkg/initliaze/logger"
	"bird/utils/cache"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	client = http.Client{}
)

func (p *BirdProxy)getDomain(w http.ResponseWriter,r *http.Request){
	fmt.Println("[getDomain] request host is",r.Host,"url",r.RequestURI)
	v,ok := cache.BirdCache.Get(r.Host)
	if !ok {
		writeJson(w,UrlNotFound,"illegal request")
		return
	}

	config := *v.(*bird_basic_service.BirdConfMessage)
	if config .IsReverseProxy {
		_,ok := config.ReverseProxyService.Method.Get(r.Method)
		if !ok {
			writeJson(w,MethodError,"illegal request")
			return
		}
		p.reverseProxyHttp(config,w,r)
	}

}

func (p *BirdProxy)reverseProxyWebSocket(config bird_basic_service.BirdConfMessage,w http.ResponseWriter,r *http.Request){

}


func (p *BirdProxy)reverseProxyHttp(config bird_basic_service.BirdConfMessage,w http.ResponseWriter,r *http.Request){
	req,err := http.NewRequest(r.Method,config.ReverseProxyService.Addr,r.Body)
	if err != nil {
		logger.Errorf(p,"[reverseProxy] newRequest err:%v",err)
		writeJson(w,ServiceError,"illegal request")
		return
	}
	for k,_ := range r.Header {
		req.Header[k] = r.Header[k]
	}
	resp,err := client.Do(req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		logger.Errorf(p,"[reverseProxy] client Do err %v",err)
		writeJson(w,ServiceError,"illegal request")
		return
	}

	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf(p,"[reverseProxy] readAll resp body err:%v",err)
		writeJson(w,ServiceError,"illegal request")
		return
	}

	for k,_ := range w.Header() {
		delete(w.Header(),k)
	}
	for k,_ := range resp.Header {
		w.Header()[k] = resp.Header[k]
	}
	w.Write(body)
	return
}

func (p *BirdProxy)proxyHttp(config bird_basic_service.BirdConfMessage,w http.ResponseWriter,r *http.Request){

}