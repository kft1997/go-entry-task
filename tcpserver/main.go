package main

import "go-entry-task/tcpserver/server"

func main() {
	server.Init()
	server.Run("127.0.0.1:8080")
}
