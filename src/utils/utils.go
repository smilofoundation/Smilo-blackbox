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
	"os"
	"path"
	"strings"
)

const (
	BlackBoxVersion = "Smilo Black Box 0.1.0"
	UpcheckMessage  = "I'm up!"

	HeaderFrom = "bb0x-from"
	HeaderTo   = "bb0x-to"
	HeaderKey  = "bb0x-key"
)

func BuildFilename(filename string) string {
	currentDir, _ := os.Getwd()
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
	//if !strings.HasPrefix(currentDir, "/"){
	//	newDBFile = path.Join(currentDir, workDir)
	//}
	newDBFile = path.Join(workDir, filename)
	return newDBFile
}
