package config

import (
	flag "github.com/spf13/pflag"
	"os"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"Smilo-blackbox/src/crypt"
	"strconv"
)

var config Config



func Init() {
	InitFlags()
    flag.Parse()
	LoadConfig(flag.Lookup("configfile").Value.String())
	mergeConfigValues()
}

func mergeConfigValues () {
	setValueOnNotDefault("port", string(config.Server.Port))
	setValueOnNotDefault("socket", config.UnixSocket)
	setValueOnNotDefault("hostname", string(config.HostName))
}

func setValueOnNotDefault(flagName string, flagValue string) {
	fg := flag.Lookup(flagName)
	if fg.Value.String() == fg.DefValue && flagValue != "" {
		fg.Value.Set(flagValue)
	}
}

func InitFlags() {

	//flag.String("generate-keys", "", "Generate a new keypair")
	flag.String("configfile", "./blackBox.conf", "Config file name")
	flag.Int("port", 9000, "Local port to the Public API")
	flag.String("socket", "./blackBox.ipc", "IPC socket to the Private API")
	flag.String("othernodes", "", "\"Boot nodes\" to connect")
	flag.String("publickeys", "", "Public keys")
	flag.String("privatekeys", "", "Private keys")
	flag.String("storage", "./blackBox.db", "Database file name")

	//flag.Bool("tls", false, "Use TLS to secure HTTP communications")
	//flag.String("tlsservercert", "", "The server certificate to be used")
	//flag.String("tlsserverkey", "", "The server private key")
	flag.String("hostname" , "http://localhost", "HostName for public API")


}


func LoadConfig(configPath string) error {
	byteValue, err := readAllFile(configPath)
	if err != nil {
		return err
	}

	json.Unmarshal(byteValue, &config)
    parseConfigValues()
	return nil;
}

func parseConfigValues() {
	for _, keyPair := range config.Keys.KeyData {
		primaryKey, err := ReadPrimaryKey(keyPair.PrivateKeyFile)
		publicKey, err2 := ReadPublicKey(keyPair.PublicKeyFile)
		if err != nil || err2 != nil { continue }
		crypt.PutKeyPair(crypt.KeyPair{ PrimaryKey:primaryKey, PublicKey:publicKey})
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

func GetSocketFile() string {
	return flag.Lookup("socket").Value.String()
}

func GetHostName() string {
	return flag.Lookup("hostname").Value.String()
}

func GetServerPort() int {
	ret, _ := strconv.Atoi(flag.Lookup("port").Value.String())
	return ret
}

func GetDatabaseName() string {
	return flag.Lookup("storage").Value.String()
}