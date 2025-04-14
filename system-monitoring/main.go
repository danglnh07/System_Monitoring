package main

import "system-monitoring/server"

func main() {
	//Create server and start
	server := server.NewServer()
	server.Start()
}
