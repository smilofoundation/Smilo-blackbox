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
	pub := GetStringSlice(PublicKeysStr)
	require.True(t, len(pub) > 0)

	priv := GetStringSlice(PrivateKeysStr)
	require.True(t, len(priv) > 0)

	publicKey, err := ReadPublicKey(pub[0])
	require.Empty(t, err, "could not open public key")

	configPrivateKey, err := ReadPrimaryKey(priv[0])
	require.Empty(t, err, "could not open private key")

	privateKey := crypt.GetPrivateKey(publicKey)

	require.NotEqual(t, string(configPrivateKey), string(privateKey), "could not open config file")

}
