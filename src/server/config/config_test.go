package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	"Smilo-blackbox/src/crypt"
	"encoding/base64"
)

const configFile = "./config_test.toml"

func TestLoadConfig(t *testing.T) {
	err := ConfigLoad(configFile)
	require.Empty(t, err, "could not open config file")

	conf := AllSettings()
	require.NotEmpty(t, conf)
}

func TestPublicPrivateKeysLoad(t *testing.T) {
	err := ConfigLoad(configFile)
	require.Empty(t, err, "could not open config file")

	conf := AllSettings()
	require.NotEmpty(t, conf)
	PublicKeysStr := GetStringSlice(PublicKeysStr)
	require.True(t, len(PublicKeysStr) > 0, "PublicKeysStr len is zero")

	PrivateKeysStr := GetStringSlice(PrivateKeysStr)
	require.True(t, len(PrivateKeysStr) > 0, "PrivateKeysStr len is zero")

	publicKey, err := ReadPublicKey(PublicKeysStr[0])
	require.Empty(t, err, "could not open public key")
	require.True(t, len(publicKey) > 0, "publicKey len is zero")

	configPrivateKey, err := ReadPrimaryKey(PrivateKeysStr[0])
	require.Empty(t, err, "could not open private key")
	require.True(t, len(configPrivateKey) > 0, "configPrivateKey len is zero")

	crypt.PutKeyPair(crypt.KeyPair{PublicKey:publicKey,PrimaryKey:configPrivateKey})

	privateKey := crypt.GetPrivateKey(publicKey)

	expected := base64.StdEncoding.EncodeToString(configPrivateKey)
	actual := base64.StdEncoding.EncodeToString(privateKey)

	require.Equal(t, expected, actual, "configPrivateKey should be equal to privateKey")

}
