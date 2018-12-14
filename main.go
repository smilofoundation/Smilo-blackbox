package main

import (
	"Smilo-blackbox/src/server"

	"Smilo-blackbox/src/server/config"
	"fmt"
	"os"
	"time"

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
	config.Init()
	dir, _ := os.Getwd() // gives us the source path if we haven't installed.
	config.WorkDir.Value = dir
}

func main() {

	defer handlePanic()

	server.StartServer()
}

func handlePanic() {
	if r := recover(); r != nil {
		log.WithError(fmt.Errorf("%+v", r)).Error(fmt.Sprintf("Application BlackBox panic"))
	}
	time.Sleep(time.Second * 5)
}
