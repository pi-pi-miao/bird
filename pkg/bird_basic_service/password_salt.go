package bird_basic_service

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"sync"
)

var (
	buffer = &bytes.Buffer{}
	lock   = &sync.Mutex{}
	PasswordSlat = ""
)

func CreateToken(account,password string)string{
	lock.Lock()
	defer lock.Unlock()
	buffer.Reset()
	buffer.WriteString(account)
	buffer.WriteString(password)
	buffer.WriteString(PasswordSlat)
	return fmt.Sprintf("%v",sha256.Sum256(buffer.Bytes()))
}

func CompareToken(account,password,token string)bool{
	lock.Lock()
	defer lock.Unlock()
	buffer.Reset()
	buffer.WriteString(account)
	buffer.WriteString(password)
	buffer.WriteString(PasswordSlat)
	return token == fmt.Sprintf("%v",sha256.Sum256(buffer.Bytes()))
}
