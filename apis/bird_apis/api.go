package bird_apis

import (
	"bird/pkg/bird_controller"
	"net/http"
)

func BirdApis(handler http.Handler){
	bird := handler.(*http.ServeMux)
	mux := http.NewServeMux()
	mux.Handle("/", &myHandler{})
	bird.HandleFunc("/",bird_controller.Example)
	bird.HandleFunc("/login",bird_controller.Login)
}

