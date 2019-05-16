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

package crypt

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
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

// ComputePublicKey will compute a key based on the secret
func ComputePublicKey(secret []byte) ([]byte, error) {
	return tweetnacl.ScalarMultBase(secret)
}

// GenerateKeys will generate key/pub and save into file
func GenerateKeys(generateKeys string) error {
	log.WithField("generateKeys", generateKeys).Info("Going to generate encryption keys")
	files := strings.Split(generateKeys, ",")
	for i := range files {
		keyPair, err := tweetnacl.CryptoBoxKeyPair()
		if err != nil {
			log.WithError(err).Error("Could not tweetnacl.CryptoBoxKeyPair")
			return err
		}
		err = WritePrivateKeyFile(base64.StdEncoding.EncodeToString(keyPair.SecretKey), files[i]+".key")
		if err != nil {
			log.WithError(err).Error("Could not WritePrivateKeyFile")
			return err
		}
		err = WritePublicKeyFile(base64.StdEncoding.EncodeToString(keyPair.PublicKey), files[i]+".pub")
		if err != nil {
			log.WithError(err).Error("Could not WritePublicKeyFile")
			return err
		}
	}
	return nil
}

// WritePrivateKeyFile creates a json file with the private key
func WritePrivateKeyFile(key string, filename string) error {
	targetObject := map[string]interface{}{
		"type": "unlocked",
		"data": map[string]interface{}{
			"bytes": key,
		},
	}
	jsonBytes, err := json.Marshal(targetObject)
	if err != nil {
		log.WithError(err).Error("Could not json.Marshal")
		return err
	}

	dir, err := os.Getwd() // gives us the source path
	if err != nil {
		log.WithError(err).Error("Could not os.Getwd()")
		return err
	}

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
