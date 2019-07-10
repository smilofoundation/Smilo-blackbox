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

	"strings"

	"Smilo-blackbox/src/data"
	"Smilo-blackbox/src/server/syncpeer"
	"Smilo-blackbox/src/utils"
)

var (
	log    *logrus.Entry
	config Config

	//GenerateKeys (cli) uses it for key pair
	GenerateKeys = cli.StringFlag{Name: "generate-keys", Value: "", Usage: "Generate a new keypair"}
	//ConfigFile (cli) uses it for config file name
	ConfigFile = cli.StringFlag{Name: "configfile", Value: "blackbox.conf", Usage: "Config file name"}
	//DBEngine (cli) uses it for db engine
	DBEngine = cli.StringFlag{Name: "dbengine", Value: "boltdb", Usage: "DB engine name"}
	//DBFile (cli) uses it for db file name
	DBFile = cli.StringFlag{Name: "dbfile", Value: "blackbox.db", Usage: "DB file name"}
	//PeersDBFile (cli) uses it for peer db file
	PeersDBFile = cli.StringFlag{Name: "peersdbfile", Value: "blackbox-peers.db", Usage: "Peers DB file name"}
	//Port (cli) uses it for local api public port
	Port = cli.StringFlag{Name: "port", Value: "9000", Usage: "Local port to the Public API"}
	//Socket (cli) uses it for socket
	Socket = cli.StringFlag{Name: "socket", Value: "blackbox.ipc", Usage: "IPC socket to the Private API"}
	//OtherNodes (cli) uses it for other nodes
	OtherNodes = cli.StringFlag{Name: "othernodes", Value: "", Usage: "\"Boot nodes\" to connect"}
	//PublicKeys (cli) uses it for  pub
	PublicKeys = cli.StringFlag{Name: "publickeys", Value: "", Usage: "Public keys"}
	//PrivateKeys (cli) uses it for pk
	PrivateKeys = cli.StringFlag{Name: "privatekeys", Value: "", Usage: "Private keys"}
	//Storage (cli) uses it for  db name
	Storage = cli.StringFlag{Name: "storage", Value: "blackbox.db", Usage: "Database file name"}
	//HostName (cli) uses it for hostname
	HostName = cli.StringFlag{Name: "hostname", Value: "http://localhost", Usage: "HostName for public API"}

	//WorkDir (cli) uses it for work dir
	WorkDir = cli.StringFlag{Name: "workdir", Value: "../../", Usage: ""}
	//IsTLS (cli) uses it for enable/disable https
	IsTLS = cli.BoolFlag{Name: "tls", Usage: "Enable HTTPs communication"}
	//ServCert (cli) uses it for cert
	ServCert = cli.StringFlag{Name: "serv_cert", Value: "", Usage: ""}
	//ServKey (cli) uses it for key
	ServKey = cli.StringFlag{Name: "serv_key", Value: "", Usage: ""}

	//MaxPeersNetwork = cli.StringFlag{Name: "maxpeersnetwork", Value: "", Usage: ""}

	//P2PDestination (cli) uses it for p2p dest
	P2PDestination = cli.StringFlag{Name: "p2p_dest", Value: "", Usage: ""}
	//P2PPort (cli) uses it for p2p port
	P2PPort = cli.StringFlag{Name: "p2p_port", Value: "", Usage: ""}
	//CPUProfiling (cli) uses it for CPU profiling data filename
	CPUProfiling = cli.StringFlag{Name: "cpuprofile", Value: "", Usage: "CPU profiling data filename"}
	//P2PEnabled (cli) uses it for enable / disable p2p
	P2PEnabled = cli.BoolFlag{Name: "p2p", Usage: "Enable p2p communication"}
	//RootCert  (cli) uses it for certs
	RootCert = cli.StringFlag{Name: "root_cert", Value: "", Usage: ""}
)

func initLog() {
	log = logrus.WithFields(logrus.Fields{
		"app":     "blackbox",
		"package": "config",
	})
}

//Init will init cli and logs
func Init(app *cli.App) {
	initLog()
	setCommandList(app)
}

func setCommandList(app *cli.App) {
	app.Flags = []cli.Flag{GenerateKeys, ConfigFile, DBEngine, DBFile, PeersDBFile, Port, Socket, OtherNodes, PublicKeys, PrivateKeys, Storage, HostName, WorkDir, IsTLS, ServCert, ServKey, RootCert, CPUProfiling, P2PEnabled}
}

//LoadConfig will load cfg
func LoadConfig(configPath string) error {
	byteValue, err := readAllFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		log.WithError(err).Error("Could not json.Unmarshal config")
		return err
	}
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
	if config.HostName != "" {
		HostName.Value = config.HostName
	}
	if IsTLS.Destination == nil {
		var local bool
		IsTLS.Destination = &local
	}
	if config.Server.TLSCert != "" {
		ServCert.Value = config.Server.TLSCert
	}
	if config.Server.TLSKey != "" {
		ServKey.Value = config.Server.TLSKey
	}
	if ServKey.Value != "" || ServCert.Value != "" {
		*IsTLS.Destination = true
	}
	RootCertArray := config.RootCA
	if RootCertArray == nil {
		RootCertArray = []string{}
	}
	if RootCert.Value != "" {
		RootCertArray = append(RootCertArray, strings.Split(RootCert.Value, ",")...)
	}
	for _, cert := range RootCertArray {

		certData, err := ioutil.ReadFile(utils.BuildFilename(cert))
		if err != nil {
			log.Errorf("Failed to read %q, %v", cert, err)
		}
		if !syncpeer.AppendCertificate(certData) {
			log.Errorf("Failed to append %q to RootCAs", cert)
		}

	}
	if config.DBEngine != "" {
		DBEngine.Value = config.DBEngine
	}
	if config.DBFile != "" {
		DBFile.Value = config.DBFile
	}
	if config.PeersDBFile != "" {
		PeersDBFile.Value = config.PeersDBFile
	}
	data.SetFilename(utils.BuildFilename(DBFile.Value))
	data.SetEngine(DBEngine.Value)
	syncpeer.SetHostURL(HostName.Value + ":" + Port.Value)
	for _, peerdata := range config.Peers {
		syncpeer.PeerAdd(peerdata.URL)
	}
}

//ReadPrimaryKey will read pk
func ReadPrimaryKey(pkFile string) ([]byte, error) {
	byteValue, err := readAllFile(pkFile)
	if err != nil {
		return nil, err
	}

	var privateKey PrivateKey
	err = json.Unmarshal(byteValue, &privateKey)
	if err != nil {
		log.WithError(err).Error("Could not json.Unmarshal privateKey")
		return nil, err
	}

	var decodedPrivateKey = make([]byte, 33)

	_, err = base64.StdEncoding.Decode(decodedPrivateKey, []byte(privateKey.Data.Bytes))

	return decodedPrivateKey[0:32], err
}

//ReadPublicKey will read pub
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
	defer func() {
		err := plainFile.Close()
		log.WithError(err).Error("Could not plainFile.Close")
	}()
	if err != nil {
		return nil, err
	}
	byteValue, _ := ioutil.ReadAll(plainFile)
	return byteValue, nil
}
