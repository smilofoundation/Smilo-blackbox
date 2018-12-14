package config

import (
	"os"
	"testing"

	"Smilo-blackbox/src/crypt"
)

const configFile = "./config_test.conf"

func TestMain(m *testing.M) {
	retcode := m.Run()
	os.Exit(retcode)
}

func TestLoadConfig(t *testing.T) {
	err := LoadConfig(configFile)

	if err != nil {
		t.Fail()
	}

	publicKey, err := ReadPublicKey("./public.pub")
	if err != nil {
		t.Fail()
	}

	configPrivateKey, err := ReadPrimaryKey("./private.key")
	if err != nil {
		t.Fail()
	}

	privateKey := crypt.GetPrivateKey(publicKey)

	if string(privateKey) != string(configPrivateKey) {
		t.Fail()
	}
}
