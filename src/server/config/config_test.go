package config

import (
	"testing"

	"encoding/base64"

	"github.com/stretchr/testify/require"

	"Smilo-blackbox/src/crypt"
)

const configFile = "./config_test.conf"

func TestLoadConfig(t *testing.T) {
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

}
