package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"

	"Smilo-blackbox/src/crypt"
)

var (
	log    *logrus.Entry
	config Config

	//flag.String("generate-keys", "", "Generate a new keypair")
	ConfigFile  = cli.StringFlag{Name: "configfile", Value: "blackbox.conf", Usage: "Config file name"}
	DBFile      = cli.StringFlag{Name: "dbfile", Value: "blackbox.db", Usage: "DB file name"}
	Port        = cli.StringFlag{Name: "port", Value: "9000", Usage: "Local port to the Public API"}
	Socket      = cli.StringFlag{Name: "socket", Value: "blackbox.ipc", Usage: "IPC socket to the Private API"}
	OtherNodes  = cli.StringFlag{Name: "othernodes", Value: "", Usage: "\"Boot nodes\" to connect"}
	PublicKeys  = cli.StringFlag{Name: "publickeys", Value: "", Usage: "Public keys"}
	PrivateKeys = cli.StringFlag{Name: "privatekeys", Value: "", Usage: "Private keys"}
	Storage     = cli.StringFlag{Name: "storage", Value: "blackbox.db", Usage: "Database file name"}

	HostName = cli.StringFlag{Name: "hostname", Value: "http://localhost", Usage: "HostName for public API"}

	WorkDir  = cli.StringFlag{Name: "workdir", Value: "../../", Usage: ""}
	IsTLS    = cli.StringFlag{Name: "tls", Value: "", Usage: ""}
	ServCert = cli.StringFlag{Name: "serv_cert", Value: "", Usage: ""}
	ServKey  = cli.StringFlag{Name: "serv_key", Value: "", Usage: ""}
)

func initLog() {
	log = logrus.WithFields(logrus.Fields{
		"package": "config",
	})
}

func Init() {
	initLog()
	pflag.Parse()
	LoadConfig(ConfigFile.Value)
	//mergeConfigValues()

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine) // Binding the flags to test the initial configuration
}

//
//func mergeConfigValues() {
//	setValueOnNotDefault("port", string(config.Server.Port))
//	setValueOnNotDefault("socket", config.UnixSocket)
//	setValueOnNotDefault("hostname", string(config.HostName))
//}

//func setValueOnNotDefault(flagName string, flagValue string) {
//	fg := pflag.Lookup(flagName)
//	if fg.Value == fg.DefValue && flagValue != "" {
//		fg.Value.Set(flagValue)
//	}
//}

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
		crypt.PutKeyPair(crypt.KeyPair{PrimaryKey: primaryKey, PublicKey: publicKey})
	}
	/*
		for _, peerdata := range config.Peers {
			sync.PeerAdd(peerdata.URL)
		}
	*/
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
