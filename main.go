package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"Smilo-blackbox/src/server/config"
	"strings"
	"github.com/twystd/tweetnacl-go"
	"encoding/base64"
	"io/ioutil"
	"Smilo-blackbox/src/server"
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

	defer handlePanic()

	app := cli.NewApp()
	config.Init(app)
	config.LoadConfig(config.ConfigFile.Value)

	app.Name = "blackbox"
	app.Usage = "safe storage and exchange service for private transactions"
	app.Action = func(c *cli.Context) error {
		log.Info("xxxxx")
		if c.String("generate-keys") != "" {
			files := strings.Split(c.String("generate-keys"), ",")
			for i := range files {
				keyPair, _ := tweetnacl.CryptoBoxKeyPair()
				writePrivateKeyFile(base64.StdEncoding.EncodeToString(keyPair.SecretKey),files[i]+".key")
				writePublicKeyFile(base64.StdEncoding.EncodeToString(keyPair.PublicKey),files[i]+".pub")
			}
		} else {
			server.StartServer()
		}
		return nil
	}

	dir, _ := os.Getwd() // gives us the source path if we haven't installed.
	config.WorkDir.Value = dir
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func writePrivateKeyFile(key string, filename string) error {
	filedata := "{\"type\" : \"unlocked\",\"data\" : {\"bytes\" : \""+key+"\"}}"
	return ioutil.WriteFile(filename, []byte(filedata), os.ModePerm)
}

func writePublicKeyFile(key string, filename string) error {
	return ioutil.WriteFile(filename, []byte(key), os.ModePerm)
}
func handlePanic() {
	if r := recover(); r != nil {
		log.WithError(fmt.Errorf("%+v", r)).Error(fmt.Sprintf("Application BlackBox panic"))
	}
	time.Sleep(time.Second * 5)
}
