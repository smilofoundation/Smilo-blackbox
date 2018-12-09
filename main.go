package main

import (
	"Smilo-blackbox/src/server"

	"github.com/sirupsen/logrus"
	"Smilo-blackbox/src/server/config"
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
	config.Init()
	server.StartServer(config.GetServerPort(), config.GetSocketFile())
}
