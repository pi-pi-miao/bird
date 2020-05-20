package bird_proxy

import (
	"fmt"
	"net/http"
)

func (p *BirdProxy)ServeHTTP(w http.ResponseWriter, r *http.Request){
	// judge request domain
	fmt.Println("host",r.Host)

}


