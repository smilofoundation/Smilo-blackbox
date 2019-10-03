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

package utils

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	//BlackBoxVersion holds bb version
	BlackBoxVersion = "Smilo Black Box 0.1.0"
	//UpcheckMessage http up check msg
	UpcheckMessage = "I'm up!"

	//HeaderFrom header default
	HeaderFrom = "bb0x-from"
	//HeaderTo header default
	HeaderTo = "bb0x-to"
	//HeaderKey header default
	HeaderKey = "bb0x-key"
)

//BuildFilename will build a filename with correct path based on pwd
func BuildFilename(filename string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var workDir string
	var newDBFile string
	if flag.Lookup("test.v") != nil {
		isServer := strings.HasSuffix(currentDir, "/server")
		isData := strings.HasSuffix(currentDir, "/data")
		isConfig := strings.HasSuffix(currentDir, "/config")
		isRoot := strings.HasSuffix(currentDir, "/Smilo-blackbox")
		if isServer {
			workDir = "../../"
		} else if isData {
			workDir = "../../"
		} else if isRoot {
			workDir = ""
		} else if isConfig {
			workDir = "../../../"
		}
	} else {
		workDir = ""
	}

	newDBFile = path.Join(workDir, filename)
	return newDBFile
}

func GetType(data interface{}) reflect.Type {
	t := reflect.TypeOf(data).Elem()
	return t
}

func GetMetadata(data interface{}) (string, string) {
	t := GetType(data)
	keyField := ""
	for key := 0; key < t.NumField(); key++ {
		field := t.Field(key)
		if field.Tag.Get("key") == "true" {
			keyField = field.Name
		}
	}
	return t.Name(), keyField
}

func GetField(data interface{}, field string) interface{} {
	r := reflect.ValueOf(data).Elem()
	s := r.FieldByName(field).Interface()
	return s
}

func ReadAllFile(file string, log *logrus.Entry) ([]byte, error) {
	plainFile, err := os.Open(file)
	defer func() {
		err := plainFile.Close()
		if err != nil {
			log.WithError(err).Error("Could not plainFile.Close")
		}
	}()
	if err != nil {
		return nil, err
	}
	byteValue, err := ioutil.ReadAll(plainFile)
	return byteValue, err
}
