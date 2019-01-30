package main

import (
	"Smilo-blackbox/src/crypt"
	"Smilo-blackbox/src/server"
	"Smilo-blackbox/src/server/config"
	"os"

	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"Smilo-blackbox/src/server/syncpeer"
)

var (
	log *logrus.Entry
)

func initLog() {
	log = logrus.WithFields(logrus.Fields{
		"app":     "blackbox",
		"package": "main",
	})
}

func init() {
	initLog()
}

func main() {

	defer handlePanic()

	app := cli.NewApp()
	config.LoadConfig(config.ConfigFile.Value)
	config.Init(app)

	app.Name = "blackbox"
	app.Usage = "safe storage and exchange service for private transactions"
	app.Action = func(c *cli.Context) error {
		generateKeys := c.String("generate-keys")
		if generateKeys != "" {
			crypt.GenerateKeys(generateKeys)
		} else {
			server.StartServer()
			server.InitP2p()
			syncpeer.StartSync()
		}
		return nil
	}

	dir, _ := os.Getwd() // gives us the source path if we haven't installed.
	config.WorkDir.Value = dir
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	select {}

}

func handlePanic() {
	if r := recover(); r != nil {
		log.WithError(fmt.Errorf("%+v", r)).Error(fmt.Sprintf("Application BlackBox panic"))
	}
	time.Sleep(time.Second * 5)
}
