package config

import (
	"flag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
)

// InitFlags initializes all supported command line flags.
func InitFlags() {
	flag.String("generate-keys", "", "Generate a new keypair")
	flag.Int("port", 9000, "Local port to the Public API")
	flag.String("socket", "blackBox.ipc", "IPC socket to the Private API")
	flag.String("othernodes", "", "\"Boot nodes\" to connect")
	flag.String("publickeys", "", "Public keys")
	flag.String("privatekeys", "", "Private keys")
	flag.String("storage", "./blackBox.db", "Database file name")
	flag.String("configfile", "./blackBox.conf", "")

	flag.Bool("tls", false, "Use TLS to secure HTTP communications")
	flag.String("tlsservercert", "", "The server certificate to be used")
	flag.String("tlsserverkey", "", "The server private key")
	flag.String("networkinterface" , "localhost", "The network interface to bind the server to")


	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine) // Binding the flags to test the initial configuration
}

// ParseCommandLine parses all provided command line arguments.
func ParseCommandLine() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

// LoadConfig loads all configuration settings in the provided configPath location.
func LoadConfig(configPath string) error {
	viper.SetConfigType("json")
	viper.SetConfigFile(configPath)
	return viper.ReadInConfig()
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

func ReadPrimaryKey(pkFile string) ([]byte, error) {
	jsonFile, err := os.Open(pkFile)
	defer jsonFile.Close()
	if err != nil {
		return nil, err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var privateKey PrivateKey
	json.Unmarshal(byteValue, &privateKey)

	var decodedPrivateKey = make([]byte, 33)

	_, err = base64.StdEncoding.Decode(decodedPrivateKey, []byte(privateKey.Data.Bytes))

	return decodedPrivateKey[0:32], err
}

func ReadPublicKey(pubFile string) ([]byte, error) {
	plainFile, err := os.Open(pubFile)
	defer plainFile.Close()
	if err != nil {
		return []byte{}, err
	}
	byteValue, _ := ioutil.ReadAll(plainFile)

	var publicKey = make([]byte, 33)

	_, err = base64.StdEncoding.Decode(publicKey, byteValue)

	return publicKey[0:32], err
}