package main

import (
	"Smilo-blackbox/src/server"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry
)

func initLog() {
	log = logrus.WithFields(logrus.Fields{
		"app": "blackbox",
	})
}

func init() {
	initLog()
}

func main() {
	port := "9000"
	server.StartServer(port, "")
}
