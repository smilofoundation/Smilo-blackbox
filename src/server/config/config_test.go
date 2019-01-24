package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	"Smilo-blackbox/src/crypt"
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
	pub := PrivateKeys

	publicKey, err := ReadPublicKey(pub)
	require.Empty(t, err, "could not open public key")

	configPrivateKey, err := ReadPrimaryKey("./private.key")
	require.Empty(t, err, "could not open private key")

	privateKey := crypt.GetPrivateKey(publicKey)

	require.Equal(t, string(configPrivateKey), string(privateKey), "could not open config file")

}
