package server

import (
	"go-entry-task/rpcsdk/rpc"
	"log"
	"net/http"
)

type userInfo struct {
	Username string
	Nickname string
	Url      string
}

type response struct {
	Code int
	Msg  interface{}
}

var server *Server

var errResponse = response{http.StatusInternalServerError, nil}

func Init() {
	if err := rpc.Init("127.0.0.1:8080"); err != nil {
		log.Fatalln(err)
	}
	server = new(Server)
	http.HandleFunc("/", server.ViewHandler)
	http.HandleFunc("/login", server.LoginHandler)
	http.HandleFunc("/query", server.QueryHandler)
	http.HandleFunc("/editnick", server.EditNickHandler)
	http.HandleFunc("/uploadpic", server.EditPicHandler)
	http.Handle("/picfile/", http.FileServer(http.Dir("")))
}

func Run(addr string) {
	http.ListenAndServe(addr, nil)
}
