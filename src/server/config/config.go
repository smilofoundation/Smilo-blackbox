package config

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	log    *logrus.Entry
	config Config

	GenerateKeysStr = "generate-keys"
	GenerateKeys    = ""

	ConfigFileStr = "configfile"
	ConfigFile    = "configfile"
	DBFileStr     = "dbfile"
	DBFile        = ""

	PortStr = "port"
	Port    = ""

	WorkDirStr = "workdir"
	WorkDir    = ""

	SocketStr = "socket"
	Socket    = ""

	OtherNodesStr = "othernodes"
	OtherNodes    = ""

	PublicKeysStr = "publickeys"
	PublicKeys    = ""

	PrivateKeysStr = "privatekeys"
	PrivateKeys    = ""

	AlwaysSendToStr = "alwayssendto"
	AlwaysSendTo    = ""

	VerbosityStr = "verbosity"
	Verbosity    = 0

	HostNameStr = "hostname"
	HostName    = ""

	IsTLSStr = "tls"
	IsTLS    = ""

	ServerCertStr = "server_cert"
	ServerCert    = ""

	ServerKeyStr = "server_key"
	ServerKey    = ""

	MaxPeersNetworkStr = "maxpeersnetwork"
	MaxPeersNetwork    = ""

	P2PDestinationStr = "p2p_dest"
	P2PDestination    = ""

	P2PPortStr = "p2p_port"
	P2PPort    = ""
)

func initLog() {
	log = logrus.WithFields(logrus.Fields{
		"app":     "blackbox",
		"package": "config",
	})
}

func Init() {
	initLog()

	flag.String(GenerateKeysStr, "", "Generate a new keypair")

	flag.String(ConfigFileStr, "blackbox.conf", "Config file name")

	flag.String(DBFileStr, "blackbox.db", "DB file name")

	flag.String(PortStr, "9000", "Local port to the Public API")

	flag.String(WorkDirStr, "../../", "")

	flag.String(SocketStr, "blackbox.ipc", "IPC socket to the Private API")

	flag.String(OtherNodesStr, "", "\"Boot nodes\" to connect")

	flag.String(PublicKeysStr, "", "Public keys")

	flag.String(PrivateKeysStr, "", "Private keys")

	flag.String(AlwaysSendToStr, "", "List of public keys for nodes to send all transactions too")

	flag.Int(VerbosityStr, 1, "Verbosity level of logs")

	flag.String(HostNameStr, "http://localhost", "HostName for public API")

	flag.String(IsTLSStr, "", "")

	flag.String(ServerKeyStr, "", "")

	flag.String(ServerCertStr, "", "")

	flag.String(MaxPeersNetworkStr, "", "")

	flag.String(P2PDestinationStr, "", "")

	flag.String(P2PPortStr, "30300", "")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)
	setCommandList()
}

func ConfigLoad(path string) error {
	//Init()
	viper.SetConfigType("hcl")
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()

	//Init()

	GenerateKeys = GetString(GenerateKeysStr)
	ConfigFile = GetString(ConfigFileStr)
	DBFile = GetString(DBFileStr)
	Port = GetString(PortStr)
	WorkDir = GetString(WorkDirStr)
	Socket = GetString(SocketStr)
	OtherNodes = GetString(OtherNodesStr)
	PublicKeys = GetString(PublicKeysStr)
	PrivateKeys = GetString(PrivateKeysStr)
	AlwaysSendTo = GetString(AlwaysSendToStr)
	Verbosity = GetInt(VerbosityStr)
	HostName = GetString(HostNameStr)

	IsTLS = GetString(IsTLS)
	ServerKey = GetString(ServerKeyStr)
	ServerCert = GetString(ServerCertStr)
	MaxPeersNetwork = GetString(MaxPeersNetworkStr)
	P2PDestination = GetString(P2PDestinationStr)
	P2PPort = GetString(P2PPortStr)

	return err
}

func RefreshAllKeys() {

}

func AllSettings() map[string]interface{} {
	return viper.AllSettings()
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func setCommandList() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)

	args := os.Args
	if len(args) == 1 {
		//log.Fatalln("No args defined")
	}

	for _, arg := range args[1:] {
		if strings.Contains(arg, ".conf") {
			err := ConfigLoad(arg)
			if err != nil {
				log.Fatalln(err)
			}
			break
		}
	}
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	//app.Flags = []cli.Flag{GenerateKeys, ConfigFile, DBFile, Port, Socket, OtherNodes, PublicKeys, PrivateKeys, Storage, HostName, WorkDir, IsTLS, ServCert, ServKey}

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
