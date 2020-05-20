package bird_controller

import "net/http"

func Example(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("hello"))
}
