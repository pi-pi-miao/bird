package bird_basic_service

import (
	"bird/config"
	"bird/pkg/bird_proxy"
	"bird/pkg/initliaze/logger"
	"bird/utils/cache"
	"bird/utils/json"
	"bird/utils/server"
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/cornelk/hashmap"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	hosts = "/home/wanglei/juejin/example/writefile/file"
)


// Full quantity update
type BirdConfMessage struct {
	Id          string   `json:"id"`
	Code 		string   `json:"code"`
	Data 		string   `json:"data"`
	Method      string   `json:"method"`     // get update delete

	BirdDomain   string   `json:"domain"`
	Port         string    `json:"port,omitempty"`

	// server config
	ReadTimeout  int64    `toml:"read_time_out" json:"read_time_out"`
	WriteTimeout int64    `toml:"write_time_out" json:"write_time_out"`
	ShutdownTime int64    `toml:"shutdown_time" json:"shutdown_time"`

	// log config
	AlarmUrl     string   `toml:"alarm_url" json:"alarm_url"`
	IsDebug      bool     `toml:"is_debug" json:"is_debug"`
	LogLevel     string   `toml:"log_level" json:"log_level"`
	LogPath      string   `toml:"log_pah" json:"log_pah"`

	IsReverseProxy bool    `json:"is_reverse_proxy "`
	ReverseProxyService Service `json:"reverse_proxy_service"`
	Route      Service  `json:"route,omitempty"`
}

// service
type Service struct {
	Addr          string    `json:"addr"`
	Method        *hashmap.HashMap  `json:"method"`
	IsWebSocket   bool      `json:"is_websocket"`
}


type birdAuth struct {
	Name        string   `json:"name"`
	Password    string   `json:"password"`
}

type BirdConfigBasic struct {
	lock *sync.Mutex
	conn net.Conn
	ch   chan bool
	dataCh chan []byte
	buffer *strings.Builder
	fileBuffer *strings.Builder
}

func NewBirdConfigBasic()*BirdConfigBasic {
	basic :=  &BirdConfigBasic{
		lock:       &sync.Mutex{},
		ch:         make(chan bool,1),
		dataCh:     make(chan []byte,1000),
		buffer:     &strings.Builder{},
		fileBuffer: &strings.Builder{},
	}
	basic.ch <- true
	return basic
}

func (b *BirdConfigBasic)Basic(){
	for {
		<- b.ch
		config.BirdConf.BirdConfig.Lock.RLock()
		conn,err := net.Dial("tcp",config.BirdConf.BirdConfig.BirdConfAddr)
		if err != nil {
			config.BirdConf.BirdConfig.Lock.RUnlock()
			b.dial("[Basic] server dial basic conf center err :%v",err)
			continue
		}
		config.BirdConf.BirdConfig.Lock.RUnlock()
		b.conn = conn
		go b.read()
		go b.getConfig()
		b.auth()
	}
}

func (b *BirdConfigBasic)auth(){
	config.BirdConf.BirdConfig.Lock.RLock()
	body,_ := json.Marshal(&birdAuth{
		Name:config.BirdConf.BirdConfig.Name,
		Password:config.BirdConf.BirdConfig.Password,
	})
	config.BirdConf.BirdConfig.Lock.RUnlock()
	if _,err := b.conn.Write(body);err != nil {
		b.conn.Close()
		b.dial("[dial] basic dial again and this conn write err %v",err)
	}
}

func (b *BirdConfigBasic)read(){
	defer func() {
		if err := recover();err != nil {
			logger.Errorf(config.BirdConf,"[irdConfigBasic read] goroutine panic err:%v",err)
		}
	}()
	sizeData := make([]byte, 2)
	for {
		select {
		case <- config.BasicChStop:
			logger.Error(config.BirdConf,"[Basic] server is finishing")
			return
		default:
		}
		if _,err := io.ReadFull(b.conn, sizeData); err != nil {
			b.conn.Close()
			logger.Errorf(config.BirdConf,"[Basic] read conn err %v",err)
			b.dial("[dial] basic dial again and this conn read header err :%v",err)
			return
		}
		data := make([]byte, binary.LittleEndian.Uint16(sizeData))
		if _, err := io.ReadFull(b.conn, data); err != nil {
			b.conn.Close()
			logger.Errorf(config.BirdConf,"[Basic] read conn err %v",err)
			b.dial("[dial] basic dial agin and this conn read body err :%v",err)
			return
		}
		b.dataCh <- data
	}
}

