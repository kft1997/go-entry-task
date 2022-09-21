package main

import "go-entry-task/httpserver/server"

func main() {
	server.Init()
	server.Run(":8000")
}
