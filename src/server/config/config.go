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

package config

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"

	"Smilo-blackbox/src/crypt"

	"strconv"

	"Smilo-blackbox/src/data"
	"Smilo-blackbox/src/server/syncpeer"
	"Smilo-blackbox/src/utils"
)

var (
	log    *logrus.Entry
	config Config

	GenerateKeys = cli.StringFlag{Name: "generate-keys", Value: "", Usage: "Generate a new keypair"}
	ConfigFile   = cli.StringFlag{Name: "configfile", Value: "blackbox.conf", Usage: "Config file name"}
	DBFile       = cli.StringFlag{Name: "dbfile", Value: "blackbox.db", Usage: "DB file name"}
	PeersDBFile  = cli.StringFlag{Name: "peersdbfile", Value: "blackbox-peers.db", Usage: "Peers DB file name"}
	Port         = cli.StringFlag{Name: "port", Value: "9000", Usage: "Local port to the Public API"}
	Socket       = cli.StringFlag{Name: "socket", Value: "blackbox.ipc", Usage: "IPC socket to the Private API"}
	OtherNodes   = cli.StringFlag{Name: "othernodes", Value: "", Usage: "\"Boot nodes\" to connect"}
	PublicKeys   = cli.StringFlag{Name: "publickeys", Value: "", Usage: "Public keys"}
	PrivateKeys  = cli.StringFlag{Name: "privatekeys", Value: "", Usage: "Private keys"}
	Storage      = cli.StringFlag{Name: "storage", Value: "blackbox.db", Usage: "Database file name"}

	HostName = cli.StringFlag{Name: "hostname", Value: "http://localhost", Usage: "HostName for public API"}

	WorkDir  = cli.StringFlag{Name: "workdir", Value: "../../", Usage: ""}
	IsTLS    = cli.StringFlag{Name: "tls", Value: "", Usage: ""}
	ServCert = cli.StringFlag{Name: "serv_cert", Value: "", Usage: ""}
	ServKey  = cli.StringFlag{Name: "serv_key", Value: "", Usage: ""}

	MaxPeersNetwork = cli.StringFlag{Name: "maxpeersnetwork", Value: "", Usage: ""}

	P2PDestination = cli.StringFlag{Name: "p2p_dest", Value: "", Usage: ""}

	P2PPort = cli.StringFlag{Name: "p2p_port", Value: "", Usage: ""}
)

func initLog() {
	log = logrus.WithFields(logrus.Fields{
		"app":     "blackbox",
		"package": "config",
	})
}

func Init(app *cli.App) {
	initLog()
	setCommandList(app)
}

func setCommandList(app *cli.App) {
	app.Flags = []cli.Flag{GenerateKeys, ConfigFile, DBFile, PeersDBFile, Port, Socket, OtherNodes, PublicKeys, PrivateKeys, Storage, HostName, WorkDir, IsTLS, ServCert, ServKey}

}

func LoadConfig(configPath string) error {
	byteValue, err := readAllFile(configPath)
	if err != nil {
		return err
	}

	json.Unmarshal(byteValue, &config)
	parseConfigValues()
	return nil
}

func parseConfigValues() {
	for _, keyPair := range config.Keys.KeyData {
		primaryKey, err := ReadPrimaryKey(keyPair.PrivateKeyFile)
		publicKey, err2 := ReadPublicKey(keyPair.PublicKeyFile)
		if err != nil || err2 != nil {
			continue
		}
		crypt.PutKeyPair(crypt.KeyPair{PrivateKey: primaryKey, PublicKey: publicKey})
	}
	Port.Value = strconv.FormatInt(int64(config.Server.Port), 10)
	if config.UnixSocket != "" {
		Socket.Value = config.UnixSocket
	}
	if config.DBFile != "" {
		DBFile.Value = config.DBFile
	}
	if config.PeersDBFile != "" {
		PeersDBFile.Value = config.PeersDBFile
	}
	data.SetFilename(utils.BuildFilename(DBFile.Value))
	syncpeer.SetHostUrl(HostName.Value+":"+Port.Value)
	for _, peerdata := range config.Peers {
		syncpeer.PeerAdd(peerdata.URL)
	}
}

func ReadPrimaryKey(pkFile string) ([]byte, error) {
	byteValue, err := readAllFile(pkFile)
	if err != nil {
		return nil, err
	}

	var privateKey PrivateKey
	json.Unmarshal(byteValue, &privateKey)

	var decodedPrivateKey = make([]byte, 33)

	_, err = base64.StdEncoding.Decode(decodedPrivateKey, []byte(privateKey.Data.Bytes))

	return decodedPrivateKey[0:32], err
}

func ReadPublicKey(pubFile string) ([]byte, error) {
	byteValue, err := readAllFile(pubFile)
	if err != nil {
		return nil, err
	}
	var publicKey = make([]byte, 33)

	_, err = base64.StdEncoding.Decode(publicKey, byteValue)

	return publicKey[0:32], err
}

func readAllFile(file string) ([]byte, error) {
	plainFile, err := os.Open(file)
	defer plainFile.Close()
	if err != nil {
		return nil, err
	}
	byteValue, _ := ioutil.ReadAll(plainFile)
	return byteValue, nil
}