func (b *BirdConfigBasic)dial(format string, args ...interface{}){
	time.Sleep(5*time.Second)
	b.ch <- true
	logger.Warnf(config.BirdConf,format,args)
	//logger.Errorf(config.BirdConf,format,args)
}

func (b *BirdConfigBasic)getConfig(){
	for v := range b.dataCh {
		go func() {
			defer func() {
				if err := recover();err != nil {
					logger.Errorf(config.BirdConf,"[writeDomain] goroutine panic %v",err)
				}
			}()
			b.judge(v)
		}()
	}
}

func (b *BirdConfigBasic)judge(data []byte){
	message := BirdConfMessage{
		ReverseProxyService:Service{
			Method:      hashmap.New(9),
		},
		Route:Service{
			Method:      hashmap.New(9),
		},
	}
	err := json.Unmarshal(data,&message)
	if err != nil {
		logger.Errorf(config.BirdConf,"[judge] get bird conf err :%v",err)
		b.response(fmt.Sprintf("[judge] unmarshal   bird conf err :%v",err),
			"update_domain",
					message.Id,
			message)
		return
	}
	if message.ReadTimeout <= 0 || message.WriteTimeout <= 0 || message.ShutdownTime <= 0 || len(message.AlarmUrl) == 0 ||
		len(message.LogLevel) == 0 || len(message.LogPath) == 0 {
		logger.Error(config.BirdConf,"[judge] read message Abnormal parameter")
		b.response("[judge] read message Abnormal parameter",
			"update_domain",
			message.Id,
			message)
		return
	}
	if message.Code != "200" {
		// not response
		logger.Errorf(config.BirdConf,"[judge] bird conf response code not 200 data is %v",message.Data)
		return
	}
	switch  message.Method{
	case "get":
		b.getConfigs(message)
		return
	case "update_domain":
		b.updateDomain(message)
		return
	case "update_service":
		b.updateService(message)
		return
	case "delete_domain":
		b.deleteDomain(message)
		return
	case "delete_service":
		b.deleteService(message)
		return
	}
}

func (b *BirdConfigBasic)getConfigs(messages BirdConfMessage){
	addr := fmt.Sprintf("%v:%v",messages.BirdDomain,messages.Port)
	v,ok := cache.BirdCache.Get(addr)
	if ok {
		data,err := json.Marshal(v)
		if err != nil {
			logger.Errorf(config.BirdConf,"[getConfigs] get domain:%v config err %v",addr,err)
			goto  Loop
		}
		b.conn.Write(data)
		return
	}
Loop:
	logger.Errorf(config.BirdConf,"[getConfigs] get domain:%v config err",messages.BirdDomain)
	message ,_ := json.Marshal(&BirdConfMessage{
		Code:       "",
		Data:       "get domain config err",
		Method:     "get",
		BirdDomain: messages.BirdDomain,
		Port:messages.Port,
	})
	b.conn.Write(message)
	return
}

