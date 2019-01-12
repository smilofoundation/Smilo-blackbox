package crypt

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/twystd/tweetnacl-go"
)

var (
	log *logrus.Entry
)

func initLog() {
	log = logrus.WithFields(logrus.Fields{
		"app":     "blackbox",
		"package": "crypt",
	})
}

func init() {
	initLog()
}

func GenerateKeys(generateKeys string) {
	log.WithField("generateKeys", generateKeys).Info("Going to generate encryption keys")
	files := strings.Split(generateKeys, ",")
	for i := range files {
		keyPair, _ := tweetnacl.CryptoBoxKeyPair()
		WritePrivateKeyFile(base64.StdEncoding.EncodeToString(keyPair.SecretKey), files[i]+".key")
		WritePublicKeyFile(base64.StdEncoding.EncodeToString(keyPair.PublicKey), files[i]+".pub")
	}
}

// WritePrivateKeyFile creates a json file with the private key
func WritePrivateKeyFile(key string, filename string) error {
	targetObject := map[string]interface{}{
		"type": "unlocked",
		"data": map[string]interface{}{
			"bytes": key,
		},
	}
	jsonBytes, _ := json.Marshal(targetObject)

	dir, _ := os.Getwd() // gives us the source path
	path := filepath.Join(dir, "keys/"+filename)

	log.WithField("path", path).Info("Going to Write Private Key File")
	return ioutil.WriteFile(path, jsonBytes, os.ModePerm)
}

// WritePublicKeyFile creates a file with the pubKey
func WritePublicKeyFile(key string, filename string) error {
	dir, _ := os.Getwd() // gives us the source path
	path := filepath.Join(dir, "keys/"+filename)

	log.WithField("path", path).Info("Going to Write Public Key File")

	return ioutil.WriteFile(path, []byte(key), os.ModePerm)
}
