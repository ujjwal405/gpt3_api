package router

import (
	"log"
	"net/http"

	"github.com/Ujjwal405/gpt3/client/handler"
)

func RunServer(h *handler.Userhandler, ch chan error, add string) {

	http.HandleFunc("/answer", h.Getanswer)
	http.HandleFunc("/search", h.Getsearch)
	log.Printf("Json_server starting at port : %s", add)
	if err := http.ListenAndServe(add, nil); err != nil {
		ch <- err
	}

}
