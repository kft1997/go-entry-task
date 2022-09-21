package server

import (
	"go-entry-task/rpcsdk/rpc"
	"log"
	"net"
)

var server *Server

//func init() {
//	//initDB()
//	//initRedis()
//	server = new(Server)
//	rpc.Register("login", server.LoginHandler)
//	rpc.Register("query", server.QueryHandler)
//	rpc.Register("updateNick", server.EditNickHandler)
//	//rpc.Register("logout",HandlerLogout)
//	rpc.Register("uploadpic", server.EditPicHandler)
//}

func Init() {
	server = new(Server)
	server.register()
}

//服务器运行
func Run(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("server start")
	rpc.Serve(lis)
}
