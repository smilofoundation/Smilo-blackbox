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
	"testing"

	"gopkg.in/urfave/cli.v1"

	"encoding/base64"

	"github.com/stretchr/testify/require"

	"Smilo-blackbox/src/crypt"
)

const configFile = "./config_test.conf"

func TestLoadConfig(t *testing.T) {
	app := cli.NewApp()
	Init(app)
	err := LoadConfig(configFile)
	require.Empty(t, err, "could not open config file")
}

func TestPublicPrivateKeysLoad(t *testing.T) {
	err := LoadConfig(configFile)
	require.Empty(t, err, "could not open config file")

	pubKeyFile := config.Keys.KeyData[0].PublicKeyFile
	require.True(t, len(pubKeyFile) > 0, "pubKeyFile len is zero")

	privKeyFile := config.Keys.KeyData[0].PrivateKeyFile
	require.True(t, len(privKeyFile) > 0, "privKeyFile len is zero")

	publicKey, err := ReadPublicKey(pubKeyFile)
	require.Empty(t, err, "could not open public key")
	require.True(t, len(publicKey) > 0, "publicKey len is zero")

	configPrivateKey, err := ReadPrimaryKey(privKeyFile)
	require.Empty(t, err, "could not open private key")
	require.True(t, len(configPrivateKey) > 0, "configPrivateKey len is zero")

	privateKey := crypt.GetPrivateKey(publicKey)

	expected := base64.StdEncoding.EncodeToString(configPrivateKey)
	actual := base64.StdEncoding.EncodeToString(privateKey)

	require.Equal(t, expected, actual, "configPrivateKey should be equal to privateKey")
	Hostaddr.Value = config.Server.Hostaddr
	require.NotEmpty(t, Hostaddr.Value, "Host address should not be null")

}
