package config

import (
	"testing"
	"os"
	"github.com/spf13/viper"
)

const configFile = "./config_test.conf"

func TestMain(m *testing.M) {
	retcode := m.Run()
	os.Exit(retcode)
}

func TestCommandLine(t *testing.T) {
	err := LoadConfig(configFile)

	if err != nil {
		t.Fatalf("Unable to load config file: %s, %s", configFile, err)
	}

	privateKey := viper.GetStringMap("keys")["keydata"].([]interface {})[0].(map[string]interface{})["config"].(string);
	publicKey := viper.GetStringMap("keys")["keydata"].([]interface {})[0].(map[string]interface{})["publicKey"].(string);

    if (privateKey != "./private.key") {
    	t.Fail()
	} else {
		_, err = ReadPrimaryKey(privateKey)
		if err != nil {
			t.Fail()
		}
	}

	if (publicKey != "./public.key") {
		teste,err2 := ReadPublicKey(publicKey)
		if err2 != nil {
			t.Fail()
		}
		_ = len(teste)
	}
}