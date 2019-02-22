// Copyright 2019 The Smilo-blackbox Authors
// This file is part of the Smilo-blackbox library.
//
// The Smilo-blackbox library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Smilo-blackbox library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Smilo-blackbox library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"Smilo-blackbox/src/crypt"
	"Smilo-blackbox/src/server"
	"Smilo-blackbox/src/server/config"
	"os"

	"fmt"
	"time"

	"Smilo-blackbox/src/server/syncpeer"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"runtime/pprof"
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
	config.Init(app)
	app.Name = "blackbox"
	app.Usage = "safe storage and exchange service for private transactions"
	app.Action = func(c *cli.Context) error {
		generateKeys := c.String("generate-keys")
		configFile := c.String("configfile")
		cpuProfilingFile := c.String("cpuprofile")
		p2pEnabled := c.Bool("p2p")
		if generateKeys != "" {
			crypt.GenerateKeys(generateKeys)
			os.Exit(0)
		} else {
			if cpuProfilingFile != "" {
				f, err := os.Create(cpuProfilingFile)
				if err != nil {
					log.Fatal(err)
				}
				pprof.StartCPUProfile(f)
				defer pprof.StopCPUProfile()
			}
			config.LoadConfig(configFile)
			server.StartServer()
			if p2pEnabled {
				server.InitP2p()
			}
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
