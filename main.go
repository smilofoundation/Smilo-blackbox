package main

import (
	"Smilo-blackbox/src/server"
)

func main() {
	port := "9000"
	server.StartServer(port,"")
}