func (b *BirdConfigBasic)updateDomain(message BirdConfMessage){
	b.lock.Lock()
	b.buffer.Reset()
	b.fileBuffer.Reset()

	// check port
	_,ok := cache.ServiceInformerCache.Get(fmt.Sprintf(":%v",message.Port))
	if ok {
		b.response( fmt.Sprintf("[updateDomain] port:%v actually used ", message.Port),
			"update_domain",
			message.Id,
			message)
		return
	}

	r, err := os.OpenFile(hosts,os.O_RDWR,0644)
	defer func() {
		r.Close()
		b.lock.Unlock()
	}()
	if err != nil {
		b.response(fmt.Sprintf("[updateDomain] open hosts err:%v",err),
			"update_domain",
			message.Id,
			message)
		logger.Errorf(config.BirdConf,"[updateDomain] open file err :%v",err)
		return
	}

	b.buffer.WriteString("  ")
	b.buffer.WriteString(message.BirdDomain)
	b.buffer.WriteString("  ")
	br := bufio.NewReader(r)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			b.fileBuffer.WriteString("\n")
			b.fileBuffer.WriteString("0.0.0.0")
			b.fileBuffer.WriteString("  ")
			b.fileBuffer.WriteString(b.buffer.String())
			b.fileBuffer.WriteString("  ")
			r.WriteString(b.fileBuffer.String())
			break
		}
		if err != nil && err != io.EOF{
			b.response( fmt.Sprintf("[updateDomain] readline hosts err:%v",err),
				"update_domain",
				message.Id,
				message)
			logger.Errorf(config.BirdConf,"[updateDomain] readline file err :%v",err)
			return
		}
		if strings.Index(string(line),message.BirdDomain) != -1 {
			b.response( fmt.Sprintf("[updateDomain] update domain %v in hosts already here", message.BirdDomain),
				"update_domain",
				message.Id,
				message)
			logger.Warnf(config.BirdConf,"[updateDomain] update domain %v in hosts already here", message.BirdDomain)
			return
		}
		if strings.Index(string(line),"0.0.0.0") != -1 {
			b.fileBuffer.WriteString("  ")
			b.fileBuffer.WriteString(b.buffer.String())
			b.fileBuffer.WriteString("  ")
			r.Seek(io.SeekStart,io.SeekCurrent)
			r.WriteString(b.fileBuffer.String())
			break
		}
	}

	cache.BirdCache.Set(fmt.Sprintf("%v:%v",message.BirdDomain,message.Port),&message)

	server.Server(message.ReadTimeout,message.WriteTimeout,message.ShutdownTime,fmt.Sprintf(":%v",message.Port),
		message.AlarmUrl,message.IsDebug,message.LogLevel,message.LogPath)
}

func (b *BirdConfigBasic)updateService(message BirdConfMessage){

}

func (b *BirdConfigBasic)deleteDomain(message BirdConfMessage){
	addr := fmt.Sprintf("%v:%v",message.BirdDomain,message.Port)
	domain := fmt.Sprintf(":%v",message.Port)
	_,ok := cache.BirdCache.Get(addr)
	v,_ := cache.ServiceInformerCache.Get(domain)
	if ok {
		cache.BirdCache.Del(addr)
		cache.ServiceInformerCache.Del(domain)
		informer := v.(*bird_proxy.BirdProxy).Informer
		informer <- struct{}{}
		<- informer
		logger.Warnf(config.BirdConf,"[deleteDomain] delete domain:%v success",message.BirdDomain)
		b.response( fmt.Sprintf("[deleteDomain] delete domain:%v success",message.BirdDomain),
			"deleteDomain",
			message.Id,
			message)
		return
	}
	logger.Errorf(config.BirdConf,"[deleteDomain] delete domain %v,port:%v service is not run", message.BirdDomain,message.Port)
	b.response( fmt.Sprintf("[deleteDomain] delete domain %v service is not runi", message.BirdDomain),
		"deleteDomain",
		message.Id,
		message)
	return
}

func (b *BirdConfigBasic)deleteService(message BirdConfMessage){

}

func (b *BirdConfigBasic)response(str,method,id string,message BirdConfMessage){
	data,_ := json.Marshal(&BirdConfMessage{
		Id:			id,
		Code:       "200",
		Data:       str,
		Method:     method,
		AlarmUrl:   message.AlarmUrl,
		BirdDomain: message.BirdDomain,
		Port:       message.Port,
		Route:    message.Route,
		ReverseProxyService:message.ReverseProxyService,
	})
	b.conn.Write(data)
	return
}