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
	"Smilo-blackbox/src/data"
	"Smilo-blackbox/src/server"
	"Smilo-blackbox/src/server/config"
	"Smilo-blackbox/src/utils"
	"os"

	"fmt"
	"time"

	"Smilo-blackbox/src/server/syncpeer"

	"runtime/pprof"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
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
	app.Version = utils.BlackBoxVersion

	app.Action = func(c *cli.Context) error {
		generateKeys := c.String("generate-keys")
		migrateDatabase := c.Bool("migrate-database")
		configFile := c.String("configfile")
		cpuProfilingFile := c.String("cpuprofile")
		if generateKeys != "" {
			err := crypt.GenerateKeys(generateKeys)
			if err != nil {
				os.Exit(-1)
			} else {
				os.Exit(0)
			}
		}
		if migrateDatabase {
			var dbengine string
			var dbfile string
			if configFile != "blackbox.conf" {
				LoadConfig(configFile)
				dbengine = config.DBEngine.Value
				dbfile = config.DBFile.Value
			} else {
				dbengine = c.String("dbengine")
				dbfile = c.String("dbfile")
			}
			err := data.Migrate(dbengine, utils.BuildFilename(dbfile), c.String("dbengine-dest"), utils.BuildFilename(c.String("dbfile-dest")))
			if err != nil {
				os.Exit(-1)
			} else {
				os.Exit(0)
			}
		}
		if cpuProfilingFile != "" {
			f, err := os.Create(cpuProfilingFile)
			if err != nil {
				log.Fatal(err)
			}
			err = pprof.StartCPUProfile(f)
			if err != nil {
				log.Fatal(err)
			}
			defer pprof.StopCPUProfile()
		}
		LoadConfig(configFile)
		server.StartServer()
		syncpeer.StartSync()
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

func LoadConfig(configFile string) {
	err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func handlePanic() {
	if r := recover(); r != nil {
		log.WithError(fmt.Errorf("%+v", r)).Error(fmt.Sprintf("Application BlackBox panic"))
	}
	time.Sleep(time.Second * 5)
}
