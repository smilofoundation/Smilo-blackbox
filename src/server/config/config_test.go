package config

import (
	"os"
	"testing"

	"Smilo-blackbox/src/crypt"
	"github.com/stretchr/testify/require"
)

const configFile = "./config_test.conf"

func TestMain(m *testing.M) {
	retcode := m.Run()
	os.Exit(retcode)
}

func TestLoadConfig(t *testing.T) {
	err := LoadConfig(configFile)
	require.Empty(t, err, "could not open config file")

	publicKey, err := ReadPublicKey("./public.pub")
	require.Empty(t, err, "could not open public key")

	configPrivateKey, err := ReadPrimaryKey("./private.key")
	require.Empty(t, err, "could not open private key")

	privateKey := crypt.GetPrivateKey(publicKey)

	require.Equal(t, string(configPrivateKey), string(privateKey), "could not open config file")

}